package responses

import (
	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/shared/errs"
)

type ApiJSONResponse struct {
    Code      int         `json:"code"`
    Message   string      `json:"message"`
    IsError   bool        `json:"isError"`
    IsSuccess bool        `json:"isSuccess"`
    Data      interface{} `json:"data"`
    Meta      interface{} `json:"meta"`
}

func JSONReturn(code int, message string, payload interface{}) *ApiJSONResponse {
    isSuccess := code >= 200 && code < 400
    isError := !isSuccess
	var msg string
	if message != "" {
		msg = message
	} else if isSuccess {
		msg = "Sucesso"
	} else {
        msg = "Erro"
    }
	return &ApiJSONResponse{
		Code:      code,
		Message:   msg,
		IsError:   isError,
		IsSuccess: isSuccess,
		Data:      payload,
		Meta:      nil,
	}
}

func JSON(c *gin.Context, code int, message string, payload interface{}) {
    response := JSONReturn(code, message, payload)
	c.JSON(response.Code, response)
}

func Error(c *gin.Context, err error) {
    apiError := errs.From(err)
    JSON(c, apiError.Code, apiError.Message, nil)
}
