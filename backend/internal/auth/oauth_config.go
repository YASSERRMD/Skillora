// Package auth handles Google OAuth2 login flow and session creation.
package auth

import (
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/skillora/backend/internal/config"
)

// NewGoogleOAuthConfig creates the OAuth2 config for Google.
func NewGoogleOAuthConfig() *oauth2.Config {
	backendURL := fmt.Sprintf("http://localhost:%s", config.C.Port)
	return &oauth2.Config{
		ClientID:     config.C.GoogleClientID,
		ClientSecret: config.C.GoogleClientSecret,
		RedirectURL:  backendURL + "/api/v1/auth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}
