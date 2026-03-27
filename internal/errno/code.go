package errno

import "net/http"

var (
	// common
	OK                     = errno{HTTPStatus: http.StatusOK, Code: 0, Message: "OK", Data: nil}
	InternalServerError    = errno{HTTPStatus: http.StatusInternalServerError, Code: 1000, Message: "Internal Server Error", Data: nil}
	Deprecated             = errno{HTTPStatus: http.StatusGone, Code: 1001, Message: "API Deprecated", Data: nil}
	RequestParserError     = errno{HTTPStatus: http.StatusBadRequest, Code: 1002, Message: "Request Parser Error", Data: nil}
	RequestValidateError   = errno{HTTPStatus: http.StatusBadRequest, Code: 1003, Message: "Request Validate Error", Data: nil}
	RequestConversionError = errno{HTTPStatus: http.StatusBadRequest, Code: 1004, Message: "Request Conversion Error", Data: nil}
	DataConversionError    = errno{HTTPStatus: http.StatusInternalServerError, Code: 1005, Message: "Data Conversion Error", Data: nil}
	TokenInvalidError      = errno{HTTPStatus: http.StatusUnauthorized, Code: 1101, Message: "Token Invalid Error", Data: nil}
	TokenParserError       = errno{HTTPStatus: http.StatusUnauthorized, Code: 1102, Message: "Token Parser Error", Data: nil}
	PermissionDeniedError  = errno{HTTPStatus: http.StatusForbidden, Code: 1103, Message: "Permission Denied", Data: nil}
	RetrievingFileError    = errno{HTTPStatus: http.StatusNotFound, Code: 1201, Message: "Retrieving File Error", Data: nil}
	// data error
	DatabaseError           = errno{HTTPStatus: http.StatusInternalServerError, Code: 1301, Message: "Database Error", Data: nil}
	DatabaseConversionError = errno{HTTPStatus: http.StatusInternalServerError, Code: 1302, Message: "Database Conversion Error", Data: nil}
)
