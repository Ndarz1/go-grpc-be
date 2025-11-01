package service

import (
	`context`
	`time`
	
	`github.com/Ndarz1/go-grpc-be/internal/entity`
	`github.com/Ndarz1/go-grpc-be/internal/repository`
	`github.com/Ndarz1/go-grpc-be/internal/utils`
	`github.com/Ndarz1/go-grpc-be/pb/auth`
	`github.com/google/uuid`
	`golang.org/x/crypto/bcrypt`
)

type IAuthService interface {
	Register(ctx context.Context, request *auth.RegisterRequest) (
		*auth.RegisterResponse,
		error,
	)
}

type authService struct {
	authRepository repository.IAuthRepository
}

func (as *authService) Register(ctx context.Context, request *auth.RegisterRequest) (
	*auth.RegisterResponse,
	error,
) {
	// Check email to database
	if request.Password != request.PasswordConfirmation {
		return &auth.RegisterResponse{
			Base: utils.BadRequestResponse("Password is not match"),
		}, nil
	}
	user, err := as.authRepository.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}
	// If email register, error
	if user != nil {
		return &auth.RegisterResponse{
			Base: utils.BadRequestResponse("User already exists"),
		}, nil
	}
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), 10)
	if err != nil {
		return nil, err
	}
	// If email not register, insert
	newUser := entity.User{
		Id:        uuid.NewString(),
		FullName:  request.FullName,
		Email:     request.Email,
		Password:  string(hashedPassword),
		RoleCode:  entity.UserRoleCustomer,
		CreatedAt: time.Now(),
		CreatedBy: &request.FullName,
	}
	err = as.authRepository.InsertUser(ctx, &newUser)
	if err != nil {
		return nil, err
	}
	
	return &auth.RegisterResponse{
		Base: utils.SuccessResponse("User successfully registered"),
	}, nil
	
}

func NewAuthService(authRepository repository.IAuthRepository) IAuthService {
	return &authService{
		authRepository: authRepository,
	}
}
