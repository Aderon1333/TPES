package authentification

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/Aderon1333/TPES/internal/api/rest/handlers"
	"github.com/Aderon1333/TPES/internal/models"
	"github.com/Aderon1333/TPES/pkg/utils/logfacade"
	"github.com/Aderon1333/TPES/pkg/utils/token"
)

type UserManager struct {
	db            *gorm.DB
	secret        string
	accessCookie  string
	refreshCookie string
}

func NewUserManager(db *gorm.DB, secretJWTKey string, accessCookie string, refreshCookie string) handlers.UserManagerInterface {
	return &UserManager{
		db:            db,
		secret:        secretJWTKey,
		accessCookie:  accessCookie,
		refreshCookie: refreshCookie,
	}
}

func (um *UserManager) Register(c *gin.Context, user *models.User, l *logfacade.LogFacade) error {
	var foundUser models.User
	result := um.db.First(&foundUser, "login = ?", user.Login)

	if result.Error != nil {
		hashedPassword, err := token.HashPassword(user.Password)
		if err != nil {
			l.Error(err)
			return err
		}

		newUser := models.User{Login: user.Login, Password: hashedPassword}

		result := um.db.Create(&newUser)

		if result.Error != nil {
			l.Error(result.Error)
			return result.Error
		}
	} else {
		l.Error("User with this login already exists")
		return errors.New("User with this login already exists")
	}

	return nil
}

func (um *UserManager) Login(c *gin.Context, user *models.User, l *logfacade.LogFacade) error {
	var foundUser models.User

	result := um.db.First(&foundUser, "login = ?", user.Login)

	if result.Error != nil {
		l.Error(result.Error)
		return result.Error
	}

	err := token.VerifyPassword(foundUser.Password, user.Password)
	if err != nil {
		l.Error(result.Error)
		return err
	}

	access_token, err := token.GenerateToken(time.Minute*3, foundUser.Login, um.secret)
	if err != nil {
		l.Error(result.Error)
		return err
	}

	refresh_token, err := token.GenerateToken(time.Hour*24, foundUser.Login, um.secret)
	if err != nil {
		l.Error(result.Error)
		return err
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(um.accessCookie, access_token, 180, "", "", false, true) // время жизни токена тоже в конфигу
	c.SetCookie(um.refreshCookie, refresh_token, int(time.Hour*24), "", "", false, true)

	return nil
}

func (um *UserManager) Delete(c *gin.Context, user *models.User, l *logfacade.LogFacade) error {
	result := um.db.Where("login = ?", user.Login).Delete(&user)

	if result.Error != nil {
		l.Error(result.Error)
		return result.Error
	}

	return nil
}

func (um *UserManager) Validate(c *gin.Context, user *models.User, l *logfacade.LogFacade) error {
	tokenString, err := c.Cookie(um.accessCookie)
	if err != nil {
		l.Error(err)
		return err
	}

	access_claims, err := token.ValidateToken(tokenString, um.secret)
	if err != nil {
		l.Error(err)
		return err
	}

	foundUser := um.db.First(user, "login = ?", access_claims["sub"])
	if foundUser.Error != nil {
		l.Error(foundUser.Error)
		return foundUser.Error
	}

	if float64(time.Now().Unix()) > access_claims["exp"].(float64) {
		refreshString, err := c.Cookie(um.refreshCookie)
		if err != nil {
			l.Error(err)
			return err
		}

		refresh_claims, err := token.ValidateToken(refreshString, um.secret)
		if err != nil {
			l.Error(err)
			return err
		}

		if float64(time.Now().Unix()) > refresh_claims["exp"].(float64) {
			l.Error("Refresh token expired")
			return errors.New("Refresh token expired")
		}

		access_token, err := token.GenerateToken(time.Minute*3, user.Login, um.secret)
		if err != nil {
			l.Error(err)
			return err
		}

		refresh_token, err := token.GenerateToken(time.Hour*24, user.Login, um.secret)
		if err != nil {
			l.Error(err)
			return err
		}

		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie(um.accessCookie, access_token, 180, "", "", false, true) // время жизни токена тоже в конфигу
		c.SetCookie(um.refreshCookie, refresh_token, int(time.Hour*24), "", "", false, true)
	}

	return nil
}
