package db_stub

import (
	"crud/api/users"
	"crud/storage"
	"github.com/golang/protobuf/proto"
)

func NewUsersStorage() (storage.Users, error) {
	return &usersStorage{
		db: make(map[uint64]*users.User),
	}, nil
}

type usersStorage struct {
	db            map[uint64]*users.User
	autoIncrement uint64
}

func (u *usersStorage) Create(user *users.User) (*users.User, error) {
	if user == nil {
		return nil, storage.InvalidArgumentError
	}

	if user.Id > 0 {
		return nil, storage.AlreadyExistsError
	}

	u.autoIncrement++
	user.Id = u.autoIncrement
	u.db[user.Id] = proto.Clone(user).(*users.User)

	return user, nil
}

func (u *usersStorage) Get(id uint64) (*users.User, error) {
	if user, ok := u.db[id]; ok {
		return proto.Clone(user).(*users.User), nil
	}

	return nil, storage.NotFoundError
}

func (u *usersStorage) Update(user *users.User) (*users.User, error) {
	if user == nil {
		return nil, storage.InvalidArgumentError
	}

	if _, ok := u.db[user.Id]; ok {
		u.db[user.Id] = user

		return user, nil
	}

	return nil, storage.NotFoundError
}

func (u *usersStorage) Delete(id uint64) error {
	if _, ok := u.db[id]; ok {
		delete(u.db, id)

		return nil
	}

	return storage.NotFoundError
}
