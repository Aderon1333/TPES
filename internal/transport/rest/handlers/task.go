package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Aderon1333/TPES/internal/core"
	"github.com/Aderon1333/TPES/pkg/utils/logfacade"
)

// интерфейс должен быть объявлен в месте его использования, а не реализации
type TaskManagerInterfaceDB interface {
	PutTaskInDB(ctx context.Context, task *core.Task, logger *logfacade.LogFacade) error
	GetTaskFromDB(ctx context.Context, id int, logger *logfacade.LogFacade) (*core.Task, error)
}

func checkLogger(c *gin.Context) (*logfacade.LogFacade, bool) {
	logger, exists := c.Get("logger")
	if !exists {
		return nil, false
	}

	log, ok := logger.(*logfacade.LogFacade)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid logger"})
		return nil, false
	}

	return log, true
}

func (h *Handler) getTask(c *gin.Context) {
	logger, exists := checkLogger(c)
	if !exists {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error("incorrect id type")
		return
	}

	// Вызов метода менеджера
	task, err := h.tm.GetTaskFromDB(context.TODO(), id, logger)

	if err != nil {
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *Handler) postTask(c *gin.Context) {
	var task core.Task

	logger, exists := checkLogger(c)
	if !exists {
		return
	}

	// Привязываем JSON из тела запроса к структуре User
	if err := c.ShouldBindJSON(&task); err != nil {
		// Если привязка не удалась, возвращаем ошибку
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Вызов метода менеджера
	err := h.tm.PutTaskInDB(context.TODO(), &task, logger)
	if err != nil {
		return
	}

	// Возвращаем успешный ответ с полученными данными
	c.JSON(http.StatusOK, gin.H{
		"message": "Task added successfully!",
		"task":    task,
	})
}
