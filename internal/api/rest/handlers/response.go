package handlers

import "github.com/gin-gonic/gin"

type customError struct {
	Message string `json:"message"`
}

func NewErrorResponse(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, customError{Message: message})
}
