package contracts

type EmailJobDTO struct {
	To           string      `json:"to"`
	Subject      string      `json:"subject"`
	TemplateName string      `json:"template_name"`
	Data         interface{} `json:"data"`
}
