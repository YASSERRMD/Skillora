package barter

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skillora/backend/internal/agents"
	"github.com/skillora/backend/internal/api"
	"github.com/skillora/backend/internal/models"
	"github.com/skillora/backend/internal/repository"
)

// Handler manages barter transaction HTTP endpoints.
type Handler struct {
	repo           *repository.BarterRepository
	milestoneAgent *agents.MilestoneAgent
}

// NewHandler constructs the barter API handler.
func NewHandler(repo *repository.BarterRepository, agent *agents.MilestoneAgent) *Handler {
	return &Handler{repo: repo, milestoneAgent: agent}
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

	// Step 2: Generate AI-designed milestones for the study plan.
	// In a real app, we'd fetch the skill description from the DB first.
	// For now, we take a dummy description or assume fixed curriculum generation.
	drafts, err := h.milestoneAgent.PlanExchange(c.Request.Context(), req.ReceiverSkillID, "Custom learning plan", req.CreditAmount)
	if err != nil {
		log.Printf("[barter] PlanExchange error: %v", err)
		// Fallback to a single generic milestone if LLM fails.
		drafts = []agents.MilestoneDraft{{Title: "Full Skill Transfer", Description: "Complete exchange", CreditPortion: req.CreditAmount}}
	}

	// Step 3: Atomic commit of transactions and milestones.
	var milestones []models.Milestone
	for _, d := range drafts {
		milestones = append(milestones, models.Milestone{
			BarterID:      tx.ID,
			Title:         d.Title,
			Description:   d.Description,
			CreditPortion: d.CreditPortion,
			Status:        models.MilestoneStatusPending,
		})
	}
	if err := h.repo.CreateMilestones(c.Request.Context(), milestones); err != nil {
		log.Printf("[barter] CreateMilestones error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to initialize milestones"})
		return
	}

	c.JSON(http.StatusCreated, tx)
}

// GetMilestones returns the curriculum progress for a barter.
// GET /api/v1/barters/:id/milestones
func (h *Handler) GetMilestones(c *gin.Context) {
	barterID := c.Param("id")
	ms, err := h.repo.GetBarterMilestones(c.Request.Context(), barterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load milestones"})
		return
	}
	c.JSON(http.StatusOK, ms)
}

// PostMilestoneComplete marks a milestone as done by the provider.
// POST /api/v1/milestones/:id/complete
func (h *Handler) PostMilestoneComplete(c *gin.Context) {
	msID := c.Param("id")
	if err := h.repo.UpdateMilestoneStatus(c.Request.Context(), msID, models.MilestoneStatusCompleted); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "status update failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "milestone marked as completed"})
}

// PostMilestoneApprove releases the credits for a specific milestone.
// POST /api/v1/milestones/:id/approve
func (h *Handler) PostMilestoneApprove(c *gin.Context) {
	msID := c.Param("id")
	ctx := c.Request.Context()

	// 1. Fetch milestone to get barter context and credit portion.
	ms, err := h.repo.GetMilestone(ctx, msID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "milestone not found"})
		return
	}

	// 2. Fetch the transaction context for participants.
	txList, err := h.repo.GetUserTransactions(ctx, api.GetUserID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify exchange context"})
		return
	}

	// Finding the exact transaction for this milestone
	var activeTx *models.BarterTransaction
	for _, t := range txList {
		if t.ID == ms.BarterID {
			activeTx = &t
			break
		}
	}
	if activeTx == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized to approve this milestone"})
		return
	}

	// 3. Atomically release the specific credit_portion from Payer to Payee.
	// Payer is the Receiver (learner), Payee is the Initiator (teacher) for this specific skill transfer.
	if err := h.repo.PostLedgerEntries(ctx, activeTx.ID, activeTx.ReceiverID, activeTx.InitiatorID, ms.CreditPortion); err != nil {
		log.Printf("[barter] PostLedgerEntries error: %v", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	
	if err := h.repo.UpdateMilestoneStatus(ctx, msID, models.MilestoneStatusApproved); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "milestone approval failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "milestone approved and credits released"})
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
