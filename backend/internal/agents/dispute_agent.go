package agents

import (
	"context"
	"fmt"

	"github.com/skillora/backend/internal/llm"
)

// DisputeAgent mediates conflicts between skill exchangers using AI adjudication.
type DisputeAgent struct {
	router *llm.Router
}

func NewDisputeAgent(router *llm.Router) *DisputeAgent {
	return &DisputeAgent{router: router}
}

// Mediate resolves a conflict by deciding if credits should be released or returned.
func (a *DisputeAgent) Mediate(ctx context.Context, description, initiatorClaim, receiverDefense string) (string, bool, error) {
	prompt := fmt.Sprintf(`
You are the Skillora AI Dispute Mediator.
Exchange context: "%s".
Initiator Claim: "%s".
Receiver Defense: "%s".

Decide who is in the right based on common sense and the provided data.
You must return your response in the format:
Verdict: [Release / Refund / Partial]
Reasoning: [1-2 sentences]
`, description, initiatorClaim, receiverDefense)

	resp, err := a.router.GenerateJSON(ctx, "dispute_resolution", prompt)
	if err != nil {
		return "", false, err
	}

	return resp, true, nil
}
