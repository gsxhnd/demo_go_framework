package rbac

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Knetic/govaluate"
	"go_sample_code/pkg/logger"

	"go.uber.org/zap"
)

// SubjectAttribute 主体属性 (用户属性)
type SubjectAttribute struct {
	ID       uint64            `json:"id"`
	Username string            `json:"username"`
	Role     string            `json:"role"`
	Roles    []string          `json:"roles"`
	Email    string            `json:"email"`
	IP       string            `json:"ip"`
	DeptID   uint64            `json:"dept_id"`
	Ext      map[string]string `json:"ext"` // 扩展属性
}

// ObjectAttribute 客体属性 (资源属性)
type ObjectAttribute struct {
	ID       uint64                 `json:"id"`
	Type     string                 `json:"type"`
	OwnerID  uint64                 `json:"owner_id"`
	Name     string                 `json:"name"`
	Path     string                 `json:"path"`
	Tags     []string               `json:"tags"`
	IsPublic bool                   `json:"is_public"`
	Ext      map[string]interface{} `json:"ext"` // 扩展属性
}

// ActionAttribute 动作属性
type ActionAttribute struct {
	Name   string `json:"name"`   // read, write, delete, etc.
	Method string `json:"method"` // HTTP method
}

// EnvironmentAttribute 环境属性
type EnvironmentAttribute struct {
	Time       time.Time `json:"time"`
	DayOfWeek  int       `json:"day_of_week"`  // 0-6, 0 is Sunday
	HourOfDay  int       `json:"hour_of_day"`  // 0-23
	IsWorkHour bool      `json:"is_work_hour"` // 9:00-18:00
}

// ABACRequest ABAC 权限请求
type ABACRequest struct {
	Subject     SubjectAttribute     `json:"subject"`
	Object      ObjectAttribute      `json:"object"`
	Action      ActionAttribute      `json:"action"`
	Environment EnvironmentAttribute `json:"environment"`
}

// ABACService ABAC 权限服务接口
type ABACService interface {
	// Enforce 检查权限 (基于属性表达式)
	Enforce(ctx context.Context, req *ABACRequest) (bool, error)
	// EnforceWithExpression 使用自定义表达式检查权限
	EnforceWithExpression(ctx context.Context, expr string, req *ABACRequest) (bool, error)
	// ValidateExpression 验证表达式有效性
	ValidateExpression(ctx context.Context, expr string) error
	// ReloadPolicy 重新加载策略
	ReloadPolicy(ctx context.Context) error
}

// abacService ABAC 服务实现
type abacService struct {
	config   *ABACConfig
	log      logger.Logger
	mu       sync.RWMutex
	// 存储 ABAC 策略
	policies []abacPolicy
}

// abacPolicy ABAC 策略
type abacPolicy struct {
	Expr     string
	Obj      string
	Act      string
	// 预编译的表达式
	compiled *govaluate.EvaluableExpression
}

// NewABACService 创建 ABAC 服务
func NewABACService(cfg ABACConfig, log logger.Logger) (ABACService, error) {
	if !cfg.Enabled {
		log.Warn("ABAC is disabled")
		return &disabledABACService{log: log}, nil
	}

	// 确保路径不为空
	modelPath := cfg.ModelPath
	if modelPath == "" {
		modelPath = "pkg/rbac/abac_model.conf"
	}

	// 验证文件是否存在
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		log.Warn("ABAC model file not found, using default policies")
	}

	service := &abacService{
		config:   &cfg,
		log:      log,
		policies: make([]abacPolicy, 0),
	}

	// 添加默认策略
	service.initDefaultPolicies()

	log.Info("ABAC service initialized", zap.String("model_path", modelPath))

	return service, nil
}

// NewABACServiceWithDefaults 使用默认配置创建 ABAC 服务
func NewABACServiceWithDefaults(log logger.Logger) (ABACService, error) {
	return NewABACService(DefaultConfig().ABAC, log)
}

// initDefaultPolicies 初始化默认策略
func (s *abacService) initDefaultPolicies() {
	// 使用简单的布尔表达式，不使用引号比较
	s.policies = []abacPolicy{
		// 资源所有者可以访问自己的资源
		{Expr: `r_sub_id == r_obj_owner_id`, Obj: "/api/v1/resources/*", Act: "read"},
		{Expr: `r_sub_id == r_obj_owner_id`, Obj: "/api/v1/resources/*", Act: "write"},
		{Expr: `r_sub_id == r_obj_owner_id`, Obj: "/api/v1/resources/*", Act: "delete"},
		// 公开资源可以访问
		{Expr: `r_obj_is_public == true`, Obj: "/api/v1/public/*", Act: "read"},
		// 管理员可以访问所有资源
		{Expr: `r_sub_role_matches_admin`, Obj: "/*", Act: "*"},
		// 工作时间内用户可以访问
		{Expr: `r_sub_role_is_user && r_env_is_work_hour == true`, Obj: "/api/v1/work/*", Act: "read"},
	}

	// 预编译表达式
	for i := range s.policies {
		expr, err := govaluate.NewEvaluableExpression(s.policies[i].Expr)
		if err != nil {
			s.log.Warn("failed to compile expression",
				zap.String("expr", s.policies[i].Expr),
				zap.Error(err))
			continue
		}
		s.policies[i].compiled = expr
	}
}

// Enforce 检查权限
func (s *abacService) Enforce(ctx context.Context, req *ABACRequest) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 遍历所有策略
	for _, policy := range s.policies {
		// 检查对象和动作是否匹配
		if !matchObjectAndAction(req, policy.Obj, policy.Act) {
			continue
		}

		// 使用表达式进行评估
		result, err := s.evaluatePolicy(policy, req)
		if err != nil {
			s.log.WarnCtx(ctx, "expression evaluation failed",
				zap.String("expression", policy.Expr),
				zap.Error(err))
			continue
		}

		if result {
			s.log.DebugCtx(ctx, "ABAC enforce passed",
				zap.String("expression", policy.Expr),
				zap.String("object", policy.Obj),
				zap.String("action", policy.Act))
			return true, nil
		}
	}

	s.log.DebugCtx(ctx, "ABAC enforce denied",
		zap.String("subject", req.Subject.Username),
		zap.String("object", req.Object.Path),
		zap.String("action", req.Action.Name))

	return false, nil
}

// evaluatePolicy 评估单个策略
func (s *abacService) evaluatePolicy(policy abacPolicy, req *ABACRequest) (bool, error) {
	if policy.compiled == nil {
		return false, fmt.Errorf("expression not compiled")
	}

	// 构建表达式参数
	parameters := make(map[string]interface{})
	parameters["r_sub_id"] = float64(req.Subject.ID)
	parameters["r_sub_role"] = req.Subject.Role
	parameters["r_sub_dept_id"] = float64(req.Subject.DeptID)
	parameters["r_sub_role_matches_admin"] = req.Subject.Role == "admin"
	parameters["r_sub_role_is_user"] = req.Subject.Role == "user"
	parameters["r_obj_id"] = float64(req.Object.ID)
	parameters["r_obj_owner_id"] = float64(req.Object.OwnerID)
	parameters["r_obj_is_public"] = req.Object.IsPublic
	parameters["r_env_day_of_week"] = float64(req.Environment.DayOfWeek)
	parameters["r_env_hour_of_day"] = float64(req.Environment.HourOfDay)
	parameters["r_env_is_work_hour"] = req.Environment.IsWorkHour

	result, err := policy.compiled.Evaluate(parameters)
	if err != nil {
		return false, err
	}

	if b, ok := result.(bool); ok {
		return b, nil
	}

	return false, fmt.Errorf("expression result is not boolean")
}

// EnforceWithExpression 使用自定义表达式检查权限
func (s *abacService) EnforceWithExpression(ctx context.Context, expr string, req *ABACRequest) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	policy := abacPolicy{Expr: expr}
	if err := s.prepareExpression(&policy); err != nil {
		return false, err
	}

	return s.evaluatePolicy(policy, req)
}

// ValidateExpression 验证表达式有效性
func (s *abacService) ValidateExpression(ctx context.Context, expr string) error {
	policy := abacPolicy{Expr: expr}
	return s.prepareExpression(&policy)
}

// prepareExpression 准备表达式
func (s *abacService) prepareExpression(policy *abacPolicy) error {
	var err error
	policy.compiled, err = govaluate.NewEvaluableExpression(policy.Expr)
	return err
}

// ReloadPolicy 重新加载策略
func (s *abacService) ReloadPolicy(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 重置策略
	s.policies = make([]abacPolicy, 0)
	s.initDefaultPolicies()

	s.log.InfoCtx(ctx, "ABAC policies reloaded")
	return nil
}

// matchObjectAndAction 匹配对象和动作
func matchObjectAndAction(req *ABACRequest, pattern, action string) bool {
	// 匹配对象路径
	if !matchPath(req.Object.Path, pattern) {
		return false
	}

	// 匹配动作
	if action == "*" || action == req.Action.Name || action == req.Action.Method {
		return true
	}

	return false
}

// matchPath 匹配路径模式
func matchPath(path, pattern string) bool {
	if pattern == "*" || pattern == "/*" {
		return true
	}

	// 支持简单通配符
	if len(pattern) > 0 && pattern[len(pattern)-1] == '*' {
		prefix := pattern[:len(pattern)-1]
		return len(path) >= len(prefix) && path[:len(prefix)] == prefix
	}

	return path == pattern
}

// disabledABACService ABAC 未启用时的服务实现
type disabledABACService struct {
	log logger.Logger
}

func (s *disabledABACService) Enforce(ctx context.Context, req *ABACRequest) (bool, error) {
	s.log.WarnCtx(ctx, "ABAC is disabled, allowing all access")
	return true, nil
}

func (s *disabledABACService) EnforceWithExpression(ctx context.Context, expr string, req *ABACRequest) (bool, error) {
	return true, nil
}

func (s *disabledABACService) ValidateExpression(ctx context.Context, expr string) error {
	return nil
}

func (s *disabledABACService) ReloadPolicy(ctx context.Context) error {
	return nil
}
