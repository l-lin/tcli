package trello

import "github.com/logrusorgru/aurora/v3"

var mapColors = map[string]func(interface{}) aurora.Value{
	"black":  aurora.BrightBlack,
	"blue":   aurora.BrightBlue,
	"green":  aurora.BrightGreen,
	"lime":   aurora.Green,
	"orange": aurora.Yellow,
	"pink":   aurora.Magenta,
	"purple": aurora.BrightMagenta,
	"red":    aurora.BrightRed,
	"sky":    aurora.BrightCyan,
	"yellow": aurora.BrightYellow,
}

type Boards []Board
type Board struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	ShortURL         string `json:"shortUrl"`
	DateLastActivity string `json:"dateLastActivity"`
}

type Lists []List
type List struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	IDBoard string `json:"idBoard"`
}

type Cards []Card
type Card struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"desc"`
	IDBoard     string `json:"idBoard"`
	IDList      string `json:"idList"`
	Labels      `json:"labels"`
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
