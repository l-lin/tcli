package renderer

import (
	"bytes"
	"fmt"
	"github.com/cheynewallace/tabby"
	"github.com/l-lin/tcli/trello"
	"github.com/rs/zerolog/log"
	"text/tabwriter"
)

type InTable struct {
	lr                          Labels
	cdr                         Description
	minWidth, tabWidth, padding int
	padChar                     byte
	flags                       uint
}

func NewInTableRenderer(lr Labels, cdr Description) Renderer {
	return InTable{
		lr:       lr,
		cdr:      cdr,
		minWidth: 0,
		tabWidth: 0,
		padding:  4,
		padChar:  ' ',
		flags:    0,
	}
}

func (b InTable) RenderBoards(boards trello.Boards) string {
	var buffer bytes.Buffer
	w := tabwriter.NewWriter(&buffer, b.minWidth, b.tabWidth, b.padding, b.padChar, b.flags)
	t := tabby.NewCustom(w)
	t.AddHeader("Name", "Last activity date")
	for _, board := range boards {
		line := make([]interface{}, 2)
		line[0] = board.Name
		line[1] = board.DateLastActivity
		t.AddLine(line...)
	}
	t.Print()
	return buffer.String()
}

func (b InTable) RenderBoard(board trello.Board) string {
	var buffer bytes.Buffer
	w := tabwriter.NewWriter(&buffer, b.minWidth, b.tabWidth, b.padding, b.padChar, b.flags)
	t := tabby.NewCustom(w)
	t.AddLine("ID:", board.ID)
	t.AddLine("Short link:", board.ShortLink)
	t.AddLine("Short URL:", board.ShortURL)
	t.AddLine("Name:", board.Name)
	t.AddLine("Last activity date:", board.DateLastActivity)
	t.Print()
	return buffer.String()
}

func (b InTable) RenderLists(lists trello.Lists) string {
	var buffer bytes.Buffer
	w := tabwriter.NewWriter(&buffer, b.minWidth, b.tabWidth, b.padding, b.padChar, b.flags)
	t := tabby.NewCustom(w)
	t.AddHeader("Name")
	for _, list := range lists {
		t.AddLine(list.Name)
	}
	t.Print()
	return buffer.String()
}

func (b InTable) RenderList(list trello.List) string {
	var buffer bytes.Buffer
	w := tabwriter.NewWriter(&buffer, b.minWidth, b.tabWidth, b.padding, b.padChar, b.flags)
	t := tabby.NewCustom(w)
	t.AddLine("ID:", list.ID)
	t.AddLine("Name:", list.Name)
	t.Print()
	return buffer.String()
}

func (b InTable) RenderCards(cards trello.Cards) string {
	var buffer bytes.Buffer
	w := tabwriter.NewWriter(&buffer, b.minWidth, b.tabWidth, b.padding, b.padChar, b.flags)
	t := tabby.NewCustom(w)
	t.AddHeader("Name", "Position", "Labels")
	for _, card := range cards.SortedByPos() {
		line := make([]interface{}, 3)
		line[0] = card.Name
		line[1] = card.Pos
		line[2] = b.lr.Render(card.Labels)
		t.AddLine(line...)
	}
	t.Print()
	return buffer.String()
}

func (b InTable) RenderCard(card trello.Card) string {
	var buffer bytes.Buffer
	w := tabwriter.NewWriter(&buffer, b.minWidth, b.tabWidth, b.padding, b.padChar, b.flags)
	t := tabby.NewCustom(w)
	t.AddLine("ID:", card.ID)
	t.AddLine("Name:", card.Name)
	t.AddLine("Position:", card.Pos)
	t.AddLine("Short link:", card.ShortLink)
	t.AddLine("Short URL:", card.ShortURL)
	t.AddLine("Labels:", b.lr.Render(card.Labels))
	t.AddLine("Description:", "")
	renderedDescription, err := b.cdr.Render(card.Desc)
	if err != nil {
		log.Debug().
			Err(err).
			Str("idCard", card.ID).
			Msg("could not render card description")
	} else {
		t.AddLine(renderedDescription)
	}
	t.Print()
	return buffer.String()
}

func (b InTable) RenderComments(comments trello.Comments) string {
	var buffer bytes.Buffer
	w := tabwriter.NewWriter(&buffer, b.minWidth, b.tabWidth, b.padding, b.padChar, b.flags)
	t := tabby.NewCustom(w)
	for _, comment := range comments.SortedByDateDesc() {
		t.AddLine()
		t.AddHeader(renderCommentHeader(comment))
		renderedText, err := b.cdr.Render(comment.Data.Text)
		if err != nil {
			log.Debug().
				Err(err).
				Str("idComment", comment.ID).
				Msg("could not render comment text")
		} else {
			t.AddLine(renderedText)
		}
	}
	t.Print()
	return buffer.String()
}

func (b InTable) RenderComment(comment trello.Comment) string {
	var buffer bytes.Buffer
	w := tabwriter.NewWriter(&buffer, b.minWidth, b.tabWidth, b.padding, b.padChar, b.flags)
	t := tabby.NewCustom(w)
	t.AddLine()
	t.AddHeader(renderCommentHeader(comment))
	renderedText, err := b.cdr.Render(comment.Data.Text)
	if err != nil {
		log.Debug().
			Err(err).
			Str("idComment", comment.ID).
			Msg("could not render comment text")
	} else {
		t.AddLine(renderedText)
	}
	t.Print()
	return buffer.String()
}

func renderCommentHeader(comment trello.Comment) string {
	return fmt.Sprintf("%s @ %s [%s]", comment.MemberCreator.Username, comment.Date, comment.ID)
}
