package service

import (
	`context`
	`errors`
	`fmt`
	`os`
	`strings`
	`time`
	
	`github.com/Ndarz1/go-grpc-be/internal/entity`
	`github.com/Ndarz1/go-grpc-be/internal/repository`
	`github.com/Ndarz1/go-grpc-be/internal/utils`
	`github.com/Ndarz1/go-grpc-be/pb/auth`
	`github.com/golang-jwt/jwt/v5`
	`github.com/google/uuid`
	gocache `github.com/patrickmn/go-cache`
	`golang.org/x/crypto/bcrypt`
	`google.golang.org/grpc/codes`
	`google.golang.org/grpc/metadata`
	`google.golang.org/grpc/status`
)

type IAuthService interface {
	Register(ctx context.Context, request *auth.RegisterRequest) (
		*auth.RegisterResponse,
		error,
	)
	Login(ctx context.Context, request *auth.LoginRequest) (*auth.LoginResponse, error)
	Logout(ctx context.Context, request *auth.LogoutRequest) (*auth.LogoutResponse, error)
}

type authService struct {
	authRepository repository.IAuthRepository
	cacheService   *gocache.Cache
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

func (as *authService) Login(ctx context.Context, request *auth.LoginRequest) (*auth.LoginResponse, error) {
	// Check email
	user, err := as.authRepository.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return &auth.LoginResponse{
			Base: utils.BadRequestResponse("User is not registered"),
		}, nil
	}
	// Check password match
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
		}
		return nil, err
	}
	// Generate JWT
	now := time.Now()
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, entity.JwtClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   user.Id,
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24)),
				IssuedAt:  jwt.NewNumericDate(now),
			},
			Email:    user.Email,
			FullName: user.FullName,
			Role:     user.RoleCode,
		},
	)
	secretKey := os.Getenv("JWT_SECRET")
	accessToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}
	// Send response
	return &auth.LoginResponse{
		Base:        utils.SuccessResponse("Login successful"),
		AccessToken: accessToken,
	}, nil
}

func (as *authService) Logout(ctx context.Context, request *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	// Get token from metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Metadata not received")
	}
	bearerToken, ok := md["authorization"]
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Metadata not received")
	}
	
	if len(bearerToken) == 1 {
		return nil, status.Errorf(codes.Unauthenticated, "Metadata not received")
	}
	// Bearer
	tokenSplit := strings.Split(bearerToken[0], " ")
	if len(tokenSplit) != 2 {
		return nil, status.Errorf(codes.Unauthenticated, "Metadata not received")
	}
	if tokenSplit[0] != "Bearer" {
		return nil, status.Errorf(codes.Unauthenticated, "Metadata not received")
	}
	jwtToken := bearerToken[1]
	// Back token to entity jwt
	tokenClaims, err := jwt.ParseWithClaims(
		jwtToken, &entity.JwtClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		},
	)
	
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Metadata not received")
	}
	if !tokenClaims.Valid {
		return nil, status.Errorf(codes.Unauthenticated, "Metadata not received")
	}
	
	var claims *entity.JwtClaims
	if claims, ok = tokenClaims.Claims.(*entity.JwtClaims); !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Metadata not received")
	}
	
	// Insert token to cache
	as.cacheService.Set(jwtToken, "", time.Duration(claims.ExpiresAt.Time.Unix()-time.Now().Unix())*time.Second)
	
	// Send response
	return &auth.LogoutResponse{
		Base: utils.SuccessResponse("Logout successful"),
	}, nil
}

func NewAuthService(authRepository repository.IAuthRepository, cacheService *gocache.Cache) IAuthService {
	return &authService{
		authRepository: authRepository,
		cacheService:   cacheService,
	}
}
