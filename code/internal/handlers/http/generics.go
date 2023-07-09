package http

import (
	"github.com/google/uuid"
)

type (
	QueryStruct struct {
		ID     string `json:"id" description:"id do objeto a ser removido/consultado"`
		UserID string `json:"user_id" description:"id do usuário requerente"`
	}
)

func (q *QueryStruct) parseToUuid() (uuid.UUID, uuid.UUID, error) {
	id, err := uuid.Parse(q.ID)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	uid, err := uuid.Parse(q.UserID)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	return id, uid, err
}
