package renderer

import "github.com/l-lin/tcli/trello"

type Edit interface {
	Marshal(trello.CardToEdit, trello.Lists) ([]byte, error)
	Unmarshal([]byte, *trello.CardToEdit) error
}
