package dtos

type CreateOrUpdateSkillRequest struct {
	Name string `json:"name" binding:"required,max=255"`
}

type SkillSearchResponse struct {
	Skills []SkillSummary     `json:"skills"`
	Page   PaginationResponse `json:"page"`
}
