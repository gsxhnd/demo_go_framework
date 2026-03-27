package errno

import (
	"errors"
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
	wrapper := errors.New("wrapper")
	err := errors.Join(wrapper, RequestValidateError.WithData("bad field"))

	decoded := Decode(nil, err)
	assert.Equal(t, RequestValidateError.Code, decoded.GetCode())
	assert.Equal(t, RequestValidateError.Message, decoded.GetMessage())
	assert.Equal(t, "bad field", decoded.GetData())
}

func TestDecodePrioritizesError(t *testing.T) {
	decoded := Decode(map[string]any{"ok": true}, errors.New("boom"))

	assert.Equal(t, InternalServerError.Code, decoded.GetCode())
	assert.Equal(t, InternalServerError.Message, decoded.GetMessage())
	assert.Equal(t, "boom", decoded.GetData())
}

func TestWithDataReturnsCopy(t *testing.T) {
	updated := RequestParserError.WithData("invalid")
	require.NotNil(t, updated)

	assert.Equal(t, RequestParserError.Code, updated.Code)
	assert.Equal(t, RequestParserError.Message, updated.Message)
	assert.Equal(t, "invalid", updated.Data)
	assert.Nil(t, RequestParserError.Data)
}

func TestDecodeCarriesHTTPStatus(t *testing.T) {
	assert.Equal(t, OK.HTTPStatus, Decode(nil, nil).GetHTTPStatus())
	assert.Equal(t, RequestParserError.HTTPStatus, Decode(nil, RequestParserError.WithData("bad body")).GetHTTPStatus())
	assert.Equal(t, InternalServerError.HTTPStatus, Decode(nil, errors.New("boom")).GetHTTPStatus())
}
