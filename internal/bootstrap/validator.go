package bootstrap

import (
	"trieu_mock_project_go/internal/dtos"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidations() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterStructValidation(dtos.ProjectRequestStructLevelValidation, dtos.CreateOrUpdateProjectRequest{})
	}
}
