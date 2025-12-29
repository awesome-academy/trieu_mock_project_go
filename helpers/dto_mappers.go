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

func MapUserToUserSummary(user *models.User) *dtos.UserSummary {
	if user == nil {
		return nil
	}
	return &dtos.UserSummary{
		ID:   user.ID,
		Name: user.Name,
	}
}

func MapUsersToUserSummaries(users []models.User) []dtos.UserSummary {
	summaries := make([]dtos.UserSummary, 0, len(users))
	for _, user := range users {
		summaries = append(summaries, *MapUserToUserSummary(&user))
	}
	return summaries
}

func MapTeamToTeamDto(team *models.Team) *dtos.Team {
	if team == nil {
		return nil
	}
	return &dtos.Team{
		ID:          team.ID,
		Name:        team.Name,
		Description: team.Description,
		CreatedAt:   team.CreatedAt,
		UpdatedAt:   team.UpdatedAt,

		Leader:   *MapUserToUserSummary(&team.Leader),
		Members:  MapUsersToUserSummaries(team.Members),
		Projects: MapProjectsToProjectSummaries(team.Projects),
	}
}

func MapTeamsToTeamDtos(teams []models.Team) []dtos.Team {
	teamDtos := make([]dtos.Team, 0, len(teams))
	for _, team := range teams {
		teamDtos = append(teamDtos, *MapTeamToTeamDto(&team))
	}
	return teamDtos
}

func MapTeamMemberToTeamMemberSummary(member *models.TeamMember) *dtos.TeamMemberSummary {
	if member == nil {
		return nil
	}
	return &dtos.TeamMemberSummary{
		ID:       member.User.ID,
		Name:     member.User.Name,
		Email:    member.User.Email,
		JoinedAt: member.JoinedAt,
	}
}

func MapTeamMembersToTeamMemberSummaries(members []models.TeamMember) []dtos.TeamMemberSummary {
	teamMemberSummaries := make([]dtos.TeamMemberSummary, 0, len(members))
	for _, member := range members {
		teamMemberSummaries = append(teamMemberSummaries, *MapTeamMemberToTeamMemberSummary(&member))
	}
	return teamMemberSummaries
}
