package repository

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"qwqserver/internal/model"
)

// PermissionRepository 权限仓库接口
type PermissionRepository interface {
	CreatePermission(ctx context.Context, permission *model.Permission) error
	GetPermissionByKey(ctx context.Context, key model.PermissionKey) (*model.Permission, error)
	DeletePermission(ctx context.Context, id uint) error
	ListPermissions(ctx context.Context) ([]*model.Permission, error)

	CreateRole(ctx context.Context, role *model.Role) error
	GetRoleByName(ctx context.Context, name string) (*model.Role, error)
	DeleteRole(ctx context.Context, id uint) error
	ListRoles(ctx context.Context) ([]*model.Role, error)

	AddPermissionToRole(ctx context.Context, roleID, permissionID uint) error
	RemovePermissionFromRole(ctx context.Context, roleID, permissionID uint) error
	GetRolePermissions(ctx context.Context, roleID uint) ([]*model.Permission, error)

	AssignRoleToUser(ctx context.Context, userID, roleID uint) error
	RemoveRoleFromUser(ctx context.Context, userID, roleID uint) error
	GetUserRoles(ctx context.Context, userID uint) ([]*model.Role, error)
	GetUserPermissions(ctx context.Context, userID uint) ([]*model.Permission, error)
}

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

// CreatePermission 创建权限
func (r *permissionRepository) CreatePermission(ctx context.Context, permission *model.Permission) error {
	if permission.Key == "" {
		return errors.New("permission key cannot be empty")
	}

	return r.db.WithContext(ctx).Create(permission).Error
}

// GetPermissionByKey 根据Key获取权限
func (r *permissionRepository) GetPermissionByKey(ctx context.Context, key model.PermissionKey) (*model.Permission, error) {
	if key == "" {
		return nil, errors.New("permission key cannot be empty")
	}

	var permission model.Permission
	err := r.db.WithContext(ctx).Where("key = ?", key).First(&permission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("permission with key '%s' not found", key)
		}
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}
	return &permission, nil
}

// DeletePermission 删除权限
func (r *permissionRepository) DeletePermission(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("permission ID cannot be zero")
	}

	// 先删除关联关系
	if err := r.db.WithContext(ctx).Where("permission_id = ?", id).Delete(&model.RolePermission{}).Error; err != nil {
		return fmt.Errorf("failed to remove permission from roles: %w", err)
	}

	return r.db.WithContext(ctx).Delete(&model.Permission{}, id).Error
}

// ListPermissions 列出所有权限
func (r *permissionRepository) ListPermissions(ctx context.Context) ([]*model.Permission, error) {
	var permissions []*model.Permission
	err := r.db.WithContext(ctx).Find(&permissions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}
	return permissions, nil
}

// CreateRole 创建角色
func (r *permissionRepository) CreateRole(ctx context.Context, role *model.Role) error {
	if role.Name == "" {
		return errors.New("role name cannot be empty")
	}

	return r.db.WithContext(ctx).Create(role).Error
}

// GetRoleByName 根据名称获取角色
func (r *permissionRepository) GetRoleByName(ctx context.Context, name string) (*model.Role, error) {
	if name == "" {
		return nil, errors.New("role name cannot be empty")
	}

	var role model.Role
	err := r.db.WithContext(ctx).Where("name = ?", name).Preload("Permissions").First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("role '%s' not found", name)
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	return &role, nil
}

// DeleteRole 删除角色
func (r *permissionRepository) DeleteRole(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("role ID cannot be zero")
	}

	// 先删除关联关系
	if err := r.db.WithContext(ctx).Where("role_id = ?", id).Delete(&model.RolePermission{}).Error; err != nil {
		return fmt.Errorf("failed to remove permissions from role: %w", err)
	}

	if err := r.db.WithContext(ctx).Where("role_id = ?", id).Delete(&model.UserRole{}).Error; err != nil {
		return fmt.Errorf("failed to remove role from users: %w", err)
	}

	return r.db.WithContext(ctx).Delete(&model.Role{}, id).Error
}

// ListRoles 列出所有角色
func (r *permissionRepository) ListRoles(ctx context.Context) ([]*model.Role, error) {
	var roles []*model.Role
	err := r.db.WithContext(ctx).Preload("Permissions").Find(&roles).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}
	return roles, nil
}

// AddPermissionToRole 添加权限到角色
func (r *permissionRepository) AddPermissionToRole(ctx context.Context, roleID, permissionID uint) error {
	if roleID == 0 {
		return errors.New("role ID cannot be zero")
	}
	if permissionID == 0 {
		return errors.New("permission ID cannot be zero")
	}

	// 检查是否已存在
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.RolePermission{}).
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check existing permission: %w", err)
	}

	if count > 0 {
		return nil // 已存在，无需添加
	}

	return r.db.WithContext(ctx).Create(&model.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
	}).Error
}

// RemovePermissionFromRole 从角色移除权限
func (r *permissionRepository) RemovePermissionFromRole(ctx context.Context, roleID, permissionID uint) error {
	if roleID == 0 {
		return errors.New("role ID cannot be zero")
	}
	if permissionID == 0 {
		return errors.New("permission ID cannot be zero")
	}

	return r.db.WithContext(ctx).Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Delete(&model.RolePermission{}).Error
}

// GetRolePermissions 获取角色的权限
func (r *permissionRepository) GetRolePermissions(ctx context.Context, roleID uint) ([]*model.Permission, error) {
	if roleID == 0 {
		return nil, errors.New("role ID cannot be zero")
	}

	var permissions []*model.Permission
	err := r.db.WithContext(ctx).Model(&model.Permission{}).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}
	return permissions, nil
}

// AssignRoleToUser 为用户分配角色
func (r *permissionRepository) AssignRoleToUser(ctx context.Context, userID, roleID uint) error {
	if userID == 0 {
		return errors.New("user ID cannot be zero")
	}
	if roleID == 0 {
		return errors.New("role ID cannot be zero")
	}

	// 检查是否已存在
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.UserRole{}).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check existing role: %w", err)
	}

	if count > 0 {
		return nil // 已存在，无需添加
	}

	return r.db.WithContext(ctx).Create(&model.UserRole{
		UserID: userID,
		RoleID: roleID,
	}).Error
}

// RemoveRoleFromUser 移除用户的角色
func (r *permissionRepository) RemoveRoleFromUser(ctx context.Context, userID, roleID uint) error {
	if userID == 0 {
		return errors.New("user ID cannot be zero")
	}
	if roleID == 0 {
		return errors.New("role ID cannot be zero")
	}

	return r.db.WithContext(ctx).Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&model.UserRole{}).Error
}

// GetUserRoles 获取用户的角色
func (r *permissionRepository) GetUserRoles(ctx context.Context, userID uint) ([]*model.Role, error) {
	if userID == 0 {
		return nil, errors.New("user ID cannot be zero")
	}

	var roles []*model.Role
	err := r.db.WithContext(ctx).Model(&model.Role{}).
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	return roles, nil
}

// GetUserPermissions 获取用户的所有权限
func (r *permissionRepository) GetUserPermissions(ctx context.Context, userID uint) ([]*model.Permission, error) {
	if userID == 0 {
		return nil, errors.New("user ID cannot be zero")
	}

	var permissions []*model.Permission
	err := r.db.WithContext(ctx).Model(&model.Permission{}).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", userID).
		Group("permissions.id").
		Find(&permissions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
	}
	return permissions, nil
}
