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
		then().
		doOnEmptyBoardName(func() {
			fmt.Fprintf(r.stderr, "nothing to archive\n")
		}).
		findBoard().
		doOnBoard(func(board *trello.Board) {
			fmt.Fprintf(r.stderr, "board archiving not implemented yet\n")
		}).
		then().
		findList().
		doOnList(func(list *trello.List) {
			fmt.Fprintf(r.stderr, "list archiving not implemented yet\n")
		}).
		then().
		findCard().
		doOnCard(func(card *trello.Card) {
			prompt := promptui.Prompt{
				Label:     fmt.Sprintf("Archive card '%s'?", card.Name),
				IsConfirm: true,
				Stdin:     r.stdin,
			}
			if _, err := prompt.Run(); err != nil {
				return
			}
			updatedCard := trello.NewUpdateCard(*card)
			updatedCard.Closed = true
			if _, err := r.tr.UpdateCard(updatedCard); err != nil {
				fmt.Fprintf(r.stderr, "could not archive card '%s': %s\n", card.Name, err)
			}
		}).
		then().
		findComment().
		doOnComment(func(comment *trello.Comment) {
			prompt := promptui.Prompt{
				Label:     fmt.Sprintf("Delete comment '%s'?", comment.ID),
				IsConfirm: true,
				Stdin:     r.stdin,
			}
			if _, err := prompt.Run(); err != nil {
				return
			}
			if err := r.tr.DeleteComment(comment.Data.Card.ID, comment.ID); err != nil {
				fmt.Fprintf(r.stderr, "could not delete comment '%s': %s\n", comment.ID, err)
			}
		})
	if exec.err != nil {
		fmt.Fprintf(r.stderr, "%s\n", exec.err)
	}
}
