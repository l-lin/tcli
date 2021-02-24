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
	pathResolver := trello.NewPathResolver(r.session)
	p, err := pathResolver.Resolve(arg)
	if err != nil {
		fmt.Fprintf(r.stderr, "%v\n", err)
		return
	}

	if p.BoardName == "" {
		fmt.Fprintf(r.stderr, "nothing to archive\n")
		return
	}

	var board *trello.Board
	if board, err = r.tr.FindBoard(p.BoardName); err != nil || board == nil {
		fmt.Fprintf(r.stderr, "no board found with name '%s'\n", p.BoardName)
		return
	}

	if p.ListName == "" {
		fmt.Fprintf(r.stderr, "board archiving not implemented yet\n")
		return
	}

	var list *trello.List
	if list, err = r.tr.FindList(board.ID, p.ListName); err != nil || list == nil {
		fmt.Fprintf(r.stderr, "no list found with name '%s'\n", p.ListName)
		return
	}

	if p.CardName == "" {
		fmt.Fprintf(r.stderr, "list archiving not implemented yet\n")
		return
	}

	var card *trello.Card
	if card, err = r.tr.FindCard(list.ID, p.CardName); err != nil || card == nil {
		fmt.Fprintf(r.stderr, "no card found with name '%s'\n", p.CardName)
		return
	}
	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("Archive card '%s'?", card.Name),
		IsConfirm: true,
		Stdin:     r.stdin,
	}
	if _, err = prompt.Run(); err != nil {
		return
	}
	updatedCard := trello.NewUpdateCard(*card)
	updatedCard.Closed = true
	if _, err = r.tr.UpdateCard(updatedCard); err != nil {
		fmt.Fprintf(r.stderr, "could not archive card '%s': %s\n", p.CardName, err)
	}
}
