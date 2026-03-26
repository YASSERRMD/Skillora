package models

// SkillMatch is the output from a vector similarity search for the matching engine.
type SkillMatch struct {
	SkillID         string  `json:"skill_id"`
	SkillName       string  `json:"skill_name"`
	CategoryName    string  `json:"category_name"`
	OwnerID         string  `json:"owner_id"`
	ProficiencyLevel int    `json:"proficiency_level"`
	CreditValue     int     `json:"credit_value"`
	SimilarityScore float64 `json:"similarity_score"` // 0.0 to 1.0 (cosine similarity)
}
