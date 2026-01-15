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
