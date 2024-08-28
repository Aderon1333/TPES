package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/Aderon1333/TPES/internal/broker/kafka/consumer"
	"github.com/Aderon1333/TPES/internal/broker/kafka/producer"
	"github.com/Aderon1333/TPES/pkg/middleware/gzipmiddleware"
	"github.com/Aderon1333/TPES/pkg/utils/logfacade"
)

type Handler struct {
	tm   TaskManagerService
	um   UserManagerInterface
	Prod producer.ProducerKafka
	Cons consumer.ConsumerKafka
	l    *logfacade.LogFacade
}

func NewHandler(tmi TaskManagerService, umi UserManagerInterface, Prod producer.ProducerKafka, Cons consumer.ConsumerKafka, logger *logfacade.LogFacade) *Handler {
	return &Handler{
		tm:   tmi,
		um:   umi,
		Prod: Prod,
		Cons: Cons,
		l:    logger,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	// миддлвара для gzip
	router.Use(gzipmiddleware.GzipMiddleware())

	// user endpoints
	router.POST("users/signup", h.registerUser)
	router.POST("users/login", h.loginUser)
	router.POST("users/delete", h.deleteUser)

	// tasks endpoints
	// PG
	router.GET("/tasks/:id", h.validateUser, h.getTask) // получаем по id
	router.POST("/tasks", h.validateUser, h.postTask)

	// PG Kafka
	router.GET("/kafka/tasks/:id", h.getTaskKafka) // получаем по id
	router.POST("/kafka/tasks", h.postTaskKafka)

	// MG
	router.GET("/mongodb/tasks/:id", h.getTask)
	router.POST("/mongodb/tasks", h.postTask)

	return router
}
