package rbac

import (
	"context"
	"strings"
	"time"

	"go_sample_code/pkg/jwx"
	"go_sample_code/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// MiddlewareConfig 中间件配置
type MiddlewareConfig struct {
	// JWT Provider 用于解析 token
	JWTProvider jwx.JwtProvider
	// RBAC 服务
	RBACService RBACService
	// ABAC 服务
	ABACService ABACService
	// 日志
	Log logger.Logger
	// 跳过认证的路径
	SkipPaths []string
	// 是否启用 ABAC
	EnableABAC bool
}

// authContextKey 认证上下文键
type authContextKey string

const (
	// UserClaimsKey 用户声明键
	UserClaimsKey authContextKey = "user_claims"
	// UserRolesKey 用户角色键
	UserRolesKey authContextKey = "user_roles"
)

// AuthorizeMiddleware 权限验证中间件
func AuthorizeMiddleware(cfg *MiddlewareConfig) fiber.Handler {
	skipPaths := make(map[string]bool)
	for _, path := range cfg.SkipPaths {
		skipPaths[path] = true
	}

	return func(c *fiber.Ctx) error {
		// 检查是否跳过认证
		if skipPaths[c.Path()] {
			return c.Next()
		}

		ctx := c.UserContext()

		// 获取 Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			cfg.Log.WarnCtx(ctx, "missing authorization header",
				zap.String("path", c.Path()))
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    1101,
				"message": "Missing authorization header",
			})
		}

		// 提取 Bearer token
		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			cfg.Log.WarnCtx(ctx, "invalid authorization header format")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    1101,
				"message": "Invalid authorization header format",
			})
		}

		token := tokenParts[1]

		// 验证 token
		claims, err := cfg.JWTProvider.ValidateSelfToken(token)
		if err != nil {
			cfg.Log.WarnCtx(ctx, "token validation failed",
				zap.Error(err))
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    1101,
				"message": "Invalid token",
			})
		}

		// 设置用户上下文
		c.Locals("uid", claims.Uid)
		c.Locals("user_type", claims.UserType)
		c.Locals("role", claims.Role)

		// 获取用户角色列表
		roles, err := cfg.RBACService.GetRolesForUser(ctx, claims.Role)
		if err != nil {
			cfg.Log.ErrorCtx(ctx, "failed to get user roles",
				zap.Uint64("uid", claims.Uid),
				zap.Error(err))
			// 如果获取失败，使用 token 中的角色
			roles = []string{claims.Role}
		}
		c.Locals("roles", roles)

		// 获取请求方法
		method := c.Method()
		path := c.Path()

		// RBAC 权限检查
		allowed, err := checkRBACPermission(ctx, cfg.RBACService, claims.Role, path, method)
		if err != nil {
			cfg.Log.ErrorCtx(ctx, "RBAC check failed",
				zap.String("role", claims.Role),
				zap.String("path", path),
				zap.String("method", method),
				zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    5000,
				"message": "Permission check failed",
			})
		}

		if !allowed {
			cfg.Log.WarnCtx(ctx, "permission denied",
				zap.Uint64("uid", claims.Uid),
				zap.String("role", claims.Role),
				zap.String("path", path),
				zap.String("method", method))
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"code":    4001,
				"message": "Permission denied",
			})
		}

		// ABAC 权限检查 (可选)
		if cfg.EnableABAC && cfg.ABACService != nil {
			allowed, err = checkABACPermission(ctx, cfg.ABACService, claims, c)
			if err != nil {
				cfg.Log.ErrorCtx(ctx, "ABAC check failed", zap.Error(err))
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"code":    5000,
					"message": "Permission check failed",
				})
			}

			if !allowed {
				cfg.Log.WarnCtx(ctx, "ABAC permission denied",
					zap.Uint64("uid", claims.Uid))
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"code":    4001,
					"message": "Permission denied by attribute policy",
				})
			}
		}

		cfg.Log.DebugCtx(ctx, "permission granted",
			zap.Uint64("uid", claims.Uid),
			zap.String("role", claims.Role),
			zap.String("path", path),
			zap.String("method", method))

		return c.Next()
	}
}

// checkRBACPermission 检查 RBAC 权限
func checkRBACPermission(ctx context.Context, service RBACService, role, path, method string) (bool, error) {
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
	allowed, err = service.Enforce(ctx, role, path, "*")
	if err != nil {
		return false, err
	}

	return allowed, nil
}

// checkABACPermission 检查 ABAC 权限
func checkABACPermission(ctx context.Context, service ABACService, claims *jwx.SelfTokenClaims, c *fiber.Ctx) (bool, error) {
	req := &ABACRequest{
		Subject: SubjectAttribute{
			ID:       claims.Uid,
			Username: "", // 可以从 context 中获取
			Role:     claims.Role,
			IP:       c.IP(),
		},
		Object: ObjectAttribute{
			Path: c.Path(),
		},
		Action: ActionAttribute{
			Name:   mapHTTPMethodToAction(c.Method()),
			Method: c.Method(),
		},
		Environment: EnvironmentAttribute{
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

// RequireRole 返回一个中间件，要求用户具有指定角色
func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role")
		if userRole == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    1101,
				"message": "Not authenticated",
			})
		}

		roleStr, ok := userRole.(string)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    5000,
				"message": "Invalid role type",
			})
		}

		// 检查是否匹配要求的角色
		for _, requiredRole := range roles {
			if roleStr == requiredRole {
				return c.Next()
			}
		}

		// 检查是否有继承的角色
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

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"code":    4001,
			"message": "Insufficient permissions",
		})
	}
}

// RequirePermission 返回一个中间件，要求用户具有指定权限
func RequirePermission(obj, act string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		role := c.Locals("role")
		if role == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    1101,
				"message": "Not authenticated",
			})
		}

		roleStr, ok := role.(string)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    5000,
				"message": "Invalid role type",
			})
		}

		allowed, err := checkRBACPermission(ctx, nil, roleStr, obj, act)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    5000,
				"message": "Permission check failed",
			})
		}

		if !allowed {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"code":    4001,
				"message": "Permission denied",
			})
		}

		return c.Next()
	}
}
