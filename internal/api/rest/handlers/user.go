package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Aderon1333/TPES/internal/models"
	"github.com/Aderon1333/TPES/pkg/utils/logfacade"
)

type UserManagerInterface interface {
	Register(ctx *gin.Context, user *models.User, l *logfacade.LogFacade) error
	Login(ctx *gin.Context, user *models.User, l *logfacade.LogFacade) error
	Delete(ctx *gin.Context, user *models.User, l *logfacade.LogFacade) error
	Validate(ctx *gin.Context, user *models.User, l *logfacade.LogFacade) error
}

func (h *Handler) registerUser(c *gin.Context) {
	var user models.User

	err := c.Bind(&user)
	if err != nil {
		h.l.Error("Incorrect user data")
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.um.Register(c, &user, h.l)
	if err != nil {
		h.l.Error("Error occured while registering new user")
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) loginUser(c *gin.Context) {
	var user models.User

	err := c.Bind(&user)
	if err != nil {
		h.l.Error("Incorrect user data, ", err)
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.um.Login(c, &user, h.l)
	if err != nil {
		h.l.Error("Error occured while loging in, ", err)
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) deleteUser(c *gin.Context) {
	var user models.User

	err := c.Bind(&user)
	if err != nil {
		h.l.Error("Incorrect user data")
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.um.Delete(c, &user, h.l)
	if err != nil {
		h.l.Error("Error occured while deleting user")
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) validateUser(c *gin.Context) {
	var user models.User
	err := h.um.Validate(c, &user, h.l)
	if err != nil {
		h.l.Error(err)
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set("user", &user)
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged in successfully",
	})
}
