package helpers

import (
	"trieu_mock_project_go/internal/dtos"
	"trieu_mock_project_go/models"
	"trieu_mock_project_go/types"
)

func MapTeamToTeamSummary(team *models.Team) *dtos.TeamSummary {
	if team == nil {
		return nil
	}
	return &dtos.TeamSummary{
		ID:   team.ID,
		Name: team.Name,
	}
}

func MapProjectToProjectSummary(project *models.Project) *dtos.ProjectSummary {
	if project == nil {
		return nil
	}
	return &dtos.ProjectSummary{
		ID:           project.ID,
		Name:         project.Name,
		Abbreviation: project.Abbreviation,
		StartDate:    project.StartDate,
		EndDate:      project.EndDate,
	}
}

func MapProjectsToProjectSummaries(projects []models.Project) []dtos.ProjectSummary {
	summaries := make([]dtos.ProjectSummary, 0, len(projects))
	for _, project := range projects {
		summaries = append(summaries, *MapProjectToProjectSummary(&project))
	}
	return summaries
}

func MapSkillToUserSkillSummary(skill *models.Skill) *dtos.UserSkillSummary {
	if skill == nil {
		return nil
	}
	return &dtos.UserSkillSummary{
		ID:   skill.ID,
		Name: skill.Name,
	}
}

func MapSkillsToUserSkillSummaries(skills []models.Skill) []dtos.UserSkillSummary {
	summaries := make([]dtos.UserSkillSummary, 0, len(skills))
	for _, skill := range skills {
		summaries = append(summaries, *MapSkillToUserSkillSummary(&skill))
	}
	return summaries
}

func MapUserToUserProfile(user *models.User) *dtos.UserProfile {

	currentTeam := MapTeamToTeamSummary(user.CurrentTeam)

	projects := MapProjectsToProjectSummaries(user.Projects)

	skills := MapSkillsToUserSkillSummaries(user.Skills)

	return &dtos.UserProfile{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Birthday:    &types.Date{Time: *user.Birthday},
		CurrentTeam: currentTeam,
		Position: dtos.Position{
			ID:           user.Position.ID,
			Name:         user.Position.Name,
			Abbreviation: user.Position.Abbreviation,
		},
		Projects: projects,
		Skills:   skills,
	}
}
