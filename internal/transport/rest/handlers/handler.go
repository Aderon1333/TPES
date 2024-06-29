package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/Aderon1333/TPES/pkg/middleware/gzipmiddleware"
	"github.com/Aderon1333/TPES/pkg/middleware/logmiddleware"
	"github.com/Aderon1333/TPES/pkg/utils/logfacade"
)

type Handler struct {
	tm TaskManagerInterfaceDB
}

func NewTaskHandler(tmi TaskManagerInterfaceDB) *Handler {
	return &Handler{tm: tmi}
}

func (h *Handler) InitRoutes(logger *logfacade.LogFacade) *gin.Engine {
	router := gin.New()

	// миддлвара для логгера
	router.Use(logmiddleware.LoggerMiddleware(logger))
	// миддлвара для gzip
	router.Use(gzipmiddleware.GzipMiddleware())

	router.GET("/tasks/:id", h.getTask) // получаем по id
	router.POST("/tasks", h.postTask)

	return router
}
