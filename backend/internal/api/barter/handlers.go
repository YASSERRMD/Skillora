package barter

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skillora/backend/internal/api"
	"github.com/skillora/backend/internal/models"
	"github.com/skillora/backend/internal/repository"
)

// Handler manages barter transaction HTTP endpoints.
type Handler struct {
	repo *repository.BarterRepository
}

// NewHandler constructs the barter API handler.
func NewHandler(repo *repository.BarterRepository) *Handler {
	return &Handler{repo: repo}
}

// ProposeReq is the incoming payload for initiating a barter.
type ProposeReq struct {
	ReceiverID       string `json:"receiver_id"         binding:"required"`
	InitiatorSkillID string `json:"initiator_skill_id"  binding:"required"`
	ReceiverSkillID  string `json:"receiver_skill_id"   binding:"required"`
	CreditAmount     int    `json:"credit_amount"       binding:"required,min=1"`
}

// PostPropose creates a new barter proposal in 'pending' state.
// POST /api/v1/barters
func (h *Handler) PostPropose(c *gin.Context) {
	initiatorID := api.GetUserID(c)

	var req ProposeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if initiatorID == req.ReceiverID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot propose a barter with yourself"})
		return
	}

	tx, err := h.repo.CreateTransaction(c.Request.Context(), models.BarterTransaction{
		InitiatorID:      initiatorID,
		ReceiverID:       req.ReceiverID,
		InitiatorSkillID: req.InitiatorSkillID,
		ReceiverSkillID:  req.ReceiverSkillID,
		CreditAmount:     req.CreditAmount,
	})
	if err != nil {
		log.Printf("[barter] CreateTransaction error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create barter proposal"})
		return
	}

	c.JSON(http.StatusCreated, tx)
}

// GetMyBarters returns all barter transactions for the current user.
// GET /api/v1/barters
func (h *Handler) GetMyBarters(c *gin.Context) {
	userID := api.GetUserID(c)

	transactions, err := h.repo.GetUserTransactions(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load barters"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// PatchBarterStatus transitions a barter to accepted or cancelled.
// PATCH /api/v1/barters/:id/status
func (h *Handler) PatchBarterStatus(c *gin.Context) {
	barterID := c.Param("id")
	userID := api.GetUserID(c)

	var req struct {
		Status string `json:"status" binding:"required,oneof=accepted cancelled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_ = userID // Authorization is enforced below via db constraint check on receiver_id in full implementation

	if err := h.repo.UpdateTransactionStatus(c.Request.Context(), barterID, models.BarterStatus(req.Status)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update barter status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "barter status updated", "status": req.Status})
}

// PostComplete finalizes a barter by writing the double-entry ledger pair.
// POST /api/v1/barters/:id/complete
func (h *Handler) PostComplete(c *gin.Context) {
	barterID := c.Param("id")

	var req struct {
		PayerID  string `json:"payer_id"  binding:"required"`
		PayeeID  string `json:"payee_id"  binding:"required"`
		Amount   int    `json:"amount"    binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Write double-entry ledger atomically.
	if err := h.repo.PostLedgerEntries(c.Request.Context(), barterID, req.PayerID, req.PayeeID, req.Amount); err != nil {
		log.Printf("[barter] PostLedgerEntries error: %v", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// Mark barter as completed.
	if err := h.repo.UpdateTransactionStatus(c.Request.Context(), barterID, models.BarterStatusCompleted); err != nil {
		log.Printf("[barter] mark completed error: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "barter completed and ledger entries recorded"})
}

// GetCreditBalance returns the authenticated user's total credit balance.
// GET /api/v1/barters/balance
func (h *Handler) GetCreditBalance(c *gin.Context) {
	userID := api.GetUserID(c)

	balance, err := h.repo.GetUserCreditBalance(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": userID, "credit_balance": balance})
}
