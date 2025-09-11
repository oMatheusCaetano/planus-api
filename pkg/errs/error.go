package errs

import "net/http"

type Error struct {
    Code    int `json:"code"`
    Message string
}

func (err *Error) Error() string {
    return err.Message
}

func NotFound(message string) *Error {
    return &Error{
        Code:    http.StatusNotFound,
        Message: message,
    }
}

func BadRequest(message string) *Error {
    return &Error{
        Code:    http.StatusBadRequest,
        Message: message,
    }
}

// From converts a standard error to an *Error type.
// If err is already an *Error, it returns it directly.
// Otherwise, it maps the generic error to an *Error.
func From(err error) *Error {
    if err == nil {
        return nil
    }

    // Check if err is already an *Error
    if e, ok := err.(*Error); ok {
        return e
    }
    
    // Otherwise consider as a generic error
    return mapErrorToAppError(err)
}

func mapErrorToAppError(err error) *Error {
    switch err.Error() {
    case "sql: no rows in result set":
        return &Error{
            Code:    http.StatusNotFound,
            Message: "Recurso não encontrado",
        }
    default:
        return &Error{
            Code:    http.StatusInternalServerError,
            Message: err.Error(),
        }
    }
}
