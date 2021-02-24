package executor

import (
	"fmt"
	"github.com/l-lin/tcli/trello"
	"github.com/rs/zerolog/log"
)

type cd struct {
	executor
}

func (c cd) Execute(args []string) {
	if len(args) == 0 {
		log.Debug().Msg("returning to top parent")
		c.session.CurrentBoard = nil
		c.session.CurrentList = nil
		return
	}
	if len(args) > 1 {
		fmt.Fprintf(c.stderr, "only one argument is accepted\n")
		return
	}

	arg := args[0]
	if arg == "" {
		log.Debug().Msg("returning to top parent")
		c.session.CurrentBoard = nil
		c.session.CurrentList = nil
		return
	}

	pathResolver := trello.NewPathResolver(c.session)
	p, err := pathResolver.Resolve(arg)
	if err != nil {
		fmt.Fprintf(c.stderr, "%v\n", err)
		return
	}

	if p.BoardName == "" {
		c.session.CurrentBoard = nil
		c.session.CurrentList = nil
		return
	}

	var board *trello.Board
	if board, err = c.tr.FindBoard(p.BoardName); err != nil || board == nil {
		fmt.Fprintf(c.stderr, "no board found with name '%s'\n", p.BoardName)
		return
	}

	if p.ListName == "" {
		c.session.CurrentBoard = board
		c.session.CurrentList = nil
		return
	}

	var list *trello.List
	if list, err = c.tr.FindList(board.ID, p.ListName); err != nil || list == nil {
		fmt.Fprintf(c.stderr, "no list found with name '%s'\n", p.ListName)
		return
	}

	if p.CardName != "" {
		fmt.Fprintf(c.stderr, "cannot cd on card\n")
		return
	}

	c.session.CurrentBoard = board
	c.session.CurrentList = list
}
