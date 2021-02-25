package executor

import (
	"fmt"
	"github.com/l-lin/tcli/renderer"
	"github.com/l-lin/tcli/trello"
	"github.com/manifoldco/promptui"
	"github.com/rs/zerolog/log"
	"io"
)

type edit struct {
	executor
	editor       Editor
	stdin        io.ReadCloser
	editRenderer renderer.Edit
}

func (e edit) Execute(args []string) {
	if len(args) != 0 {
		for _, arg := range args {
			e.execute(arg)
		}
	}
}

func (e edit) execute(arg string) {
	exec := start(e.tr).
		resolvePath(e.session, arg).
		doOnEmptyBoardName(func() {
			fmt.Fprintf(e.stderr, "nothing to edit\n")
		}).
		thenFindBoard().
		doOnEmptyListName(func(session *trello.Session) {
			fmt.Fprintf(e.stderr, "board edition not implemented yet\n")
		}).
		thenFindList().
		doOnEmptyCardName(func(session *trello.Session) {
			fmt.Fprintf(e.stderr, "list edition not implemented yet\n")
		}).
		doOnCardName(func(cardName string, session *trello.Session) {
			var card *trello.Card
			var err error
			if card, err = e.tr.FindCard(session.List.ID, cardName); err != nil || card == nil {
				log.Debug().
					Str("cardName", cardName).
					Msg("no card found => creating new card")
				card = &trello.Card{
					Name:    cardName,
					Desc:    "",
					IDBoard: session.Board.ID,
					IDList:  session.List.ID,
				}
				if err = e.createCard(*card); err != nil {
					fmt.Fprintf(e.stderr, "could not create card '%s': %v\n", cardName, err)
				}
			} else {
				if err = e.editCard(*card); err != nil {
					fmt.Fprintf(e.stderr, "could not edit card '%s': %v\n", cardName, err)
				}
			}
		})
	if exec.err != nil {
		fmt.Fprintf(e.stderr, "%s\n", exec.err)
		return
	}
}

func (e edit) createCard(card trello.Card) (err error) {
	var lists trello.Lists
	if lists, err = e.tr.FindLists(card.IDBoard); err != nil {
		return
	}

	var labels trello.Labels
	if labels, err = e.tr.FindLabels(card.IDBoard); err != nil {
		return
	}

	ctc := trello.NewCardToCreate(card)
	var in []byte
	if in, err = e.editRenderer.MarshalCardToCreate(ctc, lists, labels); err != nil {
		return
	}

	var out []byte
	if out, err = e.editor.Edit(in); err != nil {
		return
	}

	var editedCard trello.CardToCreate
	if err = e.editRenderer.Unmarshal(out, &editedCard); err != nil {
		return
	}
	createdCard := trello.NewCreateCard(card)
	createdCard.Name = editedCard.Name
	createdCard.Desc = editedCard.Desc
	createdCard.IDList = editedCard.IDList
	createdCard.Pos = editedCard.GetPos()
	createdCard.IDLabels = labels.FilterBy(
		editedCard.Labels,
		trello.LabelFilterOr(trello.LabelFilterByID, trello.LabelFilterByTCliColor, trello.LabelFilterByColor),
	).IDLabelsInString()

	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("Do you want to create the card '%s'?", createdCard.Name),
		IsConfirm: true,
		Stdin:     e.stdin,
	}
	if _, err = prompt.Run(); err != nil {
		fmt.Fprintf(e.stdout, "card '%s' not created\n", card.Name)
		return nil
	}
	_, err = e.tr.CreateCard(createdCard)
	return
}

func (e edit) editCard(card trello.Card) (err error) {
	var lists trello.Lists
	if lists, err = e.tr.FindLists(card.IDBoard); err != nil {
		return
	}
	var labels trello.Labels
	if labels, err = e.tr.FindLabels(card.IDBoard); err != nil {
		return
	}

	cte := trello.NewCardToEdit(card)
	var in []byte
	if in, err = e.editRenderer.MarshalCardToEdit(cte, lists, labels); err != nil {
		return
	}

	var out []byte
	if out, err = e.editor.Edit(in); err != nil {
		return
	}

	var editedCard trello.CardToEdit
	if err = e.editRenderer.Unmarshal(out, &editedCard); err != nil {
		return
	}
	updatedCard := trello.NewUpdateCard(card)
	updatedCard.Name = editedCard.Name
	updatedCard.Desc = editedCard.Desc
	updatedCard.Closed = editedCard.Closed
	updatedCard.IDList = editedCard.IDList
	updatedCard.Pos = editedCard.GetPos()
	updatedCard.IDLabels = labels.FilterBy(
		editedCard.Labels,
		trello.LabelFilterOr(trello.LabelFilterByID, trello.LabelFilterByTCliColor, trello.LabelFilterByColor),
	).IDLabelsInString()

	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("Do you want to update the card '%s'?", updatedCard.Name),
		IsConfirm: true,
		Stdin:     e.stdin,
	}
	if _, err = prompt.Run(); err != nil {
		return nil
	}
	_, err = e.tr.UpdateCard(updatedCard)
	return
}
