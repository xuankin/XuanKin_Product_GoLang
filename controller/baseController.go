package controller

import (
	"Product_Mangement_Api/models"
	"github.com/gin-gonic/gin"
)

func SendResponse(c *gin.Context, success bool, message string, data interface{}, statusCode int) {
	c.JSON(statusCode, models.APIResponse{
		Success: success,
		Message: message,
		Data:    data,
	})
}
func SendError(c *gin.Context, message string, data interface{}, statusCode int) {
	c.JSON(statusCode, models.APIResponse{
		Success: false,
		Message: message,
		Errors:  data,
	})
}
