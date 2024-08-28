package handlers

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// Сгенерированный код
	pb "github.com/Aderon1333/TPES/internal/api/grpc/gen/go/tpes"
	"github.com/Aderon1333/TPES/internal/models"
	"github.com/Aderon1333/TPES/pkg/utils/logfacade"
)

// Интерфейс на бизнес логику
type TaskManagerService interface {
	PutTaskInDB(ctx context.Context, task *models.Task, logger *logfacade.LogFacade) error
	GetTaskFromDB(ctx context.Context, id int, logger *logfacade.LogFacade) (*models.Task, error)
}

// GRPCService представляет собой сервис, предоставляющий API через gRPC
type GRPCService struct {
	pb.UnimplementedTaskHandlerServer
	tms TaskManagerService
}

// NewGRPCService создает новый экземпляр GRPCService
func NewGRPCService(taskManagerService TaskManagerService) *GRPCService {
	return &GRPCService{tms: taskManagerService}
}

// RegisterGRPCServer регистрирует все GRPC-методы в сервере
func (gs *GRPCService) RegisterGRPCServer(srv *grpc.Server) {
	pb.RegisterTaskHandlerServer(srv, gs)
}

func (gs *GRPCService) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.TaskResponse, error) {
	task, err := gs.tms.GetTaskFromDB(ctx, int(req.Id), &logfacade.LogFacade{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get task: %v", err)
	}

	return &pb.TaskResponse{Task: &pb.Task{
		Id:     task.ID,
		Status: task.Status,
		Item:   task.Item,
	}}, nil
}

func (gs *GRPCService) PutTask(ctx context.Context, req *pb.PutTaskRequest) (*pb.TaskResponse, error) {
	task := &models.Task{
		ID:     req.Task.Id,
		Status: req.Task.Status,
		Item:   req.Task.Item,
	}

	err := gs.tms.PutTaskInDB(ctx, task, &logfacade.LogFacade{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to put task: %v", err)
	}

	return &pb.TaskResponse{Task: &pb.Task{
		Id:     task.ID,
		Status: task.Status,
		Item:   task.Item,
	}}, nil
}
