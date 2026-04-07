package errno

import "net/http"

// Business errors for service layer
var (
	// User errors
	UserNotFoundError      = errno{HTTPStatus: http.StatusNotFound, Code: 2001, Message: "User not found", Data: nil}
	UserAlreadyExistsError = errno{HTTPStatus: http.StatusConflict, Code: 2002, Message: "User already exists", Data: nil}
	UserCreateFailedError  = errno{HTTPStatus: http.StatusInternalServerError, Code: 2003, Message: "Failed to create user", Data: nil}
	UserUpdateFailedError  = errno{HTTPStatus: http.StatusInternalServerError, Code: 2004, Message: "Failed to update user", Data: nil}
	UserDeleteFailedError  = errno{HTTPStatus: http.StatusInternalServerError, Code: 2005, Message: "Failed to delete user", Data: nil}
	InvalidUserIDError     = errno{HTTPStatus: http.StatusBadRequest, Code: 2006, Message: "Invalid user ID", Data: nil}
	InvalidEmailError      = errno{HTTPStatus: http.StatusBadRequest, Code: 2007, Message: "Invalid email format", Data: nil}
	InvalidUsernameError   = errno{HTTPStatus: http.StatusBadRequest, Code: 2008, Message: "Invalid username format", Data: nil}

	// Pagination errors
	InvalidPageError     = errno{HTTPStatus: http.StatusBadRequest, Code: 3001, Message: "Invalid page number", Data: nil}
	InvalidPageSizeError = errno{HTTPStatus: http.StatusBadRequest, Code: 3002, Message: "Invalid page size", Data: nil}
)
