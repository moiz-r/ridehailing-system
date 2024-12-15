package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	pb "github.com/moiz-r/ridehailing-system/common/rides"
	"github.com/moiz-r/ridehailing-system/rides-service/configs"
	"github.com/moiz-r/ridehailing-system/rides-service/internal/service"
	"github.com/moiz-r/ridehailing-system/rides-service/internal/store"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":"+cfg.MetricsPort, nil); err != nil {
			log.Fatal("Failed to start metrics server", zap.Error(err))
		}
	}()

	ridesStore, err := store.NewPostgresStore(&cfg.Database)
	if err != nil {
		log.Fatal("Error creating rides store", zap.Error(err))
	}
	if err := ridesStore.Init(); err != nil {
		log.Fatal("Error initializing rides store", zap.Error(err))
	}

	ridesService := service.NewRidesService(ridesStore, log)

	// Start gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterRidesServiceServer(grpcServer, ridesService)

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
	log.Info("Shutting down Rides gRPC server...")
	grpcServer.GracefulStop()

	log.Info("Rides gRPC server stopped")

}
