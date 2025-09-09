package errs

import (
	"net/http"
)

type Error struct {
    Code    int `json:"code"`
    Message string
}

func (err *Error) Error() string {
    return err.Message
}

func From(err interface{}) *Error {
    // Check if is already an *Error
    if e, ok := err.(*Error); ok {
        return e
    }
    
    // Otherwise consider as a generic error
    if e, ok := err.(error); ok {
        return mapErrorToAppError(e)
    }

    return &Error{
        Code:    http.StatusInternalServerError,
        Message: "Erro Desconhecido",
    }
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