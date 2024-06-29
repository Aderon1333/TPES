package logmiddleware

import (
	"github.com/gin-gonic/gin"

	"github.com/Aderon1333/TPES/pkg/utils/logfacade"
)

func LoggerMiddleware(logger *logfacade.LogFacade) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("logger", logger)
		c.Next()
	}
}
