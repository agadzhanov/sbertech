package main

import (
	"crud/api/users"
	"crud/config"
	"crud/db_stub"
	"crud/logger"
	"crud/server"
	"crud/storage"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// получаем логгер
	log, err := logger.NewLogger()
	if err != nil {
		panic(fmt.Sprintf("failed getting logger, %s", err))
	}

	// настраиваем подключение
	c := config.GetConfig()
	var l net.Listener
	l, err = net.Listen(c.Network, c.Address())
	if err != nil {
		panic(fmt.Sprintf("failed getting listener with %s, %s: %s", c.Network, c.Domain, err))
	}

	defer func() {

		if err := l.Close(); err != nil {
			log.Errorf("Failed to close %s %s: %v", c.Network, c.Domain, err)
		}
	}()

	// инициализируем сервер сервиса users
	var usersStorage storage.Users
	usersStorage, err = db_stub.NewUsersStorage()
	usersServerOptions := []interface{}{
		log,
		usersStorage,
	}
	usersServer := server.NewUsersServer(usersServerOptions)

	s := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_logrus.UnaryServerInterceptor(log),
			grpc_recovery.UnaryServerInterceptor(),
		),
	)
	users.RegisterUsersServer(s, usersServer)

	term := make(chan os.Signal)
	go func() {
		if err := s.Serve(l); err != nil {
			term <- syscall.SIGINT
		}
	}()

	log.Info("server started")

	signal.Notify(term, syscall.SIGTERM, syscall.SIGINT)
	<-term
	s.GracefulStop()
	log.Info("server stopped")
}
