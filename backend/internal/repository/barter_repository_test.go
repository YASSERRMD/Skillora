package repository_test

import (
	"testing"

	"github.com/skillora/backend/internal/models"
)

// TestBarterStatus_Values validates the consts are as expected so they never
// accidentally deviate from the DB CHECK constraint.
func TestBarterStatus_Values(t *testing.T) {
	cases := map[models.BarterStatus]string{
		models.BarterStatusPending:   "pending",
		models.BarterStatusAccepted:  "accepted",
		models.BarterStatusCompleted: "completed",
		models.BarterStatusCancelled: "cancelled",
	}
	for status, expected := range cases {
		if string(status) != expected {
			t.Errorf("BarterStatus %q expected %q", status, expected)
		}
	}
}

// TestLedgerEntryType_Values validates entry type constants match DB CHECK constraints.
func TestLedgerEntryType_Values(t *testing.T) {
	cases := map[models.LedgerEntryType]string{
		models.LedgerDebit:  "debit",
		models.LedgerCredit: "credit",
	}
	for entryType, expected := range cases {
		if string(entryType) != expected {
			t.Errorf("LedgerEntryType %q expected %q", entryType, expected)
		}
	}
}
