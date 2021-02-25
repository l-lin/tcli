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
		then().
		doOnEmptyBoardName(func() {
			fmt.Fprintf(t.stderr, "nothing to create\n")
		}).
		findBoard().
		doOnBoard(func(board *trello.Board) {
			fmt.Fprintf(t.stderr, "board creation not implemented yet\n")
		}).
		then().
		findList().
		doOnList(func(list *trello.List) {
			fmt.Fprintf(t.stderr, "list creation not implemented yet\n")
		}).
		then().
		doOnCardName(func(cardName string, session *trello.Session) {
			createCard := trello.CreateCard{
				Name:   cardName,
				IDList: session.List.ID,
			}
			if _, err := t.tr.CreateCard(createCard); err != nil {
				fmt.Fprintf(t.stderr, "could not create card '%s': %v\n", cardName, err)
			}
		}).
		findCard().
		then().
		doOnCommentText(func(commentText string, session *trello.Session) {
			createComment := trello.CreateComment{
				IDCard: session.Card.ID,
				Text:   commentText,
			}
			if _, err := t.tr.CreateComment(createComment); err != nil {
				fmt.Fprintf(t.stderr, "could not create comment '%s': %v\n", commentText, err)
			}
		})
	if exec.err != nil {
		fmt.Fprintf(t.stderr, "%s\n", exec.err)
	}
}
