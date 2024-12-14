package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	pb "github.com/moiz-r/ridehailing-system/common/users"
	"github.com/moiz-r/ridehailing-system/user-service/configs"
	"github.com/moiz-r/ridehailing-system/user-service/internal/service"
	"github.com/moiz-r/ridehailing-system/user-service/internal/store"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	log, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Failed to create logger", zap.Error(err))
	}
	defer log.Sync() // flushes buffer, if any

	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config", zap.Error(err))
	}
	userStore, err := store.NewPostgresStore(&cfg.Database)
	if err != nil {
		log.Fatal("Error creating user store", zap.Error(err))
	}
	if err := userStore.Init(); err != nil {
		log.Fatal("Error initializing user store", zap.Error(err))
	}
	userService := service.NewUserService(userStore, log)

	// Start gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userService)

	_, err = strconv.Atoi(cfg.GRPCPort)
	if err != nil {
		log.Fatal("Failed to parse gRPC port", zap.Error(err))
	}

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatal("Failed to listen", zap.Error(err))
	}
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("Failed to serve gRPC server", zap.Error(err))
		}
	}()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("Failed to serve gRPC server", zap.Error(err))
		}
	}()

	<-ctx.Done()
	log.Info("Shutting down gRPC server...")
	grpcServer.GracefulStop()
	log.Info("gRPC server stopped")

}
