//go:generate mockgen -source renderer.go -package renderer -destination renderer_mock.go
package renderer

import (
	"github.com/l-lin/tcli/trello"
)

// Renderer for displaying Trello entities in human friendly way
type Renderer interface {
	RenderBoards(trello.Boards) string
	RenderBoard(trello.Board) string
	RenderLists(trello.Lists) string
	RenderList(trello.List) string
	RenderCards(trello.Cards) string
	RenderCard(trello.Card) string
}
