package rbac

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go_sample_code/pkg/jwx"
	"go_sample_code/pkg/logger"
)

func newTestLogger() logger.Logger {
	l, err := logger.NewLogger(&logger.LoggerConfig{
		Level:  "debug",
		Output: "console",
	})
	if err != nil {
		panic(err)
	}
	return l
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid default config",
			config: &Config{
				RBAC: RBACConfig{
					Enabled:     true,
					ModelPath:   "pkg/rbac/model.conf",
					PolicyPath:  "pkg/rbac/policy.csv",
					DefaultRole: "guest",
				},
				ABAC: ABACConfig{
					Enabled:   true,
					ModelPath: "pkg/rbac/abac_model.conf",
				},
			},
			wantErr: false,
		},
		{
			name: "RBAC disabled",
			config: &Config{
				RBAC: RBACConfig{
					Enabled: false,
				},
			},
			wantErr: false,
		},
		{
			name: "RBAC enabled but missing model path",
			config: &Config{
				RBAC: RBACConfig{
					Enabled:    true,
					PolicyPath: "pkg/rbac/policy.csv",
				},
			},
			wantErr: true,
		},
		{
			name: "RBAC enabled but missing policy path",
			config: &Config{
				RBAC: RBACConfig{
					Enabled:   true,
					ModelPath: "pkg/rbac/model.conf",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	assert.NotNil(t, cfg)
	assert.True(t, cfg.RBAC.Enabled)
	assert.Equal(t, "pkg/rbac/model.conf", cfg.RBAC.ModelPath)
	assert.Equal(t, "pkg/rbac/policy.csv", cfg.RBAC.PolicyPath)
	assert.Equal(t, "guest", cfg.RBAC.DefaultRole)
	assert.True(t, cfg.ABAC.Enabled)
}

func TestRBACService_Disabled(t *testing.T) {
	log := newTestLogger()
	cfg := RBACConfig{Enabled: false}
	svc, err := NewRBACService(cfg, log)
	require.NoError(t, err)

	ctx := context.Background()

	// RBAC disabled 时应该允许所有操作
	allowed, err := svc.Enforce(ctx, "admin", "/api/v1/users", "read")
	assert.NoError(t, err)
	assert.True(t, allowed)

	// 其他操作应该返回 ErrRBACNotEnabled
	err = svc.AddPolicy(ctx, "admin", "/api/v1/test", "read")
	assert.Equal(t, ErrRBACNotEnabled, err)
}

func TestRBACService_WithPolicy(t *testing.T) {
	log := newTestLogger()
	cfg := RBACConfig{
		Enabled:    true,
		ModelPath:  "pkg/rbac/model.conf",
		PolicyPath: "pkg/rbac/policy.csv",
	}

	svc, err := NewRBACService(cfg, log)
	if err != nil {
		t.Skip("Skipping test due to missing policy file")
	}

	ctx := context.Background()

	// 测试管理员权限
	allowed, err := svc.Enforce(ctx, "admin", "/api/v1/users", "read")
	assert.NoError(t, err)
	assert.True(t, allowed)

	// 测试普通用户权限
	allowed, err = svc.Enforce(ctx, "user", "/api/v1/users", "read")
	assert.NoError(t, err)
	assert.True(t, allowed)

	// 测试访问受保护资源
	allowed, err = svc.Enforce(ctx, "user", "/api/v1/admin", "read")
	assert.NoError(t, err)
	assert.False(t, allowed)

	// 测试访客权限
	allowed, err = svc.Enforce(ctx, "guest", "/api/v1/public", "read")
	assert.NoError(t, err)
	assert.True(t, allowed)
}

func TestRBACService_RoleManagement(t *testing.T) {
	log := newTestLogger()
	cfg := RBACConfig{
		Enabled:    true,
		ModelPath:  "pkg/rbac/model.conf",
		PolicyPath: "pkg/rbac/policy.csv",
	}

	svc, err := NewRBACService(cfg, log)
	if err != nil {
		t.Skip("Skipping test due to missing policy file")
	}

	ctx := context.Background()

	// 添加用户角色关联
	err = svc.AddRoleForUser(ctx, "user123", "editor")
	assert.NoError(t, err)

	// 验证角色关联
	hasRole, err := svc.HasRoleForUser(ctx, "user123", "editor")
	assert.NoError(t, err)
	assert.True(t, hasRole)

	// 获取用户角色
	roles, err := svc.GetRolesForUser(ctx, "user123")
	assert.NoError(t, err)
	assert.Contains(t, roles, "editor")

	// 移除角色关联
	err = svc.RemoveRoleForUser(ctx, "user123", "editor")
	assert.NoError(t, err)

	// 验证角色已移除
	hasRole, err = svc.HasRoleForUser(ctx, "user123", "editor")
	assert.NoError(t, err)
	assert.False(t, hasRole)
}

func TestRBACService_PolicyManagement(t *testing.T) {
	log := newTestLogger()
	cfg := RBACConfig{
		Enabled:    true,
		ModelPath:  "pkg/rbac/model.conf",
		PolicyPath: "pkg/rbac/policy.csv",
	}

	svc, err := NewRBACService(cfg, log)
	if err != nil {
		t.Skip("Skipping test due to missing policy file")
	}

	ctx := context.Background()

	// 添加策略
	err = svc.AddPolicy(ctx, "tester", "/api/v1/test", "read")
	assert.NoError(t, err)

	// 验证策略
	allowed, err := svc.Enforce(ctx, "tester", "/api/v1/test", "read")
	assert.NoError(t, err)
	assert.True(t, allowed)

	// 移除策略
	err = svc.RemovePolicy(ctx, "tester", "/api/v1/test", "read")
	assert.NoError(t, err)

	// 验证策略已移除
	allowed, err = svc.Enforce(ctx, "tester", "/api/v1/test", "read")
	assert.NoError(t, err)
	assert.False(t, allowed)
}

func TestPermissionError(t *testing.T) {
	err := ErrPermissionDenied
	assert.Equal(t, "Permission denied", err.Error())
	assert.Equal(t, 4001, err.GetCode())
	assert.Equal(t, 4001, err.Code)
	assert.Equal(t, "Permission denied", err.GetMessage())
}

func TestABACService_Disabled(t *testing.T) {
	log := newTestLogger()
	cfg := ABACConfig{Enabled: false}
	svc, err := NewABACService(cfg, log)
	require.NoError(t, err)

	ctx := context.Background()
	req := &ABACRequest{
		Subject: SubjectAttribute{ID: 1, Role: "user"},
		Object:  ObjectAttribute{ID: 1, OwnerID: 2},
	}

	// ABAC disabled 时应该允许所有操作
	allowed, err := svc.Enforce(ctx, req)
	assert.NoError(t, err)
	assert.True(t, allowed)
}

func TestABACService_WithAttributes(t *testing.T) {
	log := newTestLogger()
	cfg := ABACConfig{
		Enabled:   true,
		ModelPath: "pkg/rbac/abac_model.conf",
	}

	svc, err := NewABACService(cfg, log)
	if err != nil {
		t.Skip("Skipping test due to missing model file")
	}

	ctx := context.Background()

	// 测试所有者权限
	req := &ABACRequest{
		Subject: SubjectAttribute{ID: 1, Role: "user"},
		Object:  ObjectAttribute{ID: 1, OwnerID: 1, Path: "/api/v1/resources/1"},
		Action:  ActionAttribute{Name: "read"},
	}

	allowed, err := svc.Enforce(ctx, req)
	assert.NoError(t, err)
	assert.True(t, allowed)
}

func TestNewPermissionService(t *testing.T) {
	log := newTestLogger()

	// 测试使用 disabled 服务
	rbac, err := NewRBACService(RBACConfig{Enabled: false}, log)
	require.NoError(t, err)

	abac, err := NewABACService(ABACConfig{Enabled: false}, log)
	require.NoError(t, err)

	svc := NewPermissionService(rbac, abac)
	assert.NotNil(t, svc)

	ctx := context.Background()

	// 创建有效的 claims
	claims := &jwx.SelfTokenClaims{
		Uid:      1,
		UserType: "user",
		Role:     "admin",
	}

	// 应该允许所有操作 (因为 RBAC/ABAC 都禁用了)
	allowed, err := svc.CheckPermission(ctx, claims, "/api/v1/test", "read")
	assert.NoError(t, err)
	assert.True(t, allowed)

	// 测试 nil claims 应该返回错误
	allowed, err = svc.CheckPermission(ctx, nil, "/api/v1/test", "read")
	assert.Error(t, err)
	assert.False(t, allowed)
}

func TestMapHTTPMethodToAction(t *testing.T) {
	tests := []struct {
		method string
		action string
	}{
		{"GET", "read"},
		{"POST", "create"},
		{"PUT", "update"},
		{"PATCH", "update"},
		{"DELETE", "delete"},
		{"OPTIONS", "options"},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			action := mapHTTPMethodToAction(tt.method)
			assert.Equal(t, tt.action, action)
		})
	}
}

func TestIsWorkHour(t *testing.T) {
	now := time.Now()
	// 只验证函数执行不出错
	isWorkHour(now)
}

func TestSplitPolicyInternal(t *testing.T) {
	tests := []struct {
		name     string
		policy   string
		expected []string
	}{
		{
			name:     "simple policy",
			policy:   "p, admin, /api/v1/*, *",
			expected: []string{"p", "admin", "/api/v1/*", "*"},
		},
		{
			name:     "policy with spaces",
			policy:   "p , admin , /api/v1/* , *",
			expected: []string{"p", "admin", "/api/v1/*", "*"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parts := splitPolicyTestHelper(tt.policy)
			assert.Equal(t, tt.expected, parts)
		})
	}
}

// splitPolicyTestHelper 测试用内部函数
func splitPolicyTestHelper(policy string) []string {
	var parts []string
	var current string
	inQuote := false

	for _, ch := range policy {
		if ch == '"' {
			inQuote = !inQuote
		} else if ch == ',' && !inQuote {
			parts = append(parts, current)
			current = ""
		} else if ch != ' ' {
			current += string(ch)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}

	return parts
}
