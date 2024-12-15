package service

import (
	"context"
	"strconv"

	pb "github.com/moiz-r/ridehailing-system/common/users"
	"github.com/moiz-r/ridehailing-system/user-service/internal/store"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserService is a struct that contains the methods to interact with the user service
type UserService struct {
	pb.UnimplementedUserServiceServer
	userStore store.UserStore
	log       *zap.Logger
}

// NewUserService creates a new UserService
func NewUserService(userStore store.UserStore, logger *zap.Logger) *UserService {
	return &UserService{
		userStore: userStore,
		log:       logger,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	s.log.Info("Creating new user", zap.String("username", req.Name))
	// Implementation for creating a user
	user, err := s.userStore.CreateUser(req.Name)

	if err != nil {
		s.log.Error("Error creating user", zap.String("username", req.Name), zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Error creating user")
	}

	return &pb.CreateUserResponse{
		UserId: strconv.Itoa(user.UserID),
	}, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	s.log.Info("Retrieving user", zap.String("user_id", req.UserId))
	// Implementation for retrieving a user
	userId, err := strconv.Atoi(req.UserId)
	if err != nil {
		s.log.Error("Error converting user id", zap.String("user_id", req.UserId), zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "Invalid user ID")
	}
	user, err := s.userStore.GetUser(userId)
	if err != nil {
		s.log.Error("Error retrieving user", zap.String("user_id", req.UserId), zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	return &pb.GetUserResponse{
		Name:   user.Name,
		UserId: strconv.Itoa(user.UserID),
	}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	s.log.Info("Deleting user", zap.String("user_id", req.UserId))
	// Implementation for deleting a user
	userId, err := strconv.Atoi(req.UserId)
	if err != nil {
		s.log.Error("Error converting user id", zap.String("user_id", req.UserId), zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "Invalid user ID")
	}
	err = s.userStore.DeleteUser(userId)
	if err != nil {
		s.log.Error("Error deleting user", zap.String("user_id", req.UserId), zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Error deleting user")
	}

	return &pb.DeleteUserResponse{
		Success: true,
	}, nil
}
