package domain

import "github.com/google/uuid"

func (o *Order) ExtractProductsIDs() []uuid.UUID {
	if len(o.Products) == 0 {
		return []uuid.UUID{}
	}
	out := make([]uuid.UUID, 0, len(o.Products))

	for _, v := range o.Products {
		out = append(out, v.ID)
	}

	return out
}
