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
		thenFindBoard().
		doOnEmptyListName(func(session *trello.Session) {
			fmt.Fprintf(c.stdout, "%s\n", c.r.RenderBoard(*session.Board))
		}).
		thenFindList().
		doOnEmptyCardName(func(session *trello.Session) {
			fmt.Fprintf(c.stdout, "%s\n", c.r.RenderList(*session.List))
		}).
		thenFindCard().
		doOnEmptyCommentID(func(session *trello.Session) {
			fmt.Fprintf(c.stdout, "%s\n", c.r.RenderCard(*session.Card))
		}).
		thenFindComment().
		andDoOnComment(func(comment *trello.Comment) {
			fmt.Fprintf(c.stdout, "%s\n", c.r.RenderComment(*comment))
		}); err != nil {
		fmt.Fprintf(c.stderr, "%s\n", err)
	}
}
