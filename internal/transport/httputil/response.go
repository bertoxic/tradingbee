package httputil

import (

	"github.com/bertoxic/tradingbee/internal/models"
	"github.com/gin-gonic/gin"
)

type JsonResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorJson  `json:"error,omitempty"`
}

type ErrorJson struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}


func ReadJson(ctx *gin.Context) {

}

func WriteJson(ctx *gin.Context, success bool, statuscode int, response *models.JsonResponse) {
	var resp JsonResponse
	ctx.Request.Header.Set("Content-Type", "application/json")
	header := ctx.Writer.Header()
	header.Write(ctx.Writer)
	resp.Success = success
	var status int
	if !success {
		resp.Data = response.Data
		resp.Message = response.Message
		status = statuscode
		resp.Error = &ErrorJson{400, "error occured"}

	} else {	
		resp.Data = response.Data
		resp.Message = response.Message
		status = 200
	}

	ctx.JSON(status, response)
}

func Success(data interface{}, message string) *JsonResponse {
	return &JsonResponse{
		Success: true,
		Data:    data,
		Message: message,
	}
}

func (js *JsonResponse) Fail(statuscode int, err error) *JsonResponse {
	return &JsonResponse{
		Success: false,
		Error: &ErrorJson{
			Code:    statuscode,
			Message: err.Error(),
		},
	}
}
