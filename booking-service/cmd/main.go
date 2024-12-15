package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/moiz-r/ridehailing-system/booking-service/configs"
	"github.com/moiz-r/ridehailing-system/booking-service/internal/service"
	"github.com/moiz-r/ridehailing-system/booking-service/internal/store"
	pb "github.com/moiz-r/ridehailing-system/common/bookings"
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

	bookingStore, err := store.NewPostgresStore(&cfg.Database)
	if err != nil {
		log.Fatal("Error creating user store", zap.Error(err))
	}
	if err := bookingStore.Init(); err != nil {
		log.Fatal("Error initializing user store", zap.Error(err))
	}
	bookingService := service.NewBookingService(bookingStore, log)

	grpcServer := grpc.NewServer()
	pb.RegisterBookingServiceServer(grpcServer, bookingService)

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
