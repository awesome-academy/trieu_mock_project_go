package dtos

type ProjectDeadlineReminderEmailDTO struct {
	To          string
	UserName    string
	ProjectName string
	DueDate     string
}
