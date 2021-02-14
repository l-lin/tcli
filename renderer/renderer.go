package renderer

import "github.com/l-lin/tcli/trello"

type Renderer interface {
	RenderBoards(trello.Boards) string
	RenderLists(trello.Lists) string
	RenderCards(trello.Cards) string
	RenderCard(trello.Card) string
}
