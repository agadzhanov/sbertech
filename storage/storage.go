package storage

import (
	"crud/api/users"
	"errors"
)

type Users interface {
	Create(user *users.User) (*users.User, error)
	Get(id uint64) (*users.User, error)
	Update(user *users.User) (*users.User, error)
	Delete(id uint64) error
}

var (
	InvalidArgumentError = errors.New("invalid argument")
	AlreadyExistsError   = errors.New("already exists")
	NotFoundError        = errors.New("not found")
)
