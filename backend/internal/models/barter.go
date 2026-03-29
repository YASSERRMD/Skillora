package models

import "time"

// BarterStatus represents the lifecycle state of a barter transaction.
type BarterStatus string

const (
	BarterStatusPending   BarterStatus = "pending"
	BarterStatusAccepted  BarterStatus = "accepted"
	BarterStatusCompleted BarterStatus = "completed"
	BarterStatusCancelled BarterStatus = "cancelled"
)

// BarterTransaction represents a skill exchange agreement between two participants.
type BarterTransaction struct {
	ID               string       `db:"id"                 json:"id"`
	InitiatorID      string       `db:"initiator_id"       json:"initiator_id"`
	ReceiverID       string       `db:"receiver_id"        json:"receiver_id"`
	InitiatorSkillID string       `db:"initiator_skill_id" json:"initiator_skill_id"`
	ReceiverSkillID  string       `db:"receiver_skill_id"  json:"receiver_skill_id"`
	CreditAmount     int          `db:"credit_amount"      json:"credit_amount"`
	Status           BarterStatus `db:"status"             json:"status"`
	CreatedAt        time.Time    `db:"created_at"         json:"created_at"`
	UpdatedAt        time.Time    `db:"updated_at"         json:"updated_at"`
}

// MilestoneStatus represents the lifecycle of a single progress milestone.
type MilestoneStatus string

const (
	MilestoneStatusPending   MilestoneStatus = "pending"
	MilestoneStatusCompleted MilestoneStatus = "completed"
	MilestoneStatusApproved  MilestoneStatus = "approved"
)

// Milestone is a staged part of a larger barter transaction agreement.
type Milestone struct {
	ID            string          `db:"id"             json:"id"`
	BarterID      string          `db:"barter_id"      json:"barter_id"`
	Title         string          `db:"title"          json:"title"`
	Description   string          `db:"description"    json:"description"`
	CreditPortion int             `db:"credit_portion" json:"credit_portion"`
	Status        MilestoneStatus `db:"status"         json:"status"`
	CreatedAt     time.Time       `db:"created_at"     json:"created_at"`
	UpdatedAt     time.Time       `db:"updated_at"     json:"updated_at"`
}

// LedgerEntryType discriminates debits from credits.
type LedgerEntryType string

const (
	LedgerDebit  LedgerEntryType = "debit"
	LedgerCredit LedgerEntryType = "credit"
)

// LedgerEntry is a single accounting line in the barter economy.
type LedgerEntry struct {
	ID            string          `db:"id"             json:"id"`
	TransactionID string          `db:"transaction_id" json:"transaction_id"`
	UserID        string          `db:"user_id"        json:"user_id"`
	EntryType     LedgerEntryType `db:"entry_type"     json:"entry_type"`
	Amount        int             `db:"amount"         json:"amount"`
	BalanceAfter  int             `db:"balance_after"  json:"balance_after"`
	CreatedAt     time.Time       `db:"created_at"     json:"created_at"`
}
