package main

import (
	`context`
	`log`
	`net`
	`os`
	
	`github.com/Ndarz1/go-grpc-be/internal/handler`
	`github.com/Ndarz1/go-grpc-be/internal/repository`
	`github.com/Ndarz1/go-grpc-be/internal/service`
	`github.com/Ndarz1/go-grpc-be/pb/auth`
	`github.com/Ndarz1/go-grpc-be/pkg/database`
	`github.com/Ndarz1/go-grpc-be/pkg/grpcmiddleware`
	`github.com/joho/godotenv`
	`google.golang.org/grpc`
	`google.golang.org/grpc/reflection`
)

func main() {
	ctx := context.Background()
	
	godotenv.Load()
	
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Panicf("failed to listen: %v", err)
	}
	
	db := database.ConnectDB(ctx, os.Getenv("DB_URI"))
	log.Println("Connected to database")
	
	authRepository := repository.NewAuthRepository(db)
	
	authService := service.NewAuthService(authRepository)
	authHandler := handler.NewAuthHandler(authService)
	
	serv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcmiddleware.ErrorMiddleware,
		),
	)
	
	auth.RegisterAuthServiceServer(serv, authHandler)
	
	if os.Getenv("ENVIRONMENT") == "dev" {
		reflection.Register(serv)
		log.Println("Reflection is registered")
	}
	
	log.Println("Server is running on :50051 port")
	if err := serv.Serve(lis); err != nil {
		log.Panicf("failed to serve: %v", err)
	}
}
