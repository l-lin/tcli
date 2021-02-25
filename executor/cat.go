package executor

import (
	"fmt"
	"github.com/l-lin/tcli/trello"
)

type cat struct {
	executor
}

func (c cat) Execute(args []string) {
	if len(args) != 0 {
		for _, arg := range args {
			c.execute(arg)
		}
	}
}

func (c cat) execute(arg string) {
	if arg == "" {
		return
	}

	if err := start(c.tr).
		resolvePath(c.session, arg).
		then().
		findBoard().
		doOnBoard(func(board *trello.Board) {
			fmt.Fprintf(c.stdout, "%s\n", c.r.RenderBoard(*board))
		}).
		then().
		findList().
		doOnList(func(list *trello.List) {
			fmt.Fprintf(c.stdout, "%s\n", c.r.RenderList(*list))
		}).
		then().
		findCard().
		doOnCard(func(card *trello.Card) {
			fmt.Fprintf(c.stdout, "%s\n", c.r.RenderCard(*card))
		}).
		then().
		findComment().
		doOnComment(func(comment *trello.Comment) {
			fmt.Fprintf(c.stdout, "%s\n", c.r.RenderComment(*comment))
		}); err != nil {
		fmt.Fprintf(c.stderr, "%s\n", err)
	}
}
