package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/skillora/backend/internal/llm"
)

// MilestoneDraft is the structure the AI returns per progress step.
type MilestoneDraft struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	CreditPortion int    `json:"credit_portion"`
}

// MilestoneAgent uses the LLM router to architect a curriculum/plan for a skill exchange.
type MilestoneAgent struct {
	router *llm.Router
}

func NewMilestoneAgent(router *llm.Router) *MilestoneAgent {
	return &MilestoneAgent{router: router}
}

// PlanExchange generates a 3-step milestone plan for the given skill and total credit value.
func (a *MilestoneAgent) PlanExchange(ctx context.Context, skillName, description string, totalCredits int) ([]MilestoneDraft, error) {
	prompt := fmt.Sprintf(`
You are the Skillora Curriculum Architect.
A user is offering to teach: "%s".
Their description: "%s".
Total barter value: %d credits.

Break this exchange into exactly 3 logical milestones (steps).
Distribution of credits must sum exactly to %d.
Output ONLY a JSON array of objects with "title", "description", and "credit_portion".
`, skillName, description, totalCredits, totalCredits)

	resp, err := a.router.GenerateJSON(ctx, "course_generation", prompt)
	if err != nil {
		return nil, fmt.Errorf("MilestoneAgent GenerateJSON: %w", err)
	}

	// Clean up potential markdown formatting in AI response.
	cleanJSON := strings.TrimSpace(resp)
	cleanJSON = strings.TrimPrefix(cleanJSON, "```json")
	cleanJSON = strings.TrimSuffix(cleanJSON, "```")
	cleanJSON = strings.TrimSpace(cleanJSON)

	var drafts []MilestoneDraft
	if err := json.Unmarshal([]byte(cleanJSON), &drafts); err != nil {
		return nil, fmt.Errorf("MilestoneAgent unmarshal: %w, resp was: %s", err, cleanJSON)
	}

	return drafts, nil
}
