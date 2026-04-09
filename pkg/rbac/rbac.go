package rbac

import (
	"context"
	"fmt"
	"os"
	"sync"

	"go_sample_code/pkg/logger"

	"github.com/casbin/casbin/v2"

	"go.uber.org/zap"
)

// RBACService RBAC 权限服务接口
type RBACService interface {
	// Enforce 检查权限
	Enforce(ctx context.Context, sub string, obj string, act string) (bool, error)
	// AddPolicy 添加策略
	AddPolicy(ctx context.Context, sub string, obj string, act string) error
	// RemovePolicy 移除策略
	RemovePolicy(ctx context.Context, sub string, obj string, act string) error
	// GetFilteredPolicy 获取过滤后的策略
	GetFilteredPolicy(ctx context.Context, fieldIndex int, fieldValues ...string) [][]string
	// AddRoleForUser 为用户分配角色
	AddRoleForUser(ctx context.Context, user string, role string) error
	// RemoveRoleForUser 移除用户的角色
	RemoveRoleForUser(ctx context.Context, user string, role string) error
	// GetRolesForUser 获取用户的所有角色
	GetRolesForUser(ctx context.Context, user string) ([]string, error)
	// GetUsersForRole 获取角色的所有用户
	GetUsersForRole(ctx context.Context, role string) ([]string, error)
	// DeleteRole 删除角色
	DeleteRole(ctx context.Context, role string) error
	// DeleteUser 删除用户的所有角色关联
	DeleteUser(ctx context.Context, user string) error
	// HasRoleForUser 检查用户是否拥有角色
	HasRoleForUser(ctx context.Context, user string, role string) (bool, error)
	// AddPermissionForRole 为角色添加权限
	AddPermissionForRole(ctx context.Context, role string, obj string, act string) error
	// RemovePermissionForRole 移除角色的权限
	RemovePermissionForRole(ctx context.Context, role string, obj string, act string) error
	// GetPermissionsForRole 获取角色的所有权限
	GetPermissionsForRole(ctx context.Context, role string) [][]string
	// ReloadPolicy 重新加载策略
	ReloadPolicy(ctx context.Context) error
}

// rbacService RBAC 服务实现
type rbacService struct {
	enforcer *casbin.Enforcer
	config   *RBACConfig
	log      logger.Logger
	mu       sync.RWMutex
}

// NewRBACService 创建 RBAC 服务
func NewRBACService(cfg RBACConfig, log logger.Logger) (RBACService, error) {
	if !cfg.Enabled {
		log.Warn("RBAC is disabled")
		return &disabledRBACService{log: log}, nil
	}

	// 确保路径不为空
	modelPath := cfg.ModelPath
	policyPath := cfg.PolicyPath
	if modelPath == "" {
		modelPath = "pkg/rbac/model.conf"
	}
	if policyPath == "" {
		policyPath = "pkg/rbac/policy.csv"
	}

	// 验证文件是否存在
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("model file not found: %s", modelPath)
	}
	if _, err := os.Stat(policyPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("policy file not found: %s", policyPath)
	}

	// 新版本 casbin 直接支持从文件加载模型和策略
	enforcer, err := casbin.NewEnforcer(modelPath, policyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create enforcer: %w", err)
	}

	// 设置日志
	enforcer.EnableLog(true)

	cfgCopy := cfg
	service := &rbacService{
		enforcer: enforcer,
		config:   &cfgCopy,
		log:      log,
	}

	log.Info("RBAC service initialized",
		zap.String("model_path", modelPath),
		zap.String("policy_path", policyPath))

	return service, nil
}

// NewRBACServiceWithDefaults 使用默认配置创建 RBAC 服务
func NewRBACServiceWithDefaults(log logger.Logger) (RBACService, error) {
	return NewRBACService(DefaultConfig().RBAC, log)
}

// Enforce 检查权限
func (s *rbacService) Enforce(ctx context.Context, sub string, obj string, act string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result, err := s.enforcer.Enforce(sub, obj, act)
	if err != nil {
		s.log.ErrorCtx(ctx, "enforce failed",
			zap.String("subject", sub),
			zap.String("object", obj),
			zap.String("action", act),
			zap.Error(err))
		return false, err
	}

	s.log.DebugCtx(ctx, "enforce result",
		zap.String("subject", sub),
		zap.String("object", obj),
		zap.String("action", act),
		zap.Bool("result", result))

	return result, nil
}

// AddPolicy 添加策略
func (s *rbacService) AddPolicy(ctx context.Context, sub string, obj string, act string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.enforcer.AddPolicy(sub, obj, act)
	if err != nil {
		s.log.ErrorCtx(ctx, "add policy failed",
			zap.String("subject", sub),
			zap.String("object", obj),
			zap.String("action", act),
			zap.Error(err))
		return err
	}

	s.log.InfoCtx(ctx, "policy added",
		zap.String("subject", sub),
		zap.String("object", obj),
		zap.String("action", act))

	return nil
}

// RemovePolicy 移除策略
func (s *rbacService) RemovePolicy(ctx context.Context, sub string, obj string, act string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.enforcer.RemovePolicy(sub, obj, act)
	if err != nil {
		s.log.ErrorCtx(ctx, "remove policy failed",
			zap.String("subject", sub),
			zap.String("object", obj),
			zap.String("action", act),
			zap.Error(err))
		return err
	}

	s.log.InfoCtx(ctx, "policy removed",
		zap.String("subject", sub),
		zap.String("object", obj),
		zap.String("action", act))

	return nil
}

// GetFilteredPolicy 获取过滤后的策略
func (s *rbacService) GetFilteredPolicy(ctx context.Context, fieldIndex int, fieldValues ...string) [][]string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result, _ := s.enforcer.GetFilteredPolicy(fieldIndex, fieldValues...)
	return result
}

// AddRoleForUser 为用户分配角色
func (s *rbacService) AddRoleForUser(ctx context.Context, user string, role string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.enforcer.AddGroupingPolicy(user, role)
	if err != nil {
		s.log.ErrorCtx(ctx, "add role for user failed",
			zap.String("user", user),
			zap.String("role", role),
			zap.Error(err))
		return err
	}

	s.log.InfoCtx(ctx, "role assigned to user",
		zap.String("user", user),
		zap.String("role", role))

	return nil
}

// RemoveRoleForUser 移除用户的角色
func (s *rbacService) RemoveRoleForUser(ctx context.Context, user string, role string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.enforcer.RemoveGroupingPolicy(user, role)
	if err != nil {
		s.log.ErrorCtx(ctx, "remove role for user failed",
			zap.String("user", user),
			zap.String("role", role),
			zap.Error(err))
		return err
	}

	s.log.InfoCtx(ctx, "role removed from user",
		zap.String("user", user),
		zap.String("role", role))

	return nil
}

// GetRolesForUser 获取用户的所有角色
func (s *rbacService) GetRolesForUser(ctx context.Context, user string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	roles, err := s.enforcer.GetRolesForUser(user)
	if err != nil {
		s.log.ErrorCtx(ctx, "get roles for user failed",
			zap.String("user", user),
			zap.Error(err))
		return nil, err
	}

	return roles, nil
}

// GetUsersForRole 获取角色的所有用户
func (s *rbacService) GetUsersForRole(ctx context.Context, role string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users, err := s.enforcer.GetUsersForRole(role)
	if err != nil {
		s.log.ErrorCtx(ctx, "get users for role failed",
			zap.String("role", role),
			zap.Error(err))
		return nil, err
	}

	return users, nil
}

// DeleteRole 删除角色
func (s *rbacService) DeleteRole(ctx context.Context, role string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 删除角色的所有用户关联
	_, err := s.enforcer.RemoveFilteredGroupingPolicy(0, role)
	if err != nil {
		s.log.ErrorCtx(ctx, "delete role users failed",
			zap.String("role", role),
			zap.Error(err))
		return err
	}

	// 删除角色的所有权限策略
	_, err = s.enforcer.RemoveFilteredPolicy(0, role)
	if err != nil {
		s.log.ErrorCtx(ctx, "delete role policies failed",
			zap.String("role", role),
			zap.Error(err))
		return err
	}

	s.log.InfoCtx(ctx, "role deleted",
		zap.String("role", role))

	return nil
}

// DeleteUser 删除用户的所有角色关联
func (s *rbacService) DeleteUser(ctx context.Context, user string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.enforcer.RemoveFilteredGroupingPolicy(0, user)
	if err != nil {
		s.log.ErrorCtx(ctx, "delete user roles failed",
			zap.String("user", user),
			zap.Error(err))
		return err
	}

	s.log.InfoCtx(ctx, "user roles deleted",
		zap.String("user", user))

	return nil
}

// HasRoleForUser 检查用户是否拥有角色
func (s *rbacService) HasRoleForUser(ctx context.Context, user string, role string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.enforcer.HasGroupingPolicy(user, role)
}

// AddPermissionForRole 为角色添加权限
func (s *rbacService) AddPermissionForRole(ctx context.Context, role string, obj string, act string) error {
	return s.AddPolicy(ctx, role, obj, act)
}

// RemovePermissionForRole 移除角色的权限
func (s *rbacService) RemovePermissionForRole(ctx context.Context, role string, obj string, act string) error {
	return s.RemovePolicy(ctx, role, obj, act)
}

// GetPermissionsForRole 获取角色的所有权限
func (s *rbacService) GetPermissionsForRole(ctx context.Context, role string) [][]string {
	return s.GetFilteredPolicy(ctx, 0, role)
}

// ReloadPolicy 重新加载策略
func (s *rbacService) ReloadPolicy(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.enforcer.LoadPolicy(); err != nil {
		s.log.ErrorCtx(ctx, "reload policy failed", zap.Error(err))
		return err
	}

	s.log.InfoCtx(ctx, "policy reloaded")
	return nil
}

// disabledRBACService RBAC 未启用时的服务实现
type disabledRBACService struct {
	log logger.Logger
}

func (s *disabledRBACService) Enforce(ctx context.Context, sub string, obj string, act string) (bool, error) {
	s.log.WarnCtx(ctx, "RBAC is disabled, allowing all access")
	return true, nil
}

func (s *disabledRBACService) AddPolicy(ctx context.Context, sub string, obj string, act string) error {
	return ErrRBACNotEnabled
}

func (s *disabledRBACService) RemovePolicy(ctx context.Context, sub string, obj string, act string) error {
	return ErrRBACNotEnabled
}

func (s *disabledRBACService) GetFilteredPolicy(ctx context.Context, fieldIndex int, fieldValues ...string) [][]string {
	return nil
}

func (s *disabledRBACService) AddRoleForUser(ctx context.Context, user string, role string) error {
	return ErrRBACNotEnabled
}

func (s *disabledRBACService) RemoveRoleForUser(ctx context.Context, user string, role string) error {
	return ErrRBACNotEnabled
}

func (s *disabledRBACService) GetRolesForUser(ctx context.Context, user string) ([]string, error) {
	return nil, ErrRBACNotEnabled
}

func (s *disabledRBACService) GetUsersForRole(ctx context.Context, role string) ([]string, error) {
	return nil, ErrRBACNotEnabled
}

func (s *disabledRBACService) DeleteRole(ctx context.Context, role string) error {
	return ErrRBACNotEnabled
}

func (s *disabledRBACService) DeleteUser(ctx context.Context, user string) error {
	return ErrRBACNotEnabled
}

func (s *disabledRBACService) HasRoleForUser(ctx context.Context, user string, role string) (bool, error) {
	return false, ErrRBACNotEnabled
}

func (s *disabledRBACService) AddPermissionForRole(ctx context.Context, role string, obj string, act string) error {
	return ErrRBACNotEnabled
}

func (s *disabledRBACService) RemovePermissionForRole(ctx context.Context, role string, obj string, act string) error {
	return ErrRBACNotEnabled
}

func (s *disabledRBACService) GetPermissionsForRole(ctx context.Context, role string) [][]string {
	return nil
}

func (s *disabledRBACService) ReloadPolicy(ctx context.Context) error {
	return ErrRBACNotEnabled
}
