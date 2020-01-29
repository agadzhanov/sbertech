package server

import (
	"context"
	"crud/api/users"
	"crud/storage"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	commonErrorTest = "request failed"
)

type srv struct {
	logger  *logrus.Entry
	storage storage.Users
}

func (s *srv) with(option interface{}) {
	switch t := option.(type) {
	case *logrus.Entry:
		s.logger = t
	case storage.Users:
		s.storage = t
	default:
		panic(fmt.Sprintf("unexpexted option type %s", t))
	}

}

func NewUsersServer(options []interface{}) users.UsersServer {
	server := &srv{}
	for index := range options {
		server.with(options[index])
	}
	if server.logger == nil {
		panic("server.NewUsersServer() logger is required")
	}
	if server.storage == nil {
		panic("server.NewUsersServer() storage is required")
	}
	return server
}

func (s *srv) CreateUser(ctx context.Context, request *users.CreateUserRequest) (*users.CreateUserResponse, error) {
	user, err := s.storage.Create(request.GetUser())
	if err == nil {
		return &users.CreateUserResponse{User: user}, nil
	}

	code := s.storageErrorToStatusCode(err)

	msg := fmt.Sprintf("storage.Create() failed: %s", err)

	s.logger.Error(msg)

	return nil, status.Error(code, commonErrorTest)
}

func (s srv) GetUser(ctx context.Context, request *users.GetUserRequest) (*users.GetUserResponse, error) {
	user, err := s.storage.Get(request.GetId())
	if err == nil {
		return &users.GetUserResponse{
			User: user,
		}, nil
	}

	msg := fmt.Sprintf("storage.Get() failed: %s", err)
	s.logger.Error(msg)

	statusCode := s.storageErrorToStatusCode(err)
	return nil, status.Error(statusCode, commonErrorTest)
}

func (s srv) UpdateUser(ctx context.Context, request *users.UpdateUserRequest) (*users.UpdateUserResponse, error) {
	user, err := s.storage.Update(request.GetUser())
	if err == nil {
		return &users.UpdateUserResponse{User: user}, nil
	}

	code := s.storageErrorToStatusCode(err)
	msg := fmt.Sprintf("storage.update() failed: %s", err)

	s.logger.Error(msg)

	return nil, status.Error(code, commonErrorTest)
}

func (s srv) DeleteUser(ctx context.Context, request *users.DeleteUserRequest) (*users.DeleteUserResponse, error) {
	_, err := s.storage.Get(request.GetId())
	if err == nil {
		err = s.storage.Delete(request.GetId())
		if err == nil {
			return &users.DeleteUserResponse{}, nil
		}
	}

	msg := fmt.Sprintf("storage.Delete() failed: %s", err)
	s.logger.Error(msg)
	statusCode := s.storageErrorToStatusCode(err)
	return nil, status.Error(statusCode, commonErrorTest)
}

// Задает соответствие ошибок storage и кодов ответа grpc
func (s *srv) storageErrorToStatusCode(err error) codes.Code {
	var code codes.Code
	switch err {
	case storage.NotFoundError:
		code = codes.NotFound
	case storage.InvalidArgumentError:
		code = codes.InvalidArgument
	case storage.AlreadyExistsError:
		code = codes.AlreadyExists
	default:
		code = codes.Internal
	}

	return code
}
