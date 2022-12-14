package finding

import (
	"context"
	"database/sql"
	"errors"

	users "github.com/rcebrian/users-service/internal"
)

//go:generate mockery --case=snake --outpkg=mocks --output=../mocks --name=FindAllUsersUseCase

type FindAllUsersUseCase interface {
	FindAll(ctx context.Context) ([]users.User, error)
}

type allUseCase struct {
	repository users.UserRepository
}

func NewFindAllUsersUseCase(repository users.UserRepository) FindAllUsersUseCase {
	return allUseCase{repository: repository}
}

// FindAll get all users from database
func (s allUseCase) FindAll(ctx context.Context) ([]users.User, error) {
	all, err := s.repository.FindAll(ctx)
	if err == nil {
		return all, nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, users.ErrNotFound
	}

	return nil, err
}
