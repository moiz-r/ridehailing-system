package service

import (
	"context"
	"strconv"

	"github.com/moiz-r/ridehailing-system/booking-service/internal/model"
	"github.com/moiz-r/ridehailing-system/booking-service/internal/store"
	pb "github.com/moiz-r/ridehailing-system/common/bookings"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type BookingService struct {
	pb.UnimplementedBookingServiceServer
	bookingStore store.BookingStore
	log          *zap.Logger
}

func NewBookingService(bookingStore store.BookingStore, logger *zap.Logger) *BookingService {
	return &BookingService{
		bookingStore: bookingStore,
		log:          logger,
	}
}

func (s *BookingService) CreateBooking(ctx context.Context, req *pb.CreateBookingRequest) (*pb.CreateBookingResponse, error) {
	// Implementation for creating a booking
	ride := model.Ride{
		Source:      req.Ride.Source,
		Destination: req.Ride.Destination,
		Distance:    req.Ride.Distance,
		Cost:        req.Ride.Cost,
	}
	userId, err := strconv.Atoi(req.UserId)
	if err != nil {
		s.log.Error("Error converting user id", zap.String("user_id", req.UserId), zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "Invalid user ID")
	}
	booking, err := s.bookingStore.CreateBooking(userId, ride)

	if err != nil {
		s.log.Error("Error creating booking", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Error creating booking")
	}

	return &pb.CreateBookingResponse{
		UserId:    strconv.Itoa(booking.UserID),
		RideId:    strconv.Itoa(booking.RideID),
		BookingId: strconv.Itoa(booking.ID),
	}, nil
}

func (s *BookingService) GetBooking(ctx context.Context, req *pb.GetBookingRequest) (*pb.GetBookingResponse, error) {
	// Implementation for retrieving a booking
	bookingId, err := strconv.Atoi(req.BookingId)
	if err != nil {
		s.log.Error("Error converting booking id", zap.String("booking_id", req.BookingId), zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "Invalid booking ID")
	}
	booking, err := s.bookingStore.GetBooking(bookingId)
	if err != nil {
		s.log.Error("Error retrieving booking", zap.String("booking_id", req.BookingId), zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "Booking not found")
	}

	return &pb.GetBookingResponse{
		UserId:      strconv.Itoa(booking.UserID),
		RideId:      strconv.Itoa(booking.RideID),
		BookingId:   strconv.Itoa(booking.BookingID),
		Source:      booking.Source,
		Destination: booking.Destination,
		Distance:    booking.Distance,
		Cost:        booking.Cost,
		Time:        timestamppb.New(booking.Time),
	}, nil
}
