package orders

import "errors"

var ErrNotFound = errors.New("order not foundx")

type Repository interface {
	Save(*Order) error
	ById(ID) (*Order, error)
}
