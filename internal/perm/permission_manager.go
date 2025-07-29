package perm

//
//import (
//	"context"
//	"errors"
//	"fmt"
//	"qwqserver/internal/model"
//	"qwqserver/internal/repository"
//)
//
//var (
//	ErrPermissionNotFound      = errors.New("permission not found")
//	ErrRoleNotFound            = errors.New("role not found")
//	ErrUserNotFound            = errors.New("user not found")
//	ErrPermissionAlreadyExists = errors.New("permission already exists")
//	ErrRoleAlreadyExists       = errors.New("role already exists")
//	ErrPermissionNotAssigned   = errors.New("permission not assigned to role")
//	ErrRoleNotAssigned         = errors.New("role not assigned to user")
//	ErrInvalidInput            = errors.New("invalid input")
//)
//
//// PermissionManager 权限管理服务
//type PermissionManager struct {
//	permRepo repository.PermissionRepository
//}
//
//func NewPermissionManager(permRepo repository.PermissionRepository) *PermissionManager {
//	return &PermissionManager{permRepo: permRepo}
//}
//
//// Initialize 初始化默认权限和角色
//func (m *PermissionManager) Initialize(ctx context.Context) error {
//	// 创建默认权限
//	defaultPermissions := []struct {
//		key         model.PermissionKey
//		description string
//	}{
//		{"user:create", "创建用户"},
//		{"user:read", "查看用户信息"},
//		{"user:update", "更新用户信息"},
//		{"user:delete", "删除用户"},
//		{"role:manage", "管理角色"},
//		{"permission:manage", "管理权限"},
//	}
//
//	for _, p := range defaultPermissions {
//		_, err := m.CreatePermission(ctx, p.key, p.description)
//		if err != nil && !errors.Is(err, ErrPermissionAlreadyExists) {
//			return fmt.Errorf("failed to create default permission: %w", err)
//		}
//	}
//
//	// 创建管理员角色
//	adminRole, err := m.CreateRole(ctx, "admin", "系统管理员")
//	if err != nil && !errors.Is(err, ErrRoleAlreadyExists) {
//		return fmt.Errorf("failed to create admin role: %w", err)
//	}
//
//	// 为管理员角色分配所有权限
//	if adminRole != nil {
//		allPerms, err := m.permRepo.ListPermissions(ctx)
//		if err != nil {
//			return fmt.Errorf("failed to list permissions: %w", err)
//		}
//
//		for _, perm := range allPerms {
//			err := m.AssignPermissionToRole(ctx, adminRole.ID, perm.ID)
//			if err != nil {
//				return fmt.Errorf("failed to assign permission to admin role: %w", err)
//			}
//		}
//	}
//
//	return nil
//}
//
//// CreatePermission 创建新权限
//func (m *PermissionManager) CreatePermission(
//	ctx context.Context,
//	key model.PermissionKey,
//	description string,
//) (*model.Permission, error) {
//	if key == "" {
//		return nil, ErrInvalidInput
//	}
//
//	// 检查权限是否已存在
//	_, err := m.permRepo.GetPermissionByKey(ctx, key)
//	if err == nil {
//		return nil, ErrPermissionAlreadyExists
//	} else if !errors.Is(err, repository.ErrRecordNotFound) {
//		return nil, err
//	}
//
//	perm := &model.Permission{
//		Key:         key,
//		Description: description,
//	}
//
//	if err := m.permRepo.CreatePermission(ctx, perm); err != nil {
//		return nil, fmt.Errorf("failed to create permission: %w", err)
//	}
//
//	return perm, nil
//}
//
//// CreateRole 创建新角色
//func (m *PermissionManager) CreateRole(
//	ctx context.Context,
//	name, description string,
//) (*model.Role, error) {
//	if name == "" {
//		return nil, ErrInvalidInput
//	}
//
//	// 检查角色是否已存在
//	_, err := m.permRepo.GetRoleByName(ctx, name)
//	if err == nil {
//		return nil, ErrRoleAlreadyExists
//	} else if !errors.Is(err, repository.ErrRecordNotFound) {
//		return nil, err
//	}
//
//	role := &model.Role{
//		Name:        name,
//		Description: description,
//	}
//
//	if err := m.permRepo.CreateRole(ctx, role); err != nil {
//		return nil, fmt.Errorf("failed to create role: %w", err)
//	}
//
//	return role, nil
//}
//
//// AssignPermissionToRole 为角色分配权限
//func (m *PermissionManager) AssignPermissionToRole(
//	ctx context.Context,
//	roleID, permissionID uint,
//) error {
//	if roleID == 0 || permissionID == 0 {
//		return ErrInvalidInput
//	}
//
//	// 检查角色是否存在
//	role, err := m.permRepo.GetRoleByID(ctx, roleID)
//	if err != nil {
//		if errors.Is(err, repository.ErrRecordNotFound) {
//			return ErrRoleNotFound
//		}
//		return err
//	}
//
//	// 检查权限是否存在
//	perm, err := m.permRepo.GetPermissionByID(ctx, permissionID)
//	if err != nil {
//		if errors.Is(err, repository.ErrRecordNotFound) {
//			return ErrPermissionNotFound
//		}
//		return err
//	}
//
//	// 检查是否已分配
//	perms, err := m.permRepo.GetRolePermissions(ctx, role.ID)
//	if err != nil {
//		return err
//	}
//
//	for _, p := range perms {
//		if p.ID == perm.ID {
//			return nil // 已存在，无需操作
//		}
//	}
//
//	return m.permRepo.AddPermissionToRole(ctx, role.ID, perm.ID)
//}
//
//// AssignRoleToUser 为用户分配角色
//func (m *PermissionManager) AssignRoleToUser(
//	ctx context.Context,
//	userID, roleID uint,
//) error {
//	if userID == 0 || roleID == 0 {
//		return ErrInvalidInput
//	}
//
//	// 检查角色是否存在
//	role, err := m.permRepo.GetRoleByID(ctx, roleID)
//	if err != nil {
//		if errors.Is(err, repository.ErrRecordNotFound) {
//			return ErrRoleNotFound
//		}
//		return err
//	}
//
//	// 检查用户是否存在
//	// 这里假设您有一个用户仓库
//	// user, err := userRepo.FindByID(ctx, userID)
//	// if err != nil {
//	//     if errors.Is(err, repository.ErrRecordNotFound) {
//	//         return ErrUserNotFound
//	//     }
//	//     return err
//	// }
//
//	// 检查是否已分配
//	roles, err := m.permRepo.GetUserRoles(ctx, userID)
//	if err != nil {
//		return err
//	}
//
//	for _, r := range roles {
//		if r.ID == role.ID {
//			return nil // 已存在，无需操作
//		}
//	}
//
//	return m.permRepo.AssignRoleToUser(ctx, userID, role.ID)
//}
//
//// CheckPermission 检查用户是否拥有指定权限
//func (m *PermissionManager) CheckPermission(
//	ctx context.Context,
//	userID uint,
//	permissionKey model.PermissionKey,
//) (bool, error) {
//	if userID == 0 || permissionKey == "" {
//		return false, ErrInvalidInput
//	}
//
//	// 获取用户所有权限
//	permissions, err := m.permRepo.GetUserPermissions(ctx, userID)
//	if err != nil {
//		return false, err
//	}
//
//	for _, perm := range permissions {
//		if perm.Key == permissionKey {
//			return true, nil
//		}
//	}
//
//	return false, nil
//}
//
//// GetUserPermissions 获取用户的所有权限
//func (m *PermissionManager) GetUserPermissions(
//	ctx context.Context,
//	userID uint,
//) ([]*model.Permission, error) {
//	if userID == 0 {
//		return nil, ErrInvalidInput
//	}
//	return m.permRepo.GetUserPermissions(ctx, userID)
//}
//
//// GetUserRoles 获取用户的所有角色
//func (m *PermissionManager) GetUserRoles(
//	ctx context.Context,
//	userID uint,
//) ([]*model.Role, error) {
//	if userID == 0 {
//		return nil, ErrInvalidInput
//	}
//	return m.permRepo.GetUserRoles(ctx, userID)
//}
//
//// RemovePermissionFromRole 从角色移除权限
//func (m *PermissionManager) RemovePermissionFromRole(
//	ctx context.Context,
//	roleID, permissionID uint,
//) error {
//	if roleID == 0 || permissionID == 0 {
//		return ErrInvalidInput
//	}
//
//	// 检查权限是否已分配
//	perms, err := m.permRepo.GetRolePermissions(ctx, roleID)
//	if err != nil {
//		return err
//	}
//
//	found := false
//	for _, perm := range perms {
//		if perm.ID == permissionID {
//			found = true
//			break
//		}
//	}
//
//	if !found {
//		return ErrPermissionNotAssigned
//	}
//
//	return m.permRepo.RemovePermissionFromRole(ctx, roleID, permissionID)
//}
//
//// RemoveRoleFromUser 移除用户的角色
//func (m *PermissionManager) RemoveRoleFromUser(
//	ctx context.Context,
//	userID, roleID uint,
//) error {
//	if userID == 0 || roleID == 0 {
//		return ErrInvalidInput
//	}
//
//	// 检查角色是否已分配
//	roles, err := m.permRepo.GetUserRoles(ctx, userID)
//	if err != nil {
//		return err
//	}
//
//	found := false
//	for _, role := range roles {
//		if role.ID == roleID {
//			found = true
//			break
//		}
//	}
//
//	if !found {
//		return ErrRoleNotAssigned
//	}
//
//	return m.permRepo.RemoveRoleFromUser(ctx, userID, roleID)
//}
//
//// DeletePermission 删除权限
//func (m *PermissionManager) DeletePermission(
//	ctx context.Context,
//	permissionID uint,
//) error {
//	if permissionID == 0 {
//		return ErrInvalidInput
//	}
//
//	// 检查权限是否被角色使用
//	roles, err := m.permRepo.ListRoles(ctx)
//	if err != nil {
//		return err
//	}
//
//	for _, role := range roles {
//		for _, perm := range role.Permissions {
//			if perm.ID == permissionID {
//				return fmt.Errorf("cannot delete permission: it is assigned to role '%s'", role.Name)
//			}
//		}
//	}
//
//	return m.permRepo.DeletePermission(ctx, permissionID)
//}
//
//// DeleteRole 删除角色
//func (m *PermissionManager) DeleteRole(
//	ctx context.Context,
//	roleID uint,
//) error {
//	if roleID == 0 {
//		return ErrInvalidInput
//	}
//
//	// 检查角色是否被用户使用
//	users, err := m.getUsersWithRole(ctx, roleID)
//	if err != nil {
//		return err
//	}
//
//	if len(users) > 0 {
//		return fmt.Errorf("cannot delete role: it is assigned to %d users", len(users))
//	}
//
//	return m.permRepo.DeleteRole(ctx, roleID)
//}
//
//// getUsersWithRole 获取拥有指定角色的用户
//func (m *PermissionManager) getUsersWithRole(
//	ctx context.Context,
//	roleID uint,
//) ([]uint, error) {
//	// 在实际项目中，您可能需要实现此方法
//	// 这里返回一个空切片作为示例
//	return []uint{}, nil
//}
//
//// HasPermissionForResource 检查用户是否对特定资源拥有权限
//func (m *PermissionManager) HasPermissionForResource(
//	ctx context.Context,
//	userID uint,
//	action model.PermissionKey,
//	resourceType string,
//	resourceID uint,
//) (bool, error) {
//	// 构造资源特定权限键
//	permKey := model.PermissionKey(fmt.Sprintf("%s:%s:%d", resourceType, action, resourceID))
//
//	// 首先检查特定资源权限
//	hasPerm, err := m.CheckPermission(ctx, userID, permKey)
//	if err != nil {
//		return false, err
//	}
//	if hasPerm {
//		return true, nil
//	}
//
//	// 检查全局权限
//	globalPermKey := model.PermissionKey(fmt.Sprintf("%s:%s", resourceType, action))
//	return m.CheckPermission(ctx, userID, globalPermKey)
//}
