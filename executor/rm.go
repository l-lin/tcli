package executor

import (
	"fmt"
	"github.com/l-lin/tcli/trello"
	"github.com/manifoldco/promptui"
	"io"
)

type rm struct {
	executor
	stdin io.ReadCloser
}

func (r rm) Execute(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(r.stderr, "missing card operand\n")
		return
	}
	for _, arg := range args {
		r.execute(arg)
	}
	return
}

func (r rm) execute(arg string) {
	if arg == "" {
		fmt.Fprintf(r.stderr, "missing card operand\n")
		return
	}
	exec := start(r.tr).
		resolvePath(r.session, arg).
		doOnEmptyBoardName(func() {
			fmt.Fprintf(r.stderr, "nothing to archive\n")
		}).
		thenFindBoard().
		doOnEmptyListName(func(_ *trello.Session) {
			fmt.Fprintf(r.stderr, "board archiving not implemented yet\n")
		}).
		thenFindList().
		doOnEmptyCardName(func(_ *trello.Session) {
			fmt.Fprintf(r.stderr, "list archiving not implemented yet\n")
		}).
		thenFindCard().
		doOnEmptyCommentID(func(session *trello.Session) {
			prompt := promptui.Prompt{
				Label:     fmt.Sprintf("Archive card '%s'?", session.Card.Name),
				IsConfirm: true,
				Stdin:     r.stdin,
			}
			if _, err := prompt.Run(); err != nil {
				return
			}
			updatedCard := trello.NewUpdateCard(*session.Card)
			updatedCard.Closed = true
			if _, err := r.tr.UpdateCard(updatedCard); err != nil {
				fmt.Fprintf(r.stderr, "could not archive card '%s': %s\n", session.Card.Name, err)
			}
		})
	if exec.err != nil {
		fmt.Fprintf(r.stderr, "%s\n", exec.err)
	}
}
