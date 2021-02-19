package renderer

import "github.com/l-lin/tcli/trello"

type Edit interface {
	MarshalCardToCreate(trello.CardToCreate, trello.Lists) ([]byte, error)
	MarshalCardToEdit(trello.CardToEdit, trello.Lists) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}
