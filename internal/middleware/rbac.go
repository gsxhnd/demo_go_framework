package middleware

import (
	"context"
	"strings"
	"time"

	"go_sample_code/internal/errno"
	"go_sample_code/pkg/logger"
	"go_sample_code/pkg/rbac"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// RBACConfig 授权中间件配置
type RBACConfig struct {
	// RBAC 服务
	RBACService rbac.RBACService
	// ABAC 服务（可选）
	ABACService rbac.ABACService
	// 日志记录器
	Log logger.Logger
	// 跳过授权的路径
	SkipPaths []string
	// 是否启用 ABAC
	EnableABAC bool
}

// rbacMiddleware 授权中间件，处理 RBAC 权限检查
func rbacMiddleware(cfg *RBACConfig) fiber.Handler {
	skipPaths := make(map[string]bool)
	for _, path := range cfg.SkipPaths {
		skipPaths[path] = true
	}

	return func(c *fiber.Ctx) error {
		// 检查是否跳过授权
		if skipPaths[c.Path()] {
			return c.Next()
		}

		// 检查用户是否已认证
		role := c.Locals("role")
		if role == nil {
			cfg.Log.WarnCtx(c.UserContext(), "user not authenticated")
			decoded := errno.Decode(nil, errno.TokenInvalidError)
			return c.Status(decoded.GetHTTPStatus()).JSON(decoded)
		}

		roleStr, ok := role.(string)
		if !ok {
			cfg.Log.ErrorCtx(c.UserContext(), "invalid role type in context")
			decoded := errno.Decode(nil, errno.InternalServerError)
			return c.Status(decoded.GetHTTPStatus()).JSON(decoded)
		}

		ctx := c.UserContext()
		method := c.Method()
		path := c.Path()

		// 获取用户 ID
		uid := getUIDFromContext(c)

		// RBAC 权限检查
		allowed, err := checkRBACPermission(ctx, cfg.RBACService, roleStr, path, method)
		if err != nil {
			cfg.Log.ErrorCtx(ctx, "RBAC check failed",
				zap.String("role", roleStr),
				zap.String("path", path),
				zap.String("method", method),
				zap.Error(err))
			decoded := errno.Decode(nil, errno.InternalServerError)
			return c.Status(decoded.GetHTTPStatus()).JSON(decoded)
		}

		if !allowed {
			cfg.Log.WarnCtx(ctx, "permission denied",
				zap.Uint64("uid", uid),
				zap.String("role", roleStr),
				zap.String("path", path),
				zap.String("method", method))
			decoded := errno.Decode(nil, errno.PermissionDeniedError)
			return c.Status(decoded.GetHTTPStatus()).JSON(decoded)
		}

		// ABAC 权限检查（可选）
		if cfg.EnableABAC && cfg.ABACService != nil {
			allowed, err = checkABACPermissionFromCtx(ctx, cfg.ABACService, c)
			if err != nil {
				cfg.Log.ErrorCtx(ctx, "ABAC check failed", zap.Error(err))
				decoded := errno.Decode(nil, errno.InternalServerError)
				return c.Status(decoded.GetHTTPStatus()).JSON(decoded)
			}

			if !allowed {
				cfg.Log.WarnCtx(ctx, "ABAC permission denied",
					zap.Uint64("uid", uid))
				decoded := errno.Decode(nil, errno.PermissionDeniedError)
				return c.Status(decoded.GetHTTPStatus()).JSON(decoded)
			}
		}

		cfg.Log.DebugCtx(ctx, "permission granted",
			zap.Uint64("uid", uid),
			zap.String("role", roleStr),
			zap.String("path", path),
			zap.String("method", method))

		return c.Next()
	}
}

// getUIDFromContext 从上下文获取用户 ID
func getUIDFromContext(c *fiber.Ctx) uint64 {
	uid := c.Locals("uid")
	if uid == nil {
		return 0
	}
	switch v := uid.(type) {
	case uint64:
		return v
	case int64:
		return uint64(v)
	case int:
		return uint64(v)
	default:
		return 0
	}
}

// checkRBACPermission 检查 RBAC 权限
func checkRBACPermission(ctx context.Context, service rbac.RBACService, role, path, method string) (bool, error) {
	// 方法映射到动作
	action := mapHTTPMethodToAction(method)

	// 检查直接权限
	allowed, err := service.Enforce(ctx, role, path, action)
	if err != nil {
		return false, err
	}

	if allowed {
		return true, nil
	}

	// 检查通配符权限
	return service.Enforce(ctx, role, path, "*")
}

// checkABACPermissionFromCtx 从 Fiber 上下文检查 ABAC 权限
func checkABACPermissionFromCtx(ctx context.Context, service rbac.ABACService, c *fiber.Ctx) (bool, error) {
	uid := getUIDFromContext(c)
	role, _ := c.Locals("role").(string)

	req := &rbac.ABACRequest{
		Subject: rbac.SubjectAttribute{
			ID:   uid,
			Role: role,
			IP:   c.IP(),
		},
		Object: rbac.ObjectAttribute{
			Path: c.Path(),
		},
		Action: rbac.ActionAttribute{
			Name:   mapHTTPMethodToAction(c.Method()),
			Method: c.Method(),
		},
		Environment: rbac.EnvironmentAttribute{
			Time:       time.Now(),
			DayOfWeek:  int(time.Now().Weekday()),
			HourOfDay:  time.Now().Hour(),
			IsWorkHour: isWorkHour(time.Now()),
		},
	}

	return service.Enforce(ctx, req)
}

// mapHTTPMethodToAction 将 HTTP 方法映射到动作
func mapHTTPMethodToAction(method string) string {
	switch method {
	case "GET":
		return "read"
	case "POST":
		return "create"
	case "PUT", "PATCH":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return strings.ToLower(method)
	}
}

// isWorkHour 检查是否在工作时间内
func isWorkHour(t time.Time) bool {
	hour := t.Hour()
	weekday := t.Weekday()

	// 周末不工作
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}

	// 工作时间 9:00 - 18:00
	return hour >= 9 && hour < 18
}

// RBAC 返回授权中间件处理函数
func RBAC(rbacService rbac.RBACService, log logger.Logger) fiber.Handler {
	return rbacMiddleware(&RBACConfig{
		RBACService: rbacService,
		Log:         log,
		SkipPaths:   []string{},
		EnableABAC:  false,
	})
}

// RBACWithConfig 返回授权中间件处理函数，支持完整配置
func RBACWithConfig(cfg *RBACConfig) fiber.Handler {
	return rbacMiddleware(cfg)
}

// RequireRole 返回一个中间件，要求用户具有指定角色
func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role")
		if userRole == nil {
			decoded := errno.Decode(nil, errno.TokenInvalidError)
			return c.Status(decoded.GetHTTPStatus()).JSON(decoded)
		}

		roleStr, ok := userRole.(string)
		if !ok {
			decoded := errno.Decode(nil, errno.InternalServerError)
			return c.Status(decoded.GetHTTPStatus()).JSON(decoded)
		}

		// 检查是否匹配要求的角色
		for _, requiredRole := range roles {
			if roleStr == requiredRole {
				return c.Next()
			}
		}

		// 检查是否有继承的角色（从 Locals 获取）
		userRoles := c.Locals("roles")
		if userRoles != nil {
			rolesList, ok := userRoles.([]string)
			if ok {
				for _, requiredRole := range roles {
					for _, ur := range rolesList {
						if ur == requiredRole {
							return c.Next()
						}
					}
				}
			}
		}

		decoded := errno.Decode(nil, errno.PermissionDeniedError)
		return c.Status(decoded.GetHTTPStatus()).JSON(decoded)
	}
}

// RequirePermission 返回一个中间件，要求用户具有指定权限
func RequirePermission(obj, act string, rbacService rbac.RBACService, log logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		role := c.Locals("role")
		if role == nil {
			decoded := errno.Decode(nil, errno.TokenInvalidError)
			return c.Status(decoded.GetHTTPStatus()).JSON(decoded)
		}

		roleStr, ok := role.(string)
		if !ok {
			decoded := errno.Decode(nil, errno.InternalServerError)
			return c.Status(decoded.GetHTTPStatus()).JSON(decoded)
		}

		allowed, err := checkRBACPermission(ctx, rbacService, roleStr, obj, act)
		if err != nil {
			log.ErrorCtx(ctx, "permission check failed", zap.Error(err))
			decoded := errno.Decode(nil, errno.InternalServerError)
			return c.Status(decoded.GetHTTPStatus()).JSON(decoded)
		}

		if !allowed {
			decoded := errno.Decode(nil, errno.PermissionDeniedError)
			return c.Status(decoded.GetHTTPStatus()).JSON(decoded)
		}

		return c.Next()
	}
}
