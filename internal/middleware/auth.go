package middleware

import (
	"strings"

	"go_sample_code/internal/errno"
	"go_sample_code/pkg/jwx"
	"go_sample_code/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// AuthConfig 认证中间件配置
type AuthConfig struct {
	// JWT Provider 用于解析和验证 token
	JWTProvider jwx.JwtProvider
	// 日志记录器
	Log logger.Logger
	// 跳过认证的路径（支持前缀匹配）
	SkipPaths []string
}

// authMiddleware 认证中间件，处理 Token 解析和用户上下文设置
func authMiddleware(cfg *AuthConfig) fiber.Handler {
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
			decoded := errno.Decode(nil, errno.TokenInvalidError)
			return c.Status(decoded.GetHTTPStatus()).JSON(decoded)
		}

		// 提取 Bearer token
		token, err := extractBearerToken(authHeader)
		if err != nil {
			cfg.Log.WarnCtx(ctx, "invalid authorization header format",
				zap.String("auth_header", authHeader))
			decoded := errno.Decode(nil, errno.TokenInvalidError)
			return c.Status(decoded.GetHTTPStatus()).JSON(decoded)
		}

		// 验证 token 并获取用户声明
		claims, err := cfg.JWTProvider.ValidateSelfToken(token)
		if err != nil {
			cfg.Log.WarnCtx(ctx, "token validation failed",
				zap.Error(err))
			decoded := errno.Decode(nil, errno.TokenInvalidError)
			return c.Status(decoded.GetHTTPStatus()).JSON(decoded)
		}

		// 设置用户上下文信息到 Locals
		setUserContext(c, claims, token)

		cfg.Log.DebugCtx(ctx, "user authenticated",
			zap.Uint64("uid", claims.Uid),
			zap.String("role", claims.Role),
			zap.String("user_type", claims.UserType))

		return c.Next()
	}
}

// extractBearerToken 从 Authorization header 中提取 Bearer token
func extractBearerToken(authHeader string) (string, error) {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errno.TokenInvalidError
	}
	return strings.TrimSpace(parts[1]), nil
}

// setUserContext 设置用户上下文信息
func setUserContext(c *fiber.Ctx, claims *jwx.SelfTokenClaims, token string) {
	// 使用 Fiber 的 Locals 存储用户信息
	c.Locals("uid", claims.Uid)
	c.Locals("user_type", claims.UserType)
	c.Locals("role", claims.Role)
	c.Locals("token", token)

	// 同时设置到 fiber.Ctx 以保持向后兼容
	c.Context().SetUserValue("uid", claims.Uid)
	c.Context().SetUserValue("user_type", claims.UserType)
	c.Context().SetUserValue("role", claims.Role)
	c.Context().SetUserValue("token", token)
}

// Auth 返回认证中间件处理函数
// 使用默认配置时，跳过路径为空
func Auth(jwtProvider jwx.JwtProvider, log logger.Logger) fiber.Handler {
	return authMiddleware(&AuthConfig{
		JWTProvider: jwtProvider,
		Log:         log,
		SkipPaths:   []string{},
	})
}

// AuthWithSkipPaths 返回认证中间件处理函数，支持跳过路径配置
func AuthWithSkipPaths(jwtProvider jwx.JwtProvider, log logger.Logger, skipPaths []string) fiber.Handler {
	return authMiddleware(&AuthConfig{
		JWTProvider: jwtProvider,
		Log:         log,
		SkipPaths:   skipPaths,
	})
}
