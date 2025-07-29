package rbac

import (
	"errors"
	"sync"
)

// 类型定义
type UserID string
type RoleID string
type Permission string

// 错误定义
var (
	ErrUserNotFound          = errors.New("user not found")
	ErrRoleNotFound          = errors.New("role not found")
	ErrPermissionNotFound    = errors.New("permission not found")
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrRoleAlreadyExists     = errors.New("role already exists")
	ErrPermissionExists      = errors.New("permission already exists")
	ErrPermissionNotAssigned = errors.New("permission not assigned to role")
	ErrRoleNotAssigned       = errors.New("role not assigned to user")
)

// RBAC 核心结构
type RBAC struct {
	users       map[UserID]*User
	roles       map[RoleID]*Role
	permissions map[Permission]struct{}
	mu          sync.RWMutex
}

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

// 初始化和构造函数
func NewRBAC() *RBAC {
	
	return &RBAC{
		users:       make(map[UserID]*User),
		roles:       make(map[RoleID]*Role),
		permissions: make(map[Permission]struct{}),
	}
}

// 添加权限
func (r *RBAC) AddPermission(p Permission) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.permissions[p] = struct{}{}
}

// 添加角色
func (r *RBAC) AddRole(roleID RoleID, perms []Permission) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.roles[roleID]; exists {
		return ErrRoleAlreadyExists
	}

	// 验证权限是否存在
	for _, p := range perms {
		if _, exists := r.permissions[p]; !exists {
			return ErrPermissionNotFound
		}
	}

	r.roles[roleID] = &Role{
		ID:          roleID,
		Permissions: perms,
	}
	return nil
}

// 添加用户
func (r *RBAC) AddUser(userID UserID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[userID]; exists {
		return ErrUserAlreadyExists
	}

	r.users[userID] = &User{
		ID:    userID,
		Roles: []RoleID{},
	}
	return nil
}

// 为用户分配角色
func (r *RBAC) AssignRole(userID UserID, roleID RoleID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[userID]
	if !ok {
		return ErrUserNotFound
	}

	if _, exists := r.roles[roleID]; !exists {
		return ErrRoleNotFound
	}

	// 检查是否已分配
	for _, rID := range user.Roles {
		if rID == roleID {
			return nil // 已存在，无需添加
		}
	}

	user.Roles = append(user.Roles, roleID)
	return nil
}

// 为角色添加权限
func (r *RBAC) AddPermissionToRole(roleID RoleID, perm Permission) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	role, ok := r.roles[roleID]
	if !ok {
		return ErrRoleNotFound
	}

	if _, exists := r.permissions[perm]; !exists {
		return ErrPermissionNotFound
	}

	// 检查是否已存在
	for _, p := range role.Permissions {
		if p == perm {
			return ErrPermissionExists
		}
	}

	role.Permissions = append(role.Permissions, perm)
	return nil
}

// 权限检查核心逻辑
func (r *RBAC) HasPermission(userID UserID, perm Permission) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[userID]
	if !ok {
		return false, ErrUserNotFound
	}

	for _, roleID := range user.Roles {
		role, ok := r.roles[roleID]
		if !ok {
			continue
		}
		for _, p := range role.Permissions {
			if p == perm {
				return true, nil
			}
		}
	}
	return false, nil
}

// 移除用户的角色
func (r *RBAC) RemoveRoleFromUser(userID UserID, roleID RoleID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[userID]
	if !ok {
		return ErrUserNotFound
	}

	found := false
	newRoles := make([]RoleID, 0, len(user.Roles))
	for _, rID := range user.Roles {
		if rID == roleID {
			found = true
			continue
		}
		newRoles = append(newRoles, rID)
	}

	if !found {
		return ErrRoleNotAssigned
	}

	user.Roles = newRoles
	return nil
}

// 从角色移除权限
func (r *RBAC) RemovePermissionFromRole(roleID RoleID, perm Permission) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	role, ok := r.roles[roleID]
	if !ok {
		return ErrRoleNotFound
	}

	found := false
	newPerms := make([]Permission, 0, len(role.Permissions))
	for _, p := range role.Permissions {
		if p == perm {
			found = true
			continue
		}
		newPerms = append(newPerms, p)
	}

	if !found {
		return ErrPermissionNotAssigned
	}

	role.Permissions = newPerms
	return nil
}

// 获取用户的所有角色
func (r *RBAC) GetUserRoles(userID UserID) ([]RoleID, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[userID]
	if !ok {
		return nil, ErrUserNotFound
	}

	// 返回副本
	roles := make([]RoleID, len(user.Roles))
	copy(roles, user.Roles)
	return roles, nil
}

// 获取角色的所有权限
func (r *RBAC) GetRolePermissions(roleID RoleID) ([]Permission, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	role, ok := r.roles[roleID]
	if !ok {
		return nil, ErrRoleNotFound
	}

	// 返回副本
	perms := make([]Permission, len(role.Permissions))
	copy(perms, role.Permissions)
	return perms, nil
}

// 获取所有权限
func (r *RBAC) GetAllPermissions() []Permission {
	r.mu.RLock()
	defer r.mu.RUnlock()

	perms := make([]Permission, 0, len(r.permissions))
	for p := range r.permissions {
		perms = append(perms, p)
	}
	return perms
}

// 删除用户
func (r *RBAC) DeleteUser(userID UserID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[userID]; !exists {
		return ErrUserNotFound
	}

	delete(r.users, userID)
	return nil
}

// 删除角色
func (r *RBAC) DeleteRole(roleID RoleID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.roles[roleID]; !exists {
		return ErrRoleNotFound
	}

	// 从所有用户中移除该角色
	for _, user := range r.users {
		newRoles := make([]RoleID, 0, len(user.Roles))
		for _, rID := range user.Roles {
			if rID != roleID {
				newRoles = append(newRoles, rID)
			}
		}
		user.Roles = newRoles
	}

	delete(r.roles, roleID)
	return nil
}

// 删除权限
func (r *RBAC) DeletePermission(perm Permission) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.permissions[perm]; !exists {
		return ErrPermissionNotFound
	}

	// 从所有角色中移除该权限
	for _, role := range r.roles {
		newPerms := make([]Permission, 0, len(role.Permissions))
		for _, p := range role.Permissions {
			if p != perm {
				newPerms = append(newPerms, p)
			}
		}
		role.Permissions = newPerms
	}

	delete(r.permissions, perm)
	return nil
}
