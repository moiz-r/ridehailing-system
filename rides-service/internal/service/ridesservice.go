package service

import (
	"context"
	"strconv"

	pb "github.com/moiz-r/ridehailing-system/common/rides"
	"github.com/moiz-r/ridehailing-system/rides-service/internal/model"
	"github.com/moiz-r/ridehailing-system/rides-service/internal/store"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RidesService struct {
	pb.UnimplementedRidesServiceServer
	ridesStore store.RidesStore
	log        *zap.Logger
}

func NewRidesService(ridesStore store.RidesStore, logger *zap.Logger) *RidesService {
	return &RidesService{
		ridesStore: ridesStore,
		log:        logger,
	}
}

func (s *RidesService) UpdateRides(ctx context.Context, req *pb.UpdateRideRequest) (*pb.UpdateRideResponse, error) {
	// Implementation for creating a booking

	rideId, err := strconv.Atoi(req.RideId)
	ride := model.Ride{
		ID:          rideId,
		Source:      req.Ride.Source,
		Destination: req.Ride.Destination,
		Distance:    req.Ride.Distance,
		Cost:        req.Ride.Cost,
	}
	if err != nil {
		s.log.Error("Error converting ride id", zap.String("ride_id", req.RideId), zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "Invalid ride ID")
	}
	rideID, err := s.ridesStore.UpdateRide(&ride)

	if err != nil {
		s.log.Error("Error updating ride", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Error updating ride")
	}

	return &pb.UpdateRideResponse{
		RideId: strconv.Itoa(rideID),
	}, nil
}
