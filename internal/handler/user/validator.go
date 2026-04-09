package user

import (
	"reflect"
	"strconv"

	"go_sample_code/internal/errno"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// parseAndValidateBody 解析并校验请求体
// 返回 true 表示校验失败（已设置响应），返回 false 表示校验通过
func (h *handler) parseAndValidateBody(c *fiber.Ctx, req any) bool {
	ctx := c.UserContext()
	if err := c.BodyParser(req); err != nil {
		h.log.ErrorCtx(ctx, "failed to parse request body", zap.Error(err))
		_ = c.Status(errno.RequestParserError.GetHTTPStatus()).JSON(errno.Decode(nil, errno.RequestParserError))
		return true
	}

	if err := h.validate.Struct(req); err != nil {
		h.log.WarnCtx(ctx, "request validation failed", zap.Error(err))
		_ = c.Status(errno.RequestValidateError.GetHTTPStatus()).JSON(errno.Decode(nil, errno.RequestValidateError))
		return true
	}

	return false
}

// parseAndValidateQuery 解析并校验查询参数
// 返回 true 表示校验失败（已设置响应），返回 false 表示校验通过
func (h *handler) parseAndValidateQuery(c *fiber.Ctx, req any) bool {
	ctx := c.UserContext()
	if err := c.QueryParser(req); err != nil {
		h.log.ErrorCtx(ctx, "failed to parse query parameters", zap.Error(err))
		_ = c.Status(errno.RequestParserError.GetHTTPStatus()).JSON(errno.Decode(nil, errno.RequestParserError))
		return true
	}

	if err := h.validate.Struct(req); err != nil {
		h.log.WarnCtx(ctx, "query validation failed", zap.Error(err))
		_ = c.Status(errno.RequestValidateError.GetHTTPStatus()).JSON(errno.Decode(nil, errno.RequestValidateError))
		return true
	}

	return false
}

// parseAndValidateParams 解析并校验路径参数
// 返回 true 表示校验失败（已设置响应），返回 false 表示校验通过
func (h *handler) parseAndValidateParams(c *fiber.Ctx, req any) bool {
	ctx := c.UserContext()

	// 手动解析 params，因为 Fiber 的 QueryParser 不能解析 path params
	if err := parseFiberParams(c, req); err != nil {
		h.log.ErrorCtx(ctx, "failed to parse path parameters", zap.Error(err))
		_ = c.Status(errno.RequestParserError.GetHTTPStatus()).JSON(errno.Decode(nil, errno.RequestParserError))
		return true
	}

	if err := h.validate.Struct(req); err != nil {
		h.log.WarnCtx(ctx, "path params validation failed", zap.Error(err))
		_ = c.Status(errno.RequestValidateError.GetHTTPStatus()).JSON(errno.Decode(nil, errno.RequestValidateError))
		return true
	}

	return false
}

// parseFiberParams 解析 Fiber 的路径参数到结构体
func parseFiberParams(c *fiber.Ctx, req any) error {
	// 获取结构体的字段信息
	reqVal := reflect.ValueOf(req).Elem()
	reqType := reqVal.Type()

	for i := 0; i < reqType.NumField(); i++ {
		field := reqType.Field(i)
		paramTag := field.Tag.Get("params")
		if paramTag == "" || paramTag == "-" {
			continue
		}

		paramValue := c.Params(paramTag)
		if paramValue == "" {
			continue
		}

		fieldVal := reqVal.Field(i)
		switch field.Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intVal, err := strconv.ParseInt(paramValue, 10, 64)
			if err != nil {
				return err
			}
			if fieldVal.CanSet() {
				fieldVal.SetInt(intVal)
			}
		case reflect.String:
			if fieldVal.CanSet() {
				fieldVal.SetString(paramValue)
			}
		}
	}
	return nil
}
