package executor

import (
	"fmt"
	"github.com/l-lin/tcli/trello"
)

type touch struct {
	executor
}

func (t touch) Execute(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(t.stderr, "missing card operand\n")
	} else {
		for _, arg := range args {
			t.execute(arg)
		}
	}
}

func (t touch) execute(arg string) {
	if arg == "" {
		fmt.Fprintf(t.stderr, "missing card operand\n")
		return
	}
	pathResolver := trello.NewPathResolver(t.session)
	p, err := pathResolver.Resolve(arg)
	if err != nil {
		fmt.Fprintf(t.stderr, "%v\n", err)
		return
	}

	if p.BoardName == "" {
		fmt.Fprintf(t.stderr, "nothing to create\n")
		return
	}

	var board *trello.Board
	if board, err = t.tr.FindBoard(p.BoardName); err != nil || board == nil {
		fmt.Fprintf(t.stderr, "no board found with name '%s'\n", p.BoardName)
		return
	}

	if p.ListName == "" {
		fmt.Fprintf(t.stderr, "board creation not implemented yet\n")
		return
	}

	var list *trello.List
	if list, err = t.tr.FindList(board.ID, p.ListName); err != nil || list == nil {
		fmt.Fprintf(t.stderr, "no list found with name '%s'\n", p.ListName)
		return
	}

	if p.CardName == "" {
		fmt.Fprintf(t.stderr, "list creation not implemented yet\n")
		return
	}

	createCard := trello.CreateCard{
		Name:   p.CardName,
		IDList: list.ID,
	}
	if _, err = t.tr.CreateCard(createCard); err != nil {
		fmt.Fprintf(t.stderr, "could not create card '%s': %v\n", p.CardName, err)
	}
}
