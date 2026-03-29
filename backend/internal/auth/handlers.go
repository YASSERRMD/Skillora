package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"github.com/skillora/backend/internal/config"
	"github.com/skillora/backend/internal/crypto"
	"github.com/skillora/backend/internal/db"
	"github.com/skillora/backend/internal/repository"
)

// GoogleUserInfo is the profile returned by Google's userinfo endpoint.
type GoogleUserInfo struct {
	Sub        string `json:"sub"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	Picture    string `json:"picture"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
}

// Handler holds OAuth dependencies.
type Handler struct {
	oauthCfg *oauth2.Config
	userRepo *repository.UserRepository
}

// NewHandler constructs an auth Handler.
func NewHandler(oauthCfg *oauth2.Config, userRepo *repository.UserRepository) *Handler {
	return &Handler{oauthCfg: oauthCfg, userRepo: userRepo}
}

// GoogleLogin redirects the user to Google's OAuth2 consent screen.
// GET /api/v1/auth/google/login
func (h *Handler) GoogleLogin(c *gin.Context) {
	// Dev Auth Bypass
	if config.C.GoogleClientID == "your_google_client_id_here" || config.C.GoogleClientID == "dummy_google_client_id_for_dev" {
		ctx := c.Request.Context()
		user, _, err := h.userRepo.UpsertGoogleUser(ctx, "dev_sub_000", "dev@skillora.local", "Local Developer", "https://api.dicebear.com/7.x/avataaars/svg?seed=Dev")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "dev user upsert failed"})
			return
		}
		jwtToken, err := crypto.GenerateJWT(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "dev token generation failed"})
			return
		}
		c.SetCookie("skillora_token", jwtToken, int((7 * 24 * time.Hour).Seconds()), "/", "", false, true)
		c.Redirect(http.StatusTemporaryRedirect, config.C.FrontendURL+"/dashboard")
		return
	}

	state, err := generateState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "state generation failed"})
		return
	}

	// Cache state in Redis with 5-minute TTL to prevent CSRF.
	if err := db.SetJSON(c.Request.Context(), "oauth_state:"+state, true, 5*time.Minute); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to cache state"})
		return
	}

	url := h.oauthCfg.AuthCodeURL(state, oauth2.AccessTypeOnline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback handles the OAuth2 callback from Google.
// GET /api/v1/auth/google/callback
func (h *Handler) GoogleCallback(c *gin.Context) {
	ctx := c.Request.Context()

	// Validate CSRF state.
	state := c.Query("state")
	var cached bool
	if err := db.GetJSON(ctx, "oauth_state:"+state, &cached); err != nil || !cached {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired state"})
		return
	}
	// Delete state immediately (single-use).
	db.RDB.Del(ctx, "oauth_state:"+state)

	// Exchange auth code for access token.
	code := c.Query("code")
	token, err := h.oauthCfg.Exchange(ctx, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token exchange failed"})
		return
	}

	// Fetch Google user profile.
	userInfo, err := fetchGoogleUserInfo(ctx, h.oauthCfg, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user info"})
		return
	}

	// Upsert user into DB.
	user, isNew, err := h.userRepo.UpsertGoogleUser(ctx, userInfo.Sub, userInfo.Email, userInfo.Name, userInfo.Picture)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user upsert failed"})
		return
	}

	// Mint 10 credits for new users (Phase 25 — called here for correctness).
	if isNew {
		// Wallet minting is implemented in Phase 25; placeholder log for now.
		_ = isNew
	}

	// Generate JWT and set as HttpOnly cookie.
	jwtToken, err := crypto.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	c.SetCookie(
		"skillora_token",
		jwtToken,
		int((7 * 24 * time.Hour).Seconds()),
		"/",
		"",    // domain — empty = current host
		true,  // Secure (HTTPS only in prod)
		true,  // HttpOnly
	)

	// SameSite=Strict cookie header (Gin doesn't expose SameSite directly, set header).
	c.Header("Set-Cookie", fmt.Sprintf(
		"skillora_token=%s; Path=/; Max-Age=%d; HttpOnly; SameSite=Strict",
		jwtToken, int((7*24*time.Hour).Seconds()),
	))

	c.Redirect(http.StatusTemporaryRedirect, config.C.FrontendURL+"/dashboard")
}

// generateState creates a cryptographically random state string.
func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// fetchGoogleUserInfo retrieves the authenticated user's Google profile.
func fetchGoogleUserInfo(ctx context.Context, cfg *oauth2.Config, token *oauth2.Token) (*GoogleUserInfo, error) {
	client := cfg.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, fmt.Errorf("userinfo request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("userinfo: unexpected status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("userinfo read: %w", err)
	}

	var info GoogleUserInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("userinfo unmarshal: %w", err)
	}
	return &info, nil
}
