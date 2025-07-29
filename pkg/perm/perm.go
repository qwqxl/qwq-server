package perm

type UserID string
type RoleID string
type Permission string

// 用户结构
type User struct {
	ID    UserID
	Roles []RoleID
}

// 角色结构
type Role struct {
	ID          RoleID
	Permissions []Permission
}

// 权限服务
type RBAC struct {
	users       map[UserID]*User
	roles       map[RoleID]*Role
	permissions map[Permission]struct{}
}

// 初始化和构造函数
func NewRBAC() *RBAC {
	return &RBAC{
		users:       make(map[UserID]*User),
		roles:       make(map[RoleID]*Role),
		permissions: make(map[Permission]struct{}),
	}
}

/* -------------- 用户/角色/权限 增删改查 API -------------- */

// 添加权限
func (r *RBAC) AddPermission(p Permission) {
	r.permissions[p] = struct{}{}
}

// 添加角色
func (r *RBAC) AddRole(roleID RoleID, perms []Permission) {
	r.roles[roleID] = &Role{
		ID:          roleID,
		Permissions: perms,
	}
}

// 添加用户
func (r *RBAC) AddUser(userID UserID) {
	r.users[userID] = &User{
		ID:    userID,
		Roles: []RoleID{},
	}
}

// 为用户分配角色
func (r *RBAC) AssignRole(userID UserID, roleID RoleID) {
	user, ok := r.users[userID]
	if !ok {
		return
	}
	user.Roles = append(user.Roles, roleID)
}

// 权限检查核心逻辑
func (r *RBAC) HasPermission(userID UserID, perm Permission) bool {
	user, ok := r.users[userID]
	if !ok {
		return false
	}

	for _, roleID := range user.Roles {
		role, ok := r.roles[roleID]
		if !ok {
			continue
		}
		for _, p := range role.Permissions {
			if p == perm {
				return true
			}
		}
	}
	return false
}
