package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Aderon1333/TPES/internal/models"
	"github.com/Aderon1333/TPES/pkg/utils/logfacade"
)

// Интерфейс должен быть объявлен в месте его использования, а не реализации
type TaskManagerService interface {
	PutTaskInDB(ctx context.Context, task *models.Task, logger *logfacade.LogFacade) error
	GetTaskFromDB(ctx context.Context, id int, logger *logfacade.LogFacade) (*models.Task, error)
}

func (h *Handler) getTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.l.Error("Incorrect id type")
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.tm.GetTaskFromDB(context.TODO(), id, h.l)

	if err != nil {
		h.l.Error(err)
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *Handler) getTaskKafka(c *gin.Context) {
	h.Prod.PlaceReq(c)
	res := h.Cons.GetReq(c)

	var id int

	errD := json.Unmarshal(res, &id)
	if errD != nil {
		h.l.Error("Cannot get ID from request")
		NewErrorResponse(c, http.StatusBadRequest, errD.Error())
		return
	} else {
		h.l.Info("ID ", id)
	}

	task, err := h.tm.GetTaskFromDB(context.TODO(), id, h.l)

	if err != nil {
		h.l.Error(err)
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *Handler) postTask(c *gin.Context) {
	var task models.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		h.l.Error(err)
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	h.l.Info(task)
	err := h.tm.PutTaskInDB(context.TODO(), &task, h.l)
	if err != nil {
		h.l.Error(err)
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	h.l.Info("Task added successfully!")
	c.JSON(http.StatusOK, gin.H{
		"message": "Task added successfully!",
		"task":    task,
	})
}

func (h *Handler) postTaskKafka(c *gin.Context) {
	h.Prod.PlaceReq(c)
	res := h.Cons.GetReq(c)

	var taskDecode models.Task

	errD := json.Unmarshal(res, &taskDecode)
	if errD != nil {
		h.l.Error("Cannot get task from request")
		NewErrorResponse(c, http.StatusBadRequest, errD.Error())
	} else {
		h.l.Info("task ", taskDecode)
	}

	if err := c.ShouldBindJSON(&taskDecode); err != nil {
		h.l.Error(err)
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	h.l.Info(taskDecode)
	err := h.tm.PutTaskInDB(context.TODO(), &taskDecode, h.l)
	if err != nil {
		h.l.Error(err)
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	h.l.Info("Task added successfully!")
	c.JSON(http.StatusOK, gin.H{
		"message": "Task added successfully!",
		"task":    taskDecode,
	})
}
