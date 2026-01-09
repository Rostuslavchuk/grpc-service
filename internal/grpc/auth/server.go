package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"sso/internal/storage"

	ssov1 "github.com/Rostuslavchuk/sso-protos/gen/go/sso"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RequestValidateLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	AppID    int64  `json:"app_id" validate:"required,gt=0"`
}
type RequestValidateRegister struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}
type RequestValidateIsAdmin struct {
	UserID int64 `json:"user_id" validate:"required,gt=0"`
}

type Auth interface {
	Login(ctx context.Context, email string, password string, appID int64) (token string, error error)
	SaveUser(ctx context.Context, email string, password string) (userID int64, error error)
	IsAdmin(ctx context.Context, userID int64) (isAdmin bool, error error)
}
type ServerAPI struct {
	ssov1.UnimplementedAuthServer // реалізує методи Register, Login, IsAdmin, вони returns Unimplemented тобто нереалізований ssov1.UnimplementedAuthServer корисний тим шо при додаванні не треба тут дописувати
	auth                          Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &ServerAPI{auth: auth}) // якщо прийде запит, що стосується сервісу Auth (як описано в протофайлі), передавай його на обробку ось цій структурі ServerAPI
}

func (s *ServerAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	reqValidLogin := &RequestValidateLogin{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		AppID:    req.GetAppId(),
	}
	var errorsMsgs []string
	if err := validator.New().Struct(reqValidLogin); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, valErr := range validationErrors {
				switch valErr.ActualTag() {
				case "required":
					errorsMsgs = append(errorsMsgs, fmt.Sprintf("Field %s is required", valErr.Field()))
				case "email":
					errorsMsgs = append(errorsMsgs, fmt.Sprintf("Field %s is not valid", valErr.Field()))
				case "min":
					errorsMsgs = append(errorsMsgs, fmt.Sprintf("Field %s must be at least %s characters long", valErr.Field(), valErr.Param()))
				case "gt":
					errorsMsgs = append(errorsMsgs, fmt.Sprintf("Field %s must be greater then %s", valErr.Field(), valErr.Param()))
				}
			}
			return nil, status.Errorf(codes.InvalidArgument, "validation error: %s", strings.Join(errorsMsgs, ", "))
		}
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), req.GetAppId())
	if err != nil {
		if errors.Is(err, storage.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}

		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *ServerAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	reqValidRegister := &RequestValidateRegister{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	var errorsMsgs []string
	if err := validator.New().Struct(reqValidRegister); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, valErr := range validationErrors {
				switch valErr.ActualTag() {
				case "required":
					errorsMsgs = append(errorsMsgs, fmt.Sprintf("Field %s is required", valErr.Field()))
				case "email":
					errorsMsgs = append(errorsMsgs, fmt.Sprintf("Field %s is not valid", valErr.Field()))
				case "min":
					errorsMsgs = append(errorsMsgs, fmt.Sprintf("Field %s must be at least %s characters long", valErr.Field(), valErr.Param()))
				}
			}
			return nil, status.Errorf(codes.InvalidArgument, "validation error: %s", strings.Join(errorsMsgs, ", "))
		}
	}

	userID, err := s.auth.SaveUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &ssov1.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *ServerAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	reqValidIsAdmin := &RequestValidateIsAdmin{
		UserID: req.GetUserId(),
	}

	var errorsMsgs []string
	if err := validator.New().Struct(reqValidIsAdmin); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, valErr := range validationErrors {
				switch valErr.ActualTag() {
				case "required":
					errorsMsgs = append(errorsMsgs, fmt.Sprintf("Field %s is required", valErr.Field()))
				case "gt":
					errorsMsgs = append(errorsMsgs, fmt.Sprintf("Field %s must be greater then %s", valErr.Field(), valErr.Param()))
				}
			}
			return nil, status.Errorf(codes.InvalidArgument, "validation error: %s", strings.Join(errorsMsgs, ", "))
		}
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			return nil, status.Error(codes.NotFound, "user is not exists")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}
