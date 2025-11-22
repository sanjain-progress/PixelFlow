package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sanjain/pixelflow/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Did not connect:", err)
	}
	defer conn.Close()

	c := pb.NewAuthServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// 1. Register
	fmt.Println("--- Registering User ---")
	regRes, err := c.Register(ctx, &pb.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		log.Println("Register failed:", err)
	} else {
		fmt.Printf("Registered User ID: %s\n", regRes.UserId)
	}

	// 2. Login
	fmt.Println("\n--- Logging In ---")
	loginRes, err := c.Login(ctx, &pb.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		log.Fatalln("Login failed:", err)
	}
	fmt.Printf("Token: %s\n", loginRes.Token)

	// 3. Validate
	fmt.Println("\n--- Validating Token ---")
	valRes, err := c.Validate(ctx, &pb.ValidateRequest{
		Token: loginRes.Token,
	})
	if err != nil {
		log.Fatalln("Validate failed:", err)
	}
	fmt.Printf("Validated User ID: %s\n", valRes.UserId)
}
