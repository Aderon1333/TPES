package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Aderon1333/TPES/internal/config"
	"github.com/Aderon1333/TPES/internal/repository/mdb"
	"github.com/Aderon1333/TPES/internal/service/manager"
	"github.com/Aderon1333/TPES/internal/transport/rest"
	"github.com/Aderon1333/TPES/internal/transport/rest/handlers"
	"github.com/Aderon1333/TPES/pkg/repository/mongodb"
	"github.com/Aderon1333/TPES/pkg/utils/logfacade"
)

func main() {
	// Логирование
	// 1) логирование в отдельный утилитарный пакет (фасад над внешним логгером + возможность донастройки)
	logger := &logfacade.LogFacade{}
	logrusLogger := logfacade.NewLogrusLogger()
	logger.SetLogger(logrusLogger)

	// Получение конфиги
	cfg := config.GetConfig(logger)

	// Подключение postgresql
	// postgreSQLClient, err := postgresql.NewClient(context.TODO(), cfg.Storage, logger)
	// if err != nil {
	// 	logger.Error(err)
	// }

	// Подключение mongodb
	mongoClient, err := mongodb.NewMongoClient(context.TODO(), cfg.Mongo, logger)
	if err != nil {
		logger.Error(err)
	}

	urlDAO := mdb.NewUrlDAO(context.TODO(), mongoClient)

	// repository := db.NewRepository(postgreSQLClient)
	mongoRepository := mdb.NewRepository(urlDAO)

	// pgManager := manager.NewTaskManagerDB(repository) // сервис
	// сервис инициализированный mongo(подвел монго под такой же интерфейс в менеджере)
	mongoManager := manager.NewTaskManagerDB(mongoRepository)

	// как быть с хендлером в этом случае?
	handler := handlers.NewTaskHandler(mongoManager) // htpp handler

	// Создание сервера
	restHttpServer := rest.Server{}

	// Запуск сервера
	err = restHttpServer.RunHTTPServer(cfg.App.IP, cfg.App.Port, handler.InitRoutes(logger))
	//err = restHttpServer.RunHTTPServer(cfg.App.IP, cfg.App.Port, handler.InitRoutes(logger))
	if err != nil {
		logger.Error(err)
	}

	// Gracefull shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	// Закрытие сервера
	// Создание контекста для закрытия сервера
	err = mongodb.DisconnectMongoClient(context.TODO(), mongoClient)
	if err != nil {
		logger.Error(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := restHttpServer.ShutdownHTTPServer(ctx); err != nil {
		logger.Error("Server Shutdown Failed:", err)
	} else {
		logger.Error("Server Exited Properly")
	}
}
