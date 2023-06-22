package users

import (
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/domain"
	"github.com/google/uuid"
)

type User struct {
	ID       string `json:"id"`
	Document string `json:"document"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

func (u *User) fromDomain(user *domain.User) {
	if u == nil {
		u = &User{}
	}

	u.ID = user.ID.String()
	u.Name = user.Name
	u.Document = user.Document
	u.Email = user.Email
}

func (u *User) toDomain() *domain.User {
	if u == nil {
		u = User{}
	}

	id, err := uuid.FromBytes([]byte(u.ID))
	if err != nil {
		return nil
	}

	return &domain.User{
		ID:       id,
		Document: u.Document,
		Name:     u.Name,
		Email:    u.Email,
	}
}
