package dtos

type Position struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
}

type PositionSearchResponse struct {
	Positions []Position         `json:"positions"`
	Page      PaginationResponse `json:"page"`
}

type CreateOrUpdatePositionRequest struct {
	Name         string `json:"name" binding:"required,max=255"`
	Abbreviation string `json:"abbreviation" binding:"required,max=50"`
}
