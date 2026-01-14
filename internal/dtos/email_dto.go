package dtos

type TeamMembershipEmailDTO struct {
	To       string
	UserName string
	TeamName string
}

type ProjectMembershipEmailDTO struct {
	To          string
	UserName    string
	ProjectName string
}

type ProjectDeadlineReminderEmailDTO struct {
	To          string
	UserName    string
	ProjectName string
	DueDate     string
}

type EmailJobDTO struct {
	To           string      `json:"to"`
	Subject      string      `json:"subject"`
	TemplateName string      `json:"template_name"`
	Data         interface{} `json:"data"`
}
