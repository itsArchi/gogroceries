package helper

import (
	"gogroceries/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendSuccess(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, domain.Response{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

func SendError(c *gin.Context, statusCode int, message string, errors interface{}) {
	c.JSON(statusCode, domain.Response{
		Status:  false,
		Message: message,
		Errors:  errors,
	})
}

func SendPagination(c *gin.Context, message string, page, limit int, data interface{}) {
	c.JSON(http.StatusOK, domain.Response{
		Status:  true,
		Message: message,
		Data: domain.PaginationResponse{
			Page:  page,
			Limit: limit,
			Data:  data,
		},
	})
}