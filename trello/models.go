package trello

import "github.com/logrusorgru/aurora/v3"

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
