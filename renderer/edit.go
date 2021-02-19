package renderer

import "github.com/l-lin/tcli/trello"

type Edit interface {
	MarshalCardToCreate(trello.CardToCreate, trello.Lists, trello.Labels) ([]byte, error)
	MarshalCardToEdit(trello.CardToEdit, trello.Lists, trello.Labels) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}
