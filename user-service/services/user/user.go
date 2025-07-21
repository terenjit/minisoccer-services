package services

import (
	"context"
	"strings"
	"time"
	"user-service/config"
	"user-service/constants"
	"user-service/domain/dto"
	"user-service/repositories"

	errWrap "user-service/common/error"
	errConstant "user-service/constants/error"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository repositories.IRegistryRepository
}

type IUserService interface {
	Login(context.Context, *dto.LoginRequest) (*dto.LoginResponse, error)
	Register(context.Context, *dto.RegisterRequest) (*dto.RegisterResponse, error)
	Update(context.Context, *dto.UpdateRequest, string) (*dto.UserResponse, error)
	GetUserLogin(context.Context) (*dto.UserResponse, error)
	GetUserByUUID(context.Context, string) (*dto.UserResponse, error)
}

type Claims struct {
	User *dto.UserResponse
	jwt.RegisteredClaims
}

func NewUserService(repository repositories.IRegistryRepository) IUserService {
	return &UserService{repository: repository}
}

func (s *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.repository.GetUser().FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, err
	}

	TokenExpireTime := time.Now().Add(time.Duration(config.Cfg.JWTExpirationTime) * time.Minute).Unix()
	data := &dto.UserResponse{
		UUID:        user.UUID,
		Name:        user.Name,
		Username:    user.Username,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Role:        strings.ToLower(user.Role.Code),
	}

	claims := &Claims{
		User: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(TokenExpireTime, 0)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tokenString, err := token.SignedString([]byte(config.Cfg.JWTSecretKey))
	if err != nil {
		return nil, err
	}

	response := &dto.LoginResponse{
		User:  *data,
		Token: tokenString,
	}

	return response, nil
}

func (s *UserService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	checkEmail, err := s.repository.GetUser().FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if checkEmail != nil {
		return nil, errWrap.WrapError(errConstant.ErrEmailExists)
	}

	checkUsername, err := s.repository.GetUser().FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if checkUsername != nil {
		return nil, errWrap.WrapError(errConstant.ErrUsernameExists)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if req.Password != req.ConfirmPass {
		return nil, errWrap.WrapError(errConstant.ErrPasswordDoesNotMatch)
	}

	req.Password = string(hashedPassword)
	req.RoleID = constants.Customer
	createUser, err := s.repository.GetUser().Register(ctx, req)
	if err != nil {
		return nil, err
	}

	response := &dto.RegisterResponse{
		User: dto.UserResponse{
			UUID:        createUser.UUID,
			Name:        createUser.Name,
			Username:    createUser.Username,
			PhoneNumber: createUser.PhoneNumber,
			Email:       createUser.Email,
		},
	}

	return response, nil
}
func (s *UserService) Update(ctx context.Context, req *dto.UpdateRequest, uuid string) (*dto.UserResponse, error) {

	getUser, err := s.repository.GetUser().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if req.Username != getUser.Username {
		checkUsername, err := s.repository.GetUser().FindByUsername(ctx, req.Username)
		if err != nil {
			return nil, err
		}
		if checkUsername != nil {
			return nil, errWrap.WrapError(errConstant.ErrUsernameExists)
		}
	}

	if req.Email != getUser.Email {
		checkEmail, err := s.repository.GetUser().FindByEmail(ctx, req.Email)
		if err != nil {
			return nil, err
		}
		if checkEmail != nil {
			return nil, errWrap.WrapError(errConstant.ErrEmailExists)
		}
	}

	if req.Password != "" {
		if req.Password != req.ConfirmPass {
			return nil, errWrap.WrapError(errConstant.ErrPasswordDoesNotMatch)
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		req.Password = string(hashedPassword)
	}

	req.RoleID = constants.Customer
	Update, err := s.repository.GetUser().Update(ctx, req, uuid)
	if err != nil {
		return nil, err
	}

	response := &dto.UserResponse{
		UUID:        Update.UUID,
		Name:        Update.Name,
		Username:    Update.Username,
		Email:       Update.Email,
		PhoneNumber: Update.PhoneNumber,
	}

	return response, nil
}

func (s *UserService) GetUserLogin(ctx context.Context) (*dto.UserResponse, error) {
	var (
		userLogin = ctx.Value(constants.UserLogin).(*dto.UserResponse)
		data      dto.UserResponse
	)

	data = dto.UserResponse{
		UUID:        userLogin.UUID,
		Name:        userLogin.Name,
		Username:    userLogin.Username,
		Email:       userLogin.Email,
		PhoneNumber: userLogin.PhoneNumber,
		Role:        userLogin.Role,
	}

	return &data, nil
}
func (s *UserService) GetUserByUUID(ctx context.Context, uuid string) (*dto.UserResponse, error) {
	user, err := s.repository.GetUser().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	data := dto.UserResponse{
		UUID:        user.UUID,
		Name:        user.Name,
		Username:    user.Username,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role.Code,
	}

	return &data, nil
}
