package executor

import (
	"fmt"
	"github.com/l-lin/tcli/trello"
	"github.com/rs/zerolog/log"
)

type cd struct {
	executor
}

func (c *cd) Execute(args []string) {
	if len(args) == 0 {
		log.Debug().Msg("returning to top parent")
		c.registerSession(&trello.Session{})
		return
	}
	if len(args) > 1 {
		fmt.Fprintf(c.stderr, "only one argument is accepted\n")
		return
	}

	arg := args[0]
	if arg == "" {
		log.Debug().Msg("returning to top parent")
		c.registerSession(&trello.Session{})
		return
	}

	exec := start(c.tr).
		resolvePath(c.session, arg).
		doOnEmptyBoardName(func() {
			c.registerSession(&trello.Session{})
		}).
		thenFindBoard().
		doOnEmptyListName(c.registerSession).
		thenFindList().
		doOnEmptyCardName(c.registerSession).
		thenFindCard().
		doOnEmptyCommentID(c.registerSession)
	if exec.err != nil {
		fmt.Fprintf(c.stderr, "%s\n", exec.err)
	} else if !exec.isFinished {
		fmt.Fprintf(c.stderr, "cannot cd on comment\n")
	}
}

func (c *cd) registerSession(session *trello.Session) {
	c.session.Board = session.Board
	c.session.List = session.List
	c.session.Card = session.Card
}
