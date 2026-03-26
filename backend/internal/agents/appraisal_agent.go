// Package agents holds autonomous agents that analyze, appraise, and facilitate.
package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/skillora/backend/internal/llm"
	"github.com/skillora/backend/internal/models"
)

// AppraisalAgent evaluates user-submitted skills to determine validity, proficiency,
// and assigns an estimated credit value for the barter economy.
type AppraisalAgent struct {
	router *llm.Router
}

// NewAppraisalAgent constructs the agent with a connection to the AI routing engine.
func NewAppraisalAgent(router *llm.Router) *AppraisalAgent {
	return &AppraisalAgent{router: router}
}

// AppraisalResult is the strict JSON output expected from the LLM.
type AppraisalResult struct {
	IsValidSkill bool   `json:"is_valid_skill"`
	Proficiency  int    `json:"proficiency"`     // 1 to 5 scale
	CreditValue  int    `json:"credit_value"`    // Suggested value (e.g. 5 to 50)
	Reasoning    string `json:"reasoning"`       // Explanation for the assessment
}

// DraftAppraisal queries the LLM to assess a newly offered skill by a user.
func (a *AppraisalAgent) DraftAppraisal(ctx context.Context, categoryName, skillName, description string) (*AppraisalResult, error) {
	// Construct the prompt using explicit JSON schema instructions.
	prompt := fmt.Sprintf(`
You are an expert appraiser for the Skillora platform, a knowledge barter economy.
A user has submitted a skill to offer. You must evaluate it.

Category: %s
Skill Name: %s
User Description: %s

Assess the skill based on:
1. Is it a valid, teachable/actionable skill? (true/false)
2. Based on the description, what is their self-evident proficiency? (1=beginner, 5=expert)
3. What is the estimated market value of an hour of this skill in credits? (range 5 to 50)
4. Provide a 1-sentence reasoning.

Output YOUR response ONLY in the following generic JSON format:
{
  "is_valid_skill": false,
  "proficiency": 0,
  "credit_value": 0,
  "reasoning": ""
}
`, categoryName, skillName, description)

	// Since this is an assessment task, UseCaseMediator or UseCaseGeneral applies. Using General.
	resultStr, err := a.router.GenerateJSON(ctx, models.UseCaseGeneral, prompt)
	if err != nil {
		return nil, fmt.Errorf("appraisal agent failed to generate assessment: %w", err)
	}

	// Clean up potential markdown formatting from unruly models (even if instructed not to).
	resultStr = strings.TrimPrefix(strings.TrimSpace(resultStr), "```json")
	resultStr = strings.TrimPrefix(resultStr, "```")
	resultStr = strings.TrimSuffix(resultStr, "```")

	var result AppraisalResult
	if err := json.Unmarshal([]byte(resultStr), &result); err != nil {
		return nil, fmt.Errorf("appraisal agent received malformed JSON: %w (raw: %s)", err, resultStr)
	}

	return &result, nil
}
