package permission

import "errors"

var (
	ErrPermissionNotFound    = errors.New("permission not found")
	ErrRoleNotFound          = errors.New("role not found")
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidInput          = errors.New("invalid input")
	ErrPermissionExists      = errors.New("permission already exists")
	ErrRoleExists            = errors.New("role already exists")
	ErrPermissionAssigned    = errors.New("permission is assigned to roles")
	ErrRoleAssigned          = errors.New("role is assigned to users")
	ErrPermissionNotAssigned = errors.New("permission not assigned to role")
	ErrRoleNotAssigned       = errors.New("role not assigned to user")
)
