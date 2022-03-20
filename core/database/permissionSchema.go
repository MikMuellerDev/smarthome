package database

// This file defines which permissions exists and describes their attributes
type Permission struct {
	Permission  PermissionType `json:"permission"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
}

type PermissionType string

const (
	PermissionAuthentication          PermissionType = "authentication"
	PermissionGetUserSwitches         PermissionType = "getUserSwitches"
	PermissionSetPower                PermissionType = "setPower"
	PermissionDeleteOldLogs           PermissionType = "deleteOldLogs"
	PermissionDeleteAllLogs           PermissionType = "deleteAllLogs"
	PermissionListLogs                PermissionType = "listLogs"
	PermissionChangeAvatar            PermissionType = "changeAvatar"
	PermissionChangeUserPermissions   PermissionType = "changeUserPermissions"
	PermissionChangeSwitchPermissions PermissionType = "changeSwitchPermissions"
	PermissionGetDebugInfo            PermissionType = "getDebugInfo"
	PermissionChangeUsers             PermissionType = "changeUsers"

	// Dangerous
	PermissionWildCard PermissionType = "*"
)

func GetPermissions() []Permission {
	permissions := []Permission{
		{
			// User is allowed to authenticate and login, if disabled, a user is `disabled`
			Permission:  PermissionAuthentication,
			Name:        "Authentication",
			Description: "Allows the user to authenticate",
		},
		{
			// User is allowed to request a list of their personal switches (which they have access to)
			Permission:  PermissionGetUserSwitches,
			Name:        "Get Personal Switches",
			Description: "Get all allowed switches for current user.",
		},
		{
			// User is allowed to request power jobs, interact with outlets, still dependent on switch permissions
			Permission:  PermissionSetPower,
			Name:        "Set Power",
			Description: "Interact with switches",
		},
		{
			// (Admin) is allowed to delete logs older than 30 days
			Permission:  PermissionDeleteOldLogs,
			Name:        "Flush Old Logs",
			Description: "Delete logs events which are older than 30 days",
		},
		{
			// (Admin) is allowed to delete all logs
			Permission:  PermissionDeleteAllLogs,
			Name:        "Flush All Logs",
			Description: "Delete all log event records",
		},
		{
			// (Admin) is allowed to request logs
			Permission:  PermissionListLogs,
			Name:        "List Logs",
			Description: "List all internal logs",
		},
		{
			// User is allowed to upload a custom avatar
			Permission:  PermissionChangeAvatar,
			Name:        "Upload / Delete / Change Avatar",
			Description: "Allows the user to customize their avatar",
		},
		{
			// (Admin) is allowed to add / delete permissions to / from users
			Permission:  PermissionChangeUserPermissions,
			Name:        "Change User Permissions",
			Description: "Add / delete permissions to / from users",
		},
		{
			// (Admin) is allowed to add / delete switch permissions to / from users
			Permission:  PermissionChangeSwitchPermissions,
			Name:        "Change User Switch Permissions",
			Description: "Add / delete switch permissions to / from users",
		},
		{
			// (Admin) is allowed to read debug information from the server
			Permission:  PermissionGetDebugInfo,
			Name:        "Display Debug Info",
			Description: "Obtain debug information about the system",
		},
		{
			// (Admin) is allowed to create new users or delete users
			Permission:  PermissionChangeUsers,
			Name:        "Add / Delete users",
			Description: "Create a new user or delete users",
		},
		{
			// WARNING: This allows a user to do everything, should only be allowed to the `admin` user
			Permission:  PermissionWildCard,
			Name:        "Permission Wildcard *",
			Description: "Allows all permissions",
		},
	}
	return permissions
}
