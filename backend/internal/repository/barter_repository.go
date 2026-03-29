package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skillora/backend/internal/models"
)

// BarterRepository handles transactions and double-entry ledger writes.
type BarterRepository struct {
	db *pgxpool.Pool
}

// NewBarterRepository constructs the repo.
func NewBarterRepository(db *pgxpool.Pool) *BarterRepository {
	return &BarterRepository{db: db}
}

// CreateTransaction inserts a new pending barter agreement between two users.
func (r *BarterRepository) CreateTransaction(ctx context.Context, tx models.BarterTransaction) (*models.BarterTransaction, error) {
	const q = `
		INSERT INTO barter_transactions
			(initiator_id, receiver_id, initiator_skill_id, receiver_skill_id, credit_amount, status)
		VALUES ($1, $2, $3, $4, $5, 'pending')
		RETURNING id, initiator_id, receiver_id, initiator_skill_id, receiver_skill_id, credit_amount, status, created_at, updated_at
	`
	row := r.db.QueryRow(ctx, q,
		tx.InitiatorID, tx.ReceiverID,
		tx.InitiatorSkillID, tx.ReceiverSkillID,
		tx.CreditAmount,
	)
	var t models.BarterTransaction
	if err := row.Scan(
		&t.ID, &t.InitiatorID, &t.ReceiverID,
		&t.InitiatorSkillID, &t.ReceiverSkillID,
		&t.CreditAmount, &t.Status, &t.CreatedAt, &t.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("CreateTransaction: %w", err)
	}
	return &t, nil
}

// GetUserTransactions returns both incoming and outgoing barters for a user.
func (r *BarterRepository) GetUserTransactions(ctx context.Context, userID string) ([]models.BarterTransaction, error) {
	const q = `
		SELECT id, initiator_id, receiver_id, initiator_skill_id, receiver_skill_id, credit_amount, status, created_at, updated_at
		FROM barter_transactions
		WHERE initiator_id = $1 OR receiver_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("GetUserTransactions: %w", err)
	}
	defer rows.Close()

	var list []models.BarterTransaction
	for rows.Next() {
		var t models.BarterTransaction
		if err := rows.Scan(
			&t.ID, &t.InitiatorID, &t.ReceiverID,
			&t.InitiatorSkillID, &t.ReceiverSkillID,
			&t.CreditAmount, &t.Status, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("GetUserTransactions scan: %w", err)
		}
		list = append(list, t)
	}
	if list == nil {
		list = make([]models.BarterTransaction, 0)
	}
	return list, rows.Err()
}

// UpdateTransactionStatus transitions a barter to accepted, completed, or cancelled.
func (r *BarterRepository) UpdateTransactionStatus(ctx context.Context, txID string, status models.BarterStatus) error {
	const q = `UPDATE barter_transactions SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, q, string(status), txID)
	if err != nil {
		return fmt.Errorf("UpdateTransactionStatus: %w", err)
	}
	return nil
}

// GetUserCreditBalance computes the running credit balance from ledger_entries.
func (r *BarterRepository) GetUserCreditBalance(ctx context.Context, userID string) (int, error) {
	const q = `
		SELECT COALESCE(
			SUM(CASE WHEN entry_type = 'credit' THEN amount ELSE -amount END),
			0
		)
		FROM ledger_entries
		WHERE user_id = $1
	`
	var balance int
	if err := r.db.QueryRow(ctx, q, userID).Scan(&balance); err != nil {
		return 0, fmt.Errorf("GetUserCreditBalance: %w", err)
	}
	return balance, nil
}

// PostLedgerEntries atomically writes two ledger rows for a completed barter.
// One debit for the payer, one credit for the payee.
func (r *BarterRepository) PostLedgerEntries(ctx context.Context, txID, payerID, payeeID string, amount int) error {
	// Fetch running balance for both parties first.
	payerBal, err := r.GetUserCreditBalance(ctx, payerID)
	if err != nil {
		return err
	}
	payeeBal, err := r.GetUserCreditBalance(ctx, payeeID)
	if err != nil {
		return err
	}

	// Validate payer has sufficient balance.
	if payerBal < amount {
		return fmt.Errorf("insufficient credits: user %s has %d, needs %d", payerID, payerBal, amount)
	}

	// Use a pgx transaction to ensure atomicity.
	pgTx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("PostLedgerEntries begin: %w", err)
	}
	defer pgTx.Rollback(ctx)

	const insertQ = `
		INSERT INTO ledger_entries (transaction_id, user_id, entry_type, amount, balance_after)
		VALUES ($1, $2, $3, $4, $5)
	`

	// Debit payer
	if _, err := pgTx.Exec(ctx, insertQ, txID, payerID, "debit", amount, payerBal-amount); err != nil {
		return fmt.Errorf("PostLedgerEntries debit: %w", err)
	}

	// Credit payee
	if _, err := pgTx.Exec(ctx, insertQ, txID, payeeID, "credit", amount, payeeBal+amount); err != nil {
		return fmt.Errorf("PostLedgerEntries credit: %w", err)
	}

	return pgTx.Commit(ctx)
}

// CreateMilestones inserts a set of progress steps for a barter agreement.
func (r *BarterRepository) CreateMilestones(ctx context.Context, ms []models.Milestone) error {
	const q = `
		INSERT INTO milestones (barter_id, title, description, credit_portion, status)
		VALUES ($1, $2, $3, $4, 'pending')
	`
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, m := range ms {
		if _, err := tx.Exec(ctx, q, m.BarterID, m.Title, m.Description, m.CreditPortion); err != nil {
			return fmt.Errorf("CreateMilestone: %w", err)
		}
	}
	return tx.Commit(ctx)
}

// GetBarterMilestones retrieves the progress track for a specific exchange.
func (r *BarterRepository) GetBarterMilestones(ctx context.Context, barterID string) ([]models.Milestone, error) {
	const q = `
		SELECT id, barter_id, title, description, credit_portion, status, created_at, updated_at
		FROM milestones
		WHERE barter_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.db.Query(ctx, q, barterID)
	if err != nil {
		return nil, fmt.Errorf("GetBarterMilestones: %w", err)
	}
	defer rows.Close()

	var list []models.Milestone
	for rows.Next() {
		var m models.Milestone
		if err := rows.Scan(
			&m.ID, &m.BarterID, &m.Title, &m.Description, &m.CreditPortion, &m.Status, &m.CreatedAt, &m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, rows.Err()
}

// UpdateMilestoneStatus transitions a single milestone and optionally releases funds.
func (r *BarterRepository) UpdateMilestoneStatus(ctx context.Context, milestoneID string, status models.MilestoneStatus) error {
	const q = `UPDATE milestones SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, q, string(status), milestoneID)
	return err
}

// GetMilestone retrieves a single milestone by ID.
func (r *BarterRepository) GetMilestone(ctx context.Context, id string) (*models.Milestone, error) {
	const q = `SELECT id, barter_id, title, description, credit_portion, status, created_at, updated_at FROM milestones WHERE id = $1`
	var m models.Milestone
	err := r.db.QueryRow(ctx, q, id).Scan(
		&m.ID, &m.BarterID, &m.Title, &m.Description, &m.CreditPortion, &m.Status, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
