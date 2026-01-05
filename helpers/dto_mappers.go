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

func MapTeamsToTeamSummaries(teams []models.Team) []dtos.TeamSummary {
	summaries := make([]dtos.TeamSummary, 0, len(teams))
	for _, team := range teams {
		summary := MapTeamToTeamSummary(&team)
		if summary != nil {
			summaries = append(summaries, *summary)
		}
	}
	return summaries
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
		summary := MapProjectToProjectSummary(&project)
		if summary != nil {
			summaries = append(summaries, *summary)
		}
	}
	return summaries
}

func MapUserSkillToUserSkillSummary(userSkill *models.UserSkill) *dtos.UserSkillSummary {
	if userSkill == nil {
		return nil
	}
	return &dtos.UserSkillSummary{
		ID:             userSkill.SkillID,
		Name:           userSkill.Skill.Name,
		Level:          userSkill.Level,
		UsedYearNumber: userSkill.UsedYearNumber,
	}
}

func MapUserSkillsToUserSkillSummaries(skills []models.UserSkill) []dtos.UserSkillSummary {
	summaries := make([]dtos.UserSkillSummary, 0, len(skills))
	for _, skill := range skills {
		summary := MapUserSkillToUserSkillSummary(&skill)
		if summary != nil {
			summaries = append(summaries, *summary)
		}
	}
	return summaries
}

func MapUserToUserProfile(user *models.User) *dtos.UserProfile {

	currentTeam := MapTeamToTeamSummary(user.CurrentTeam)

	projects := MapProjectsToProjectSummaries(user.Projects)

	skills := MapUserSkillsToUserSkillSummaries(user.UserSkill)
	var birthday *types.Date
	if user.Birthday != nil {
		birthday = &types.Date{Time: *user.Birthday}
	}
	return &dtos.UserProfile{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Birthday:    birthday,
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
		summary := MapUserToUserSummary(&user)
		if summary != nil {
			summaries = append(summaries, *summary)
		}
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
		dto := MapTeamToTeamDto(&team)
		if dto != nil {
			teamDtos = append(teamDtos, *dto)
		}
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
		summary := MapTeamMemberToTeamMemberSummary(&member)
		if summary != nil {
			teamMemberSummaries = append(teamMemberSummaries, *summary)
		}
	}
	return teamMemberSummaries
}

func MapTeamMemberToTeamMemberHistory(member *models.TeamMember) *dtos.TeamMemberHistory {
	if member == nil {
		return nil
	}
	return &dtos.TeamMemberHistory{
		ID:       member.ID,
		UserID:   member.User.ID,
		UserName: member.User.Name,
		JoinedAt: member.JoinedAt,
		LeftAt:   member.LeftAt,
	}
}

func MapTeamMembersToTeamMemberHistories(members []models.TeamMember) []dtos.TeamMemberHistory {
	histories := make([]dtos.TeamMemberHistory, 0, len(members))
	for _, member := range members {
		teamMemberHistory := MapTeamMemberToTeamMemberHistory(&member)
		if teamMemberHistory != nil {
			histories = append(histories, *teamMemberHistory)
		}
	}
	return histories
}

func MapPositionToPositionSummary(position *models.Position) *dtos.PositionSummary {
	if position == nil {
		return nil
	}
	return &dtos.PositionSummary{
		ID:   position.ID,
		Name: position.Name,
	}
}

func MapPositionsToPositionSummaries(positions []models.Position) []dtos.PositionSummary {
	summaries := make([]dtos.PositionSummary, 0, len(positions))
	for _, position := range positions {
		summary := MapPositionToPositionSummary(&position)
		if summary != nil {
			summaries = append(summaries, *summary)
		}
	}
	return summaries
}

func MapPositionToPositionDto(position *models.Position) *dtos.Position {
	if position == nil {
		return nil
	}
	return &dtos.Position{
		ID:           position.ID,
		Name:         position.Name,
		Abbreviation: position.Abbreviation,
	}
}

func MapPositionsToPositionDtos(positions []models.Position) []dtos.Position {
	positionDtos := make([]dtos.Position, 0, len(positions))
	for _, position := range positions {
		dto := MapPositionToPositionDto(&position)
		if dto != nil {
			positionDtos = append(positionDtos, *dto)
		}
	}
	return positionDtos
}

func MapUserToUserDataForSearch(user *models.User) *dtos.UserDataForSearch {
	if user == nil {
		return nil
	}
	return &dtos.UserDataForSearch{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		CurrentTeam: MapTeamToTeamSummary(user.CurrentTeam),
	}
}

func MapUsersToUserDataForSearches(users []models.User) []dtos.UserDataForSearch {
	userDtos := make([]dtos.UserDataForSearch, 0, len(users))
	for _, user := range users {
		dto := MapUserToUserDataForSearch(&user)
		if dto != nil {
			userDtos = append(userDtos, *dto)
		}
	}
	return userDtos
}

func MapSkillToSkillSummary(skill *models.Skill) *dtos.SkillSummary {
	if skill == nil {
		return nil
	}
	return &dtos.SkillSummary{
		ID:   skill.ID,
		Name: skill.Name,
	}
}

func MapSkillsToSkillSummaries(skills []models.Skill) []dtos.SkillSummary {
	skillSummaries := make([]dtos.SkillSummary, 0, len(skills))
	for _, skill := range skills {
		skillSummary := MapSkillToSkillSummary(&skill)
		if skillSummary != nil {
			skillSummaries = append(skillSummaries, *skillSummary)
		}
	}
	return skillSummaries
}
