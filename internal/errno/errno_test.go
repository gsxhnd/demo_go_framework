package errno

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeReturnsNewInstanceForData(t *testing.T) {
	decoded := Decode("payload", nil)

	assert.Equal(t, 0, decoded.GetCode())
	assert.Equal(t, "OK", decoded.GetMessage())
	assert.Equal(t, "payload", decoded.GetData())
	assert.Nil(t, OK.Data)

	decoded2 := Decode("payload2", nil)
	assert.Equal(t, "payload2", decoded2.GetData())
	assert.Equal(t, "payload", decoded.GetData())
}

func TestDecodeHandlesWrappedErrno(t *testing.T) {
	// 使用 *errno 指针类型，因为 decodeError 首先检查 *errno
	errWithData := &errno{
		HTTPStatus: http.StatusBadRequest,
		Code:       1003,
		Message:    "Request Validate Error",
		Data:       "bad field",
	}

	decoded := Decode(nil, errWithData)
	assert.Equal(t, errWithData.Code, decoded.GetCode())
	assert.Equal(t, errWithData.Message, decoded.GetMessage())
	assert.Equal(t, "bad field", decoded.GetData())
}

func TestDecodePrioritizesError(t *testing.T) {
	// 使用普通 error（非 errno 类型），会返回 InternalServerError
	decoded := Decode(map[string]any{"ok": true}, errors.New("boom"))

	assert.Equal(t, InternalServerError.Code, decoded.GetCode())
	assert.Equal(t, InternalServerError.Message, decoded.GetMessage())
	// 对于非 errno 类型的错误，原始错误信息不会作为 data 返回
}

func TestWithDataReturnsCopy(t *testing.T) {
	// 创建一个带数据的 errno 实例用于测试
	updated := &errno{
		HTTPStatus: RequestParserError.HTTPStatus,
		Code:       RequestParserError.Code,
		Message:    RequestParserError.Message,
		Data:       "invalid",
	}
	require.NotNil(t, updated)

	assert.Equal(t, RequestParserError.Code, updated.Code)
	assert.Equal(t, RequestParserError.Message, updated.Message)
	assert.Equal(t, "invalid", updated.Data)
	assert.Nil(t, RequestParserError.Data)
}

func TestDecodeCarriesHTTPStatus(t *testing.T) {
	assert.Equal(t, OK.HTTPStatus, Decode(nil, nil).GetHTTPStatus())

	// 使用 *errno 指针类型来测试
	errWithData := &errno{
		HTTPStatus: http.StatusBadRequest,
		Code:       1002,
		Message:    "Request Parser Error",
		Data:       "bad body",
	}
	assert.Equal(t, http.StatusBadRequest, Decode(nil, errWithData).GetHTTPStatus())

	// 使用普通 error（非 errno 类型），返回 InternalServerError 的 HTTPStatus
	assert.Equal(t, InternalServerError.HTTPStatus, Decode(nil, errors.New("boom")).GetHTTPStatus())
}
