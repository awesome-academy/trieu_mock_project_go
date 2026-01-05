package types

type ActivityLog struct {
	Value             string
	DescriptionFormat string
}

var (
	UserSignIn     = ActivityLog{Value: "USER_SIGN_IN", DescriptionFormat: "User [%d - %s] signed in"}
	UserSignOut    = ActivityLog{Value: "USER_SIGN_OUT", DescriptionFormat: "User [%d - %s] signed out"}
	AdminSignIn    = ActivityLog{Value: "ADMIN_SIGN_IN", DescriptionFormat: "User [%d - %s] signed in as admin"}
	AdminSignOut   = ActivityLog{Value: "ADMIN_SIGN_OUT", DescriptionFormat: "User [%d - %s] signed out as admin"}
	CreateUser     = ActivityLog{Value: "CREATE_USER", DescriptionFormat: "User [%d - %s] was created new user(ID: %d)"}
	UpdateUser     = ActivityLog{Value: "UPDATE_USER", DescriptionFormat: "User [%d - %s] was updated user (ID: %d)"}
	DeleteUser     = ActivityLog{Value: "DELETE_USER", DescriptionFormat: "User [%d - %s] was deleted user (ID: %d, Email: %s)"}
	CreatePosition = ActivityLog{Value: "CREATE_POSITION", DescriptionFormat: "User [%d - %s] created position (ID: %d, Name: %s)"}
	UpdatePosition = ActivityLog{Value: "UPDATE_POSITION", DescriptionFormat: "User [%d - %s] updated position (ID: %d, Name: %s)"}
	DeletePosition = ActivityLog{Value: "DELETE_POSITION", DescriptionFormat: "User [%d - %s] deleted position (ID: %d, Name: %s)"}
	CreateProject  = ActivityLog{Value: "CREATE_PROJECT", DescriptionFormat: "User [%d - %s] created project (ID: %d, Name: %s)"}
	UpdateProject  = ActivityLog{Value: "UPDATE_PROJECT", DescriptionFormat: "User [%d - %s] updated project (ID: %d, Name: %s)"}
	DeleteProject  = ActivityLog{Value: "DELETE_PROJECT", DescriptionFormat: "User [%d - %s] deleted project (ID: %d, Name: %s)"}
	CreateSkill    = ActivityLog{Value: "CREATE_SKILL", DescriptionFormat: "User [%d - %s] created skill (ID: %d, Name: %s)"}
	UpdateSkill    = ActivityLog{Value: "UPDATE_SKILL", DescriptionFormat: "User [%d - %s] updated skill (ID: %d, Name: %s)"}
	DeleteSkill    = ActivityLog{Value: "DELETE_SKILL", DescriptionFormat: "User [%d - %s] deleted skill (ID: %d, Name: %s)"}
	CreateTeam     = ActivityLog{Value: "CREATE_TEAM", DescriptionFormat: "User [%d - %s] created team (ID: %d, Name: %s)"}
	UpdateTeam     = ActivityLog{Value: "UPDATE_TEAM", DescriptionFormat: "User [%d - %s] updated team (ID: %d, Name: %s)"}
	DeleteTeam     = ActivityLog{Value: "DELETE_TEAM", DescriptionFormat: "User [%d - %s] deleted team (ID: %d, Name: %s)"}
	JoinTeam       = ActivityLog{Value: "JOIN_TEAM", DescriptionFormat: "User [%d - %s] updated user [%d - %s] join team (Team ID: %d)"}
	LeaveTeam      = ActivityLog{Value: "LEAVE_TEAM", DescriptionFormat: "User [%d - %s] updated user [%d - %s] leave team (Team ID: %d)"}
)
