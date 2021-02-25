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

	exec := start(t.tr).
		resolvePath(t.session, arg).
		doOnEmptyBoardName(func() {
			fmt.Fprintf(t.stderr, "nothing to create\n")
		}).
		thenFindBoard().
		doOnEmptyListName(func(_ *trello.Session) {
			fmt.Fprintf(t.stderr, "board creation not implemented yet\n")
		}).
		thenFindList().
		doOnEmptyCardName(func(_ *trello.Session) {
			fmt.Fprintf(t.stderr, "list creation not implemented yet\n")
		}).
		doOnCardName(func(cardName string, session *trello.Session) {
			createCard := trello.CreateCard{
				Name:   cardName,
				IDList: session.List.ID,
			}
			if _, err := t.tr.CreateCard(createCard); err != nil {
				fmt.Fprintf(t.stderr, "could not create card '%s': %v\n", cardName, err)
			}
		})
	if exec.err != nil {
		fmt.Fprintf(t.stderr, "%s\n", exec.err)
	} else if !exec.isFinished {
		fmt.Fprintf(t.stderr, "comment creation not implemented yet\n")
	}
}
