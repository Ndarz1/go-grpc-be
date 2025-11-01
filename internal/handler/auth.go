package handler

import (
	`context`
	
	`github.com/Ndarz1/go-grpc-be/internal/service`
	`github.com/Ndarz1/go-grpc-be/internal/utils`
	`github.com/Ndarz1/go-grpc-be/pb/auth`
)

type authHandler struct {
	auth.UnimplementedAuthServiceServer
	
	authService service.IAuthService
}

func (sh *authHandler) Register(ctx context.Context, request *auth.RegisterRequest) (
		*auth.RegisterResponse,
		error,
) {
	validationErrors, err := utils.CheckValidation(request)
	if err != nil {
		return nil, err
	}
	if validationErrors != nil {
		return &auth.RegisterResponse{
			Base: utils.ValidationErrorResponse(validationErrors),
		}, nil
	}
	
	res, err := sh.authService.Register(ctx, request)
	if err != nil {
		return nil, err
	}
	
	return res, nil
}

func NewAuthHandler(authService service.IAuthService) *authHandler {
	return &authHandler{
		authService: authService,
	}
}
