package http

type (
	DeletionStruct struct {
		ID     string `json:"id" description:"id do objeto a ser removido"`
		UserID string `json:"user_id" description:"id do usuário requerente"`
	}
)
