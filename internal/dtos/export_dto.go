package dtos

type ExportRequest struct {
	Type string `form:"type" binding:"required,oneof=user position project skill team"`
}
