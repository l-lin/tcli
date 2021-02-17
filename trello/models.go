package trello

import (
	"fmt"
	"github.com/logrusorgru/aurora/v3"
)

var mapColors = map[string]func(interface{}) aurora.Value{
	"black":  aurora.BgBrightBlack,
	"blue":   aurora.BgBrightBlue,
	"green":  aurora.BgBrightGreen,
	"lime":   aurora.BgGreen,
	"orange": aurora.BgYellow,
	"pink":   aurora.BgMagenta,
	"purple": aurora.BgBrightMagenta,
	"red":    aurora.BgBrightRed,
	"sky":    aurora.BgBrightCyan,
	"yellow": aurora.BgBrightYellow,
}

func FindBoard(boards Boards, query string) *Board {
	for _, board := range boards {
		if board.TCliID() == query ||
			board.ID == query ||
			board.ShortLink == query ||
			board.Name == query {
			return &board
		}
	}
	return nil
}

type Boards []Board
type Board struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	ShortLink        string `json:"shortLink"`
	ShortURL         string `json:"shortUrl"`
	DateLastActivity string `json:"dateLastActivity"`
}

func (b Board) TCliID() string {
	return toTCliID(b.Name, b.ShortLink)
}

func FindList(lists Lists, query string) *List {
	for _, list := range lists {
		if list.ID == query ||
			list.Name == query {
			return &list
		}
	}
	return nil
}

type Lists []List
type List struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	IDBoard string `json:"idBoard"`
}

func FindCard(cards Cards, query string) *Card {
	for _, card := range cards {
		if card.TCliID() == query ||
			card.ID == query ||
			card.ShortLink == query ||
			card.Name == query {
			return &card
		}
	}
	return nil
}

type Cards []Card
type Card struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"desc"`
	IDBoard     string `json:"idBoard"`
	IDList      string `json:"idList"`
	Closed      bool   `json:"closed"`
	ShortLink   string `json:"shortLink"`
	ShortURL    string `json:"shortUrl"`
	Labels      `json:"labels"`
}

func (c Card) TCliID() string {
	return toTCliID(c.Name, c.ShortLink)
}

type Labels []Label
type Label struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

func (l Label) Colorize(s string) aurora.Value {
	if c := mapColors[l.Color]; c != nil {
		return c(s)
	}
	return aurora.White(s)
}

type UpdateCard struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"desc"`
	IDBoard     string `json:"idBoard"`
	IDList      string `json:"idList"`
	Closed      bool   `json:"closed"`
}

func NewUpdateCard(card Card) UpdateCard {
	return UpdateCard{
		ID:          card.ID,
		Name:        card.Name,
		Description: card.Description,
		IDBoard:     card.IDBoard,
		IDList:      card.IDList,
		Closed:      card.Closed,
	}
}

func CopyUpdateCard(updateCard UpdateCard) UpdateCard {
	return UpdateCard{
		ID:          updateCard.ID,
		Name:        updateCard.Name,
		Description: updateCard.Description,
		IDBoard:     updateCard.IDBoard,
		IDList:      updateCard.IDList,
		Closed:      updateCard.Closed,
	}
}

// toTCliID converts a Trello entity into a unique ID understandable by tcli
// it's using the name, for the user experience in the completion, and the short link
// instead of the id to prevent having long lines in the completion.
func toTCliID(name, shortLink string) string {
	return fmt.Sprintf("%s [%s]", name, shortLink)
}
