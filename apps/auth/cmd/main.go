package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/sanjain/pixelflow/apps/auth/internal/db"
	"github.com/sanjain/pixelflow/apps/auth/internal/models"
	"github.com/sanjain/pixelflow/apps/auth/internal/utils"
	"github.com/sanjain/pixelflow/pkg/pb"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedAuthServiceServer
	H *db.Handler
}

func (s *server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var user models.User

	if result := s.H.DB.Where(&models.User{Email: req.Email}).First(&user); result.Error == nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user.Email = req.Email
	user.Password = hashedPassword

	if result := s.H.DB.Create(&user); result.Error != nil {
		return nil, errors.New("failed to create user")
	}

	return &pb.RegisterResponse{
		UserId: fmt.Sprintf("%d", user.ID),
	}, nil
}

func (s *server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var user models.User

	if result := s.H.DB.Where(&models.User{Email: req.Email}).First(&user); result.Error != nil {
		return nil, errors.New("user not found")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &pb.LoginResponse{
		Token: token,
	}, nil
}

func (s *server) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	userID, err := utils.ValidateToken(req.Token)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	return &pb.ValidateResponse{
		UserId: fmt.Sprintf("%d", userID),
	}, nil
}

func main() {
	h := db.Init("postgres://pixelflow:password@localhost:5432/auth_db?sslmode=disable")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	fmt.Println("Auth Service running on :50051")

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, &server{H: h})

	if err := s.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
