package usecases

import (
	"context"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/domain"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/ports"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type usersUseCase struct {
	logger   *zap.SugaredLogger
	userRepo ports.UsersRepository
}

func NewUsersUseCase(userRepo ports.UsersRepository, log *zap.SugaredLogger) ports.UsersUseCase {
	return &usersUseCase{userRepo: userRepo, logger: log}
}

func (u usersUseCase) IsUserAdmin(ctx context.Context, id uuid.UUID) (bool, error) {
	isAdmin, err := u.userRepo.IsUserAdmin(ctx, id)
	if err != nil {
		return false, err
	}

	return isAdmin, err
}

func (u usersUseCase) CreateUser(ctx context.Context, name, document, email string) (*domain.User, error) {
	user := domain.NewUser(uuid.New(), document, name, email)

	err := u.userRepo.InsertUser(ctx, user)
	if err != nil {
		u.logger.Errorw(
			"failed inserting user",
			zap.String("document", document),
			zap.Error(err),
		)
		return nil, err
	}

	return user, err
}

func (u usersUseCase) ValidateUser(ctx context.Context, document string) (uuid.UUID, error) {
	uID, err := u.userRepo.ValidateUser(ctx, document)
	if err != nil {
		u.logger.Errorw(
			"failed validating user",
			zap.String("document", document),
			zap.Error(err),
		)
		return [16]byte{}, err
	}
	return uID, err
}
