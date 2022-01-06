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
	editor        Editor
	stdin         io.ReadCloser
	editRenderer  renderer.Edit
	defaultLabels []string
	neverPrompt   bool
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
		then().
		doOnEmptyBoardName(func() {
			fmt.Fprintf(e.stderr, "nothing to edit\n")
		}).
		findBoard().
		doOnBoard(func(board *trello.Board) {
			fmt.Fprintf(e.stderr, "board edition not implemented yet\n")
		}).
		then().
		findList().
		doOnList(func(list *trello.List) {
			fmt.Fprintf(e.stderr, "list edition not implemented yet\n")
		}).
		then().
		findCard().
		doOnCard(func(card *trello.Card) {
			if err := e.editCard(*card); err != nil {
				fmt.Fprintf(e.stderr, "could not edit card '%s': %v\n", card.Name, err)
			}
		}).
		doOnCardName(func(cardName string, session *trello.Session) {
			log.Debug().
				Str("cardName", cardName).
				Msg("no card found => creating new card")
			card := &trello.Card{
				Name:    cardName,
				Desc:    "",
				IDBoard: session.Board.ID,
				IDList:  session.List.ID,
			}
			if err := e.createCard(*card); err != nil {
				fmt.Fprintf(e.stderr, "could not create card '%s': %v\n", cardName, err)
			}
		}).
		then().
		findComment().
		doOnComment(func(comment *trello.Comment) {
			if err := e.updateComment(comment); err != nil {
				fmt.Fprintf(e.stderr, "could not update comment '%s': %v\n", comment.ID, err)
			}
		}).
		doOnCommentText(func(commentText string, session *trello.Session) {
			if err := e.createComment(commentText, session.Card.ID); err != nil {
				fmt.Fprintf(e.stderr, "could not create comment '%s': %v\n", commentText, err)
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

	ctc := trello.NewCardToCreate(card, e.defaultLabels)
	var in []byte
	if in, err = e.editRenderer.MarshalCardToCreate(ctc, lists, labels); err != nil {
		return
	}

	var out []byte
	if out, err = e.editor.Edit(in, e.editRenderer.GetFileType()); err != nil {
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

	if !e.neverPrompt {
		prompt := promptui.Prompt{
			Label:     fmt.Sprintf("Do you want to create the card '%s'", createdCard.Name),
			IsConfirm: true,
			Stdin:     e.stdin,
		}
		if _, err = prompt.Run(); err != nil {
			fmt.Fprintf(e.stdout, "card '%s' not created\n", card.Name)
			return nil
		}
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
	if out, err = e.editor.Edit(in, e.editRenderer.GetFileType()); err != nil {
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

	if !e.neverPrompt {
		prompt := promptui.Prompt{
			Label:     fmt.Sprintf("Do you want to update the card '%s'", updatedCard.Name),
			IsConfirm: true,
			Stdin:     e.stdin,
		}
		if _, err = prompt.Run(); err != nil {
			return nil
		}
	}
	_, err = e.tr.UpdateCard(updatedCard)
	return
}

func (e edit) createComment(commentText, idCard string) (err error) {
	in := []byte(commentText)

	var out []byte
	if out, err = e.editor.Edit(in, markdownFileType); err != nil {
		return
	}

	editedText := string(out)

	if !e.neverPrompt {
		prompt := promptui.Prompt{
			Label:     "Do you want to create the comment?",
			IsConfirm: true,
			Stdin:     e.stdin,
		}
		if _, err = prompt.Run(); err != nil {
			return nil
		}
	}

	createComment := trello.CreateComment{
		IDCard: idCard,
		Text:   editedText,
	}
	_, err = e.tr.CreateComment(createComment)
	return
}

func (e edit) updateComment(comment *trello.Comment) (err error) {
	in := []byte(comment.Data.Text)

	var out []byte
	if out, err = e.editor.Edit(in, markdownFileType); err != nil {
		return
	}

	editedText := string(out)

	if !e.neverPrompt {
		prompt := promptui.Prompt{
			Label:     fmt.Sprintf("Do you want to update the comment '%s'?", comment.ID),
			IsConfirm: true,
			Stdin:     e.stdin,
		}
		if _, err = prompt.Run(); err != nil {
			return nil
		}
	}

	updateComment := trello.UpdateComment{
		ID:     comment.ID,
		IDCard: comment.Data.Card.ID,
		Text:   editedText,
	}
	_, err = e.tr.UpdateComment(updateComment)
	return
}
