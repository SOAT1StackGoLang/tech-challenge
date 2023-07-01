package postgres

import (
	"context"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/domain"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/ports"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const userTable = "lanchonete_users"

type usersRepositoryImpl struct {
	log *zap.SugaredLogger
	db  *gorm.DB
}

func NewPgxUsersRepository(db *gorm.DB, logger *zap.SugaredLogger) ports.UsersRepository {
	return &usersRepositoryImpl{
		log: logger,
		db:  db,
	}
}

func (u usersRepositoryImpl) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	repUser := User{}
	err := u.db.WithContext(ctx).Table(userTable).
		Select("*").Where("id = ?", id).First(repUser).Error
	if err != nil {
		u.log.Errorw(
			"failed getting user",
			zap.String("id", id.String()),
			zap.Error(err),
		)
		return nil, err
	}

	return repUser.toDomain(), nil

}

func (u usersRepositoryImpl) InsertUser(ctx context.Context, user *domain.User) error {
	repUser := User{}
	repUser.fromDomain(user)

	err := u.db.WithContext(ctx).Table(userTable).
		Create(&repUser).Error
	if err != nil {
		u.log.Errorw(
			"failed inserting user",
			zap.String("email", user.Email),
			zap.Error(err),
		)
	}

	return err
}

func (u usersRepositoryImpl) ValidateUser(ctx context.Context, document string) (uuid.UUID, error) {
	var user User
	err := u.db.WithContext(ctx).Table(userTable).
		Select("*").Where("document = ?", document).First(&user).Error
	if err != nil {
		u.log.Errorw(
			"failed validating user",
			zap.String("document", document),
			zap.Error(err),
		)
		return uuid.Nil, err
	}

	return user.ID, nil
}

func (u usersRepositoryImpl) IsUserAdmin(ctx context.Context, id uuid.UUID) (bool, error) {
	var isAdmin bool
	err := u.db.WithContext(ctx).Table(userTable).
		Select("is_admin").Where("id = ?", id).First(isAdmin).Error
	if err != nil {
		u.log.Errorw(
			"failed checking user admin status",
			zap.String("document", id.String()),
			zap.Error(err),
		)
		return false, err
	}

	return isAdmin, nil
}
