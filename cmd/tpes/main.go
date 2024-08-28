package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	// "github.com/Aderon1333/TPES/internal/api/rest/handlers"
	"github.com/Aderon1333/TPES/internal/api/rest"
	"github.com/Aderon1333/TPES/internal/api/rest/handlers"
	"github.com/Aderon1333/TPES/internal/broker/kafka/consumer"
	"github.com/Aderon1333/TPES/internal/broker/kafka/producer"
	"github.com/Aderon1333/TPES/internal/config"
	"github.com/Aderon1333/TPES/internal/repository/mg"
	"github.com/Aderon1333/TPES/internal/repository/pg"
	"github.com/Aderon1333/TPES/internal/service/authentification"
	"github.com/Aderon1333/TPES/internal/service/manager"
	"github.com/Aderon1333/TPES/pkg/broker/kafka"
	"github.com/Aderon1333/TPES/pkg/repository/mongodb"
	"github.com/Aderon1333/TPES/pkg/repository/postgresql"
	"github.com/Aderon1333/TPES/pkg/utils/logfacade"
	"github.com/IBM/sarama"
)

func main() {
	// Получение конфиги (можно еще сделать через .env)
	cfg := config.GetConfig()

	// Логирование в отдельном утилитарном пакете (фасад над внешним логгером + возможность донастройки)
	logger := &logfacade.LogFacade{}
	logFile, err := os.OpenFile(cfg.AppCfg.Log, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	logrusLogger := logfacade.NewLogrusLogger(logFile)
	logger.SetLogger(logrusLogger)

	// REST API -------------------------------------------------------------------------------------
	var handler *handlers.Handler

	// Kafka
	addr := []string{"kafka:9092"} // в конфигу
	var prod sarama.SyncProducer
	var errP error
	for i := 0; i < 5; i++ {
		prod, errP = kafka.ConnectProducer(addr)
		if errP != nil {
			logger.Info(errP.Error())
			time.Sleep(3 * time.Second)
		}
	}

	if errP != nil {
		logger.Fatal(errP.Error())
		return
	}
	sProd := producer.NewProducer(prod)

	cons, errC := kafka.ConnectConsumer(addr)
	if errC != nil {
		logger.Fatal(errC)
		return
	}
	sCons := consumer.NewConsumer(cons)

	// Db
	if cfg.AppCfg.DB == "postgresql" {
		postgreSQLClient, err := postgresql.NewClient(context.TODO(), cfg.PostgresCfg)
		if err != nil {
			logger.Error(err)
		} else {
			logger.Info("Successfully connected to PostgreSQL database")
		}

		repository := pg.NewRepository(postgreSQLClient)
		userManager := authentification.NewUserManager(postgreSQLClient, cfg.JWTCfg.Secret, cfg.JWTCfg.AccessCookie, cfg.JWTCfg.RefreshCookie) // сервис авторизации
		taskManager := manager.NewTaskManager(repository)                                                                                      // сервис
		handler = handlers.NewHandler(taskManager, userManager, sProd, sCons, logger)
	} else if cfg.AppCfg.DB == "mongodb" {
		mongoClient, err := mongodb.NewMongoClient(context.TODO(), cfg.MongoCfg, logger)
		if err != nil {
			logger.Error(err)
		}

		urlDAO := mg.NewUrlDAO(context.TODO(), mongoClient)
		mongoRepository := mg.NewRepository(urlDAO)
		mongoManager := manager.NewTaskManager(mongoRepository)
		handler = handlers.NewHandler(mongoManager, nil, sProd, sCons, logger)

		// defer mongodb.DisconnectMongoClient(context.TODO(), mongoClient)
		// Создание контекста для закрытия сервера (а как теперь дисконнектить?)
		// err = DisconnectMongoClient(context.TODO(), client)
		// if err != nil {
		// 	logger.Error(err)
		// }
	} else {
		handler = nil
		logger.Error("Database type is not available")
	}

	// Создание сервера
	restHttpServer := rest.Server{}

	logger.Info("Starting server on port ", cfg.AppCfg.Port)
	// Запуск сервера
	if handler != nil {
		err := restHttpServer.RunHTTPServer(cfg.AppCfg.IP, cfg.AppCfg.Port, handler.InitRoutes())

		if err != nil {
			logger.Error(err)
		}
	}

	// Gracefull shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	// Закрытие сервера
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := restHttpServer.ShutdownHTTPServer(ctx); err != nil {
		logger.Error("Server Shutdown Failed:", err)
	} else {
		logger.Error("Server Exited Properly")
	}

	// gRPC -----------------------------------------------------------------------------------------
	// var grpcHandler *handlers.GRPCService

	// if cfg.AppCfg.DB == "postgresql" {
	// 	postgreSQLClient, err := postgresql.NewClient(context.TODO(), cfg.PostgresCfg)
	// 	if err != nil {
	// 		logger.Error(err)
	// 	} else {
	// 		logger.Info("Successfully connected to PostgreSQL database")
	// 	}

	// 	repository := pg.NewRepository(postgreSQLClient)
	// 	//userManager := authentification.NewUserManager(postgreSQLClient, cfg.JWTCfg.Secret, cfg.JWTCfg.AccessCookie, cfg.JWTCfg.RefreshCookie) // сервис авторизации
	// 	taskManager := manager.NewTaskManager(repository)  // сервис задач
	// 	grpcHandler = handlers.NewGRPCService(taskManager) // grpc обработчик
	// }

	// // Запуск gRPC сервера
	// lis, err := net.Listen("tcp", ":8080")
	// if err != nil {
	// 	log.Fatalf("Failed to listen: %v", err)
	// }

	// // Создание gRPC сервера
	// grpcServer := grpc.NewServer()

	// // Регистрация обработчиков с gRPC сервером
	// grpcHandler.RegisterGRPCServer(grpcServer)

	// log.Printf("gRPC server listening on port 8080")

	// if err := grpcServer.Serve(lis); err != nil {
	// 	log.Fatalf("Failed to serve: %v", err)
	// }
	// ----------------------------------------------------------------------------------------------
}
