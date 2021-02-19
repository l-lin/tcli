package renderer

import (
	"bytes"
	"github.com/cheynewallace/tabby"
	"github.com/l-lin/tcli/trello"
	"github.com/rs/zerolog/log"
	"sort"
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
	t.AddHeader("Name", "ID", "Short URL", "Last activity date")
	for _, board := range boards {
		line := make([]interface{}, 4)
		line[0] = board.Name
		line[1] = board.ID
		line[2] = board.ShortURL
		line[3] = board.DateLastActivity
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
	t.AddHeader("Name", "ID")
	for _, list := range lists {
		line := make([]interface{}, 2)
		line[0] = list.Name
		line[1] = list.ID
		t.AddLine(line...)
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
	t.AddHeader("Name", "ID", "Short URL", "Labels")
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Pos < cards[j].Pos
	})
	for _, card := range cards {
		line := make([]interface{}, 4)
		line[0] = card.Name
		line[1] = card.ID
		line[2] = card.ShortURL
		line[3] = b.lr.Render(card.Labels)
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
	renderedDescription, err := b.cdr.Render(card.Description)
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
