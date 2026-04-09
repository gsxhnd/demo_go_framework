package rbac

import (
	"context"

	"go_sample_code/pkg/jwx"
)

// PermissionService 权限服务接口
// 整合 RBAC 和 ABAC 功能，提供统一的权限管理接口
type PermissionService interface {
	// CheckPermission 检查用户是否有权限访问指定资源
	CheckPermission(ctx context.Context, claims *jwx.SelfTokenClaims, obj string, act string) (bool, error)

	// AssignRole 为用户分配角色
	AssignRole(ctx context.Context, userID string, role string) error

	// RevokeRole 撤销用户的角色
	RevokeRole(ctx context.Context, userID string, role string) error

	// GetUserRoles 获取用户的所有角色
	GetUserRoles(ctx context.Context, userID string) ([]string, error)

	// HasRole 检查用户是否拥有指定角色
	HasRole(ctx context.Context, userID string, role string) (bool, error)

	// AddPermission 添加权限
	AddPermission(ctx context.Context, role string, obj string, act string) error

	// RemovePermission 移除权限
	RemovePermission(ctx context.Context, role string, obj string, act string) error

	// GetRolePermissions 获取角色的所有权限
	GetRolePermissions(ctx context.Context, role string) [][]string

	// CreateRole 创建新角色
	CreateRole(ctx context.Context, role string, description string) error

	// DeleteRole 删除角色
	DeleteRole(ctx context.Context, role string) error

	// CheckABAC 检查 ABAC 权限
	CheckABAC(ctx context.Context, claims *jwx.SelfTokenClaims, req *ABACRequest) (bool, error)
}

// permissionService 权限服务实现
type permissionService struct {
	rbac RBACService
	abac ABACService
}

// NewPermissionService 创建权限服务
func NewPermissionService(rbac RBACService, abac ABACService) PermissionService {
	return &permissionService{
		rbac: rbac,
		abac: abac,
	}
}

// CheckPermission 检查用户是否有权限访问指定资源
func (s *permissionService) CheckPermission(ctx context.Context, claims *jwx.SelfTokenClaims, obj string, act string) (bool, error) {
	if claims == nil {
		return false, ErrSubjectInvalid
	}

	// 先检查 RBAC
	allowed, err := s.rbac.Enforce(ctx, claims.Role, obj, act)
	if err != nil {
		return false, err
	}

	if allowed {
		return true, nil
	}

	// RBAC 未通过，检查 ABAC
	if s.abac != nil {
		req := &ABACRequest{
			Subject: SubjectAttribute{
				ID:   claims.Uid,
				Role: claims.Role,
			},
			Object: ObjectAttribute{
				Path: obj,
			},
			Action: ActionAttribute{
				Name: act,
			},
		}
		return s.abac.Enforce(ctx, req)
	}

	return false, nil
}

// AssignRole 为用户分配角色
func (s *permissionService) AssignRole(ctx context.Context, userID string, role string) error {
	return s.rbac.AddRoleForUser(ctx, userID, role)
}

// RevokeRole 撤销用户的角色
func (s *permissionService) RevokeRole(ctx context.Context, userID string, role string) error {
	return s.rbac.RemoveRoleForUser(ctx, userID, role)
}

// GetUserRoles 获取用户的所有角色
func (s *permissionService) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	return s.rbac.GetRolesForUser(ctx, userID)
}

// HasRole 检查用户是否拥有指定角色
func (s *permissionService) HasRole(ctx context.Context, userID string, role string) (bool, error) {
	return s.rbac.HasRoleForUser(ctx, userID, role)
}

// AddPermission 添加权限
func (s *permissionService) AddPermission(ctx context.Context, role string, obj string, act string) error {
	return s.rbac.AddPermissionForRole(ctx, role, obj, act)
}

// RemovePermission 移除权限
func (s *permissionService) RemovePermission(ctx context.Context, role string, obj string, act string) error {
	return s.rbac.RemovePermissionForRole(ctx, role, obj, act)
}

// GetRolePermissions 获取角色的所有权限
func (s *permissionService) GetRolePermissions(ctx context.Context, role string) [][]string {
	return s.rbac.GetPermissionsForRole(ctx, role)
}

// CreateRole 创建新角色
func (s *permissionService) CreateRole(ctx context.Context, role string, description string) error {
	// 角色创建实际上是在策略文件中添加
	// 这里可以添加初始化逻辑
	return nil
}

// DeleteRole 删除角色
func (s *permissionService) DeleteRole(ctx context.Context, role string) error {
	return s.rbac.DeleteRole(ctx, role)
}

// CheckABAC 检查 ABAC 权限
func (s *permissionService) CheckABAC(ctx context.Context, claims *jwx.SelfTokenClaims, req *ABACRequest) (bool, error) {
	if s.abac == nil {
		return true, nil
	}
	return s.abac.Enforce(ctx, req)
}
