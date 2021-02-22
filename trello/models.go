package trello

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
)

func FindBoard(boards Boards, query string) *Board {
	sanitizedQuery := sanitize(query)
	for _, board := range boards {
		if board.TCliID() == sanitizedQuery ||
			board.ID == query ||
			board.ShortLink == query ||
			board.Name == query ||
			board.SanitizedName() == sanitizedQuery {
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
	return toTCliID(b.SanitizedName(), b.ShortLink)
}

func (b Board) SanitizedName() string {
	return sanitize(b.Name)
}

func FindList(lists Lists, query string) *List {
	sanitizedQuery := sanitize(query)
	for _, list := range lists {
		if list.TCliID() == sanitizedQuery ||
			list.ID == query ||
			list.Name == query ||
			list.SanitizedName() == sanitizedQuery {
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

func (l List) TCliID() string {
	return toTCliID(l.SanitizedName(), l.ID)
}

func (l List) SanitizedName() string {
	return sanitize(l.Name)
}

func FindCard(cards Cards, query string) *Card {
	sanitizedQuery := sanitize(query)
	for _, card := range cards {
		if card.TCliID() == sanitizedQuery ||
			card.ID == query ||
			card.ShortLink == query ||
			card.Name == query ||
			card.SanitizedName() == sanitizedQuery {
			return &card
		}
	}
	return nil
}

type Cards []Card
type Card struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Desc      string  `json:"desc"`
	IDBoard   string  `json:"idBoard"`
	IDList    string  `json:"idList"`
	Closed    bool    `json:"closed"`
	ShortLink string  `json:"shortLink"`
	ShortURL  string  `json:"shortUrl"`
	Pos       float64 `json:"pos"`
	Labels    `json:"labels"`
}

func (c Card) TCliID() string {
	return toTCliID(c.SanitizedName(), c.ShortLink)
}

func (c Card) SanitizedName() string {
	return sanitize(c.Name)
}

type Labels []Label

func (l Labels) String() string {
	if len(l) == 0 {
		return ""
	}
	sb := strings.Builder{}
	for i := 0; i < len(l); i++ {
		sb.WriteString(l[i].ID)
		if i < len(l)-1 {
			sb.WriteString(",")
		}
	}
	return sb.String()
}

type LabelFilter func(s string, l Label) bool

var LabelFilterByID = func(s string, l Label) bool {
	return s == l.ID
}
var LabelFilterByTCliColor = func(s string, l Label) bool {
	return s == l.ToTCliColor()
}
var LabelFilterByColor = func(s string, l Label) bool {
	return s == l.Color
}
var LabelFilterOr = func(filters ...LabelFilter) LabelFilter {
	return func(s string, l Label) bool {
		for _, filter := range filters {
			if filter(s, l) {
				return true
			}
		}
		return false
	}
}

func (l Labels) FilterBy(labels []string, filter LabelFilter) Labels {
	filtered := Labels{}
	for _, s := range labels {
		for _, label := range l {
			if filter(s, label) {
				filtered = append(filtered, label)
				continue
			}
		}
	}
	return filtered
}

func (l Labels) IDLabelsInString() string {
	var idLabels []string
	for _, label := range l {
		idLabels = append(idLabels, label.ID)
	}
	return strings.Join(idLabels, ",")
}

func (l Labels) ToSliceTCliColors() []string {
	s := make([]string, len(l))
	for i := 0; i < len(l); i++ {
		s[i] = l[i].ToTCliColor()
	}
	return s
}

type Label struct {
	ID      string `json:"id"`
	IDBoard string `json:"idBoard"`
	Name    string `json:"name"`
	Color   string `json:"color"`
}

func (l Label) ToTCliColor() string {
	if l.Name == "" {
		return l.Color
	}
	return fmt.Sprintf("%s [%s]", l.Color, l.Name)
}

// CARD CREATION ---------------------------------------------------------------------------------------

// CreateCard represents the resources used to create a new card
// See https://developer.atlassian.com/cloud/trello/rest/api-group-cards/#api-cards-post for more info
type CreateCard struct {
	Name     string      `json:"name"`
	Desc     string      `json:"desc"`
	IDList   string      `json:"idList"`
	IDLabels string      `json:"idLabels,omitempty"`
	Closed   bool        `json:"closed,omitempty"`
	Pos      interface{} `json:"pos,omitempty"` // "top", "bottom" or a positive float
}

func NewCreateCard(card Card) CreateCard {
	return CreateCard{
		Name:   card.Name,
		Desc:   card.Desc,
		IDList: card.IDList,
		Closed: card.Closed,
		Pos:    card.Pos,
	}
}

func NewCardToCreate(card Card) CardToCreate {
	return CardToCreate{
		Name:   card.Name,
		Desc:   card.Desc,
		IDList: card.IDList,
		Pos:    strconv.FormatFloat(card.Pos, 'f', 2, 64),
	}
}

// CardToCreate is the representation used in the card creation
// it's different from the other card representation because we do not want to expose everything to the user
// like for instance, the card ID as the user
type CardToCreate struct {
	Name   string   `yaml:"name"`
	Desc   string   `yaml:"desc"`
	IDList string   `yaml:"idList"`
	Pos    string   `yaml:"pos,omitempty"` // "top", "bottom" or a positive float
	Labels []string `yaml:"labels"`
}

func (ctc CardToCreate) GetPos() interface{} {
	return getPos(ctc.Pos)
}

// CARD UPDATE ---------------------------------------------------------------------------------------

// UpdateCard represents the resources used to update a card
// See https://developer.atlassian.com/cloud/trello/rest/api-group-cards/#api-cards-id-put for more info
type UpdateCard struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Desc     string      `json:"desc"`
	IDBoard  string      `json:"idBoard"`
	IDList   string      `json:"idList"`
	IDLabels string      `json:"idLabels,omitempty"`
	Closed   bool        `json:"closed,omitempty"`
	Pos      interface{} `json:"pos,omitempty"` // "top", "bottom" or a positive float
}

func NewUpdateCard(card Card) UpdateCard {
	return UpdateCard{
		ID:       card.ID,
		Name:     card.Name,
		Desc:     card.Desc,
		IDBoard:  card.IDBoard,
		IDList:   card.IDList,
		IDLabels: card.Labels.String(),
		Closed:   card.Closed,
		Pos:      card.Pos,
	}
}

func NewCardToEdit(card Card) CardToEdit {
	return CardToEdit{
		Name:   card.Name,
		Desc:   card.Desc,
		Closed: card.Closed,
		IDList: card.IDList,
		Labels: card.Labels.ToSliceTCliColors(),
		Pos:    strconv.FormatFloat(card.Pos, 'f', 2, 64),
	}
}

// CardToEdit is the representation used in the card edition
// it's different from the other card representation because we do not want to expose everything to the user
// like for instance, the card ID as the user
type CardToEdit struct {
	Name   string   `yaml:"name"`
	Desc   string   `yaml:"desc"`
	Closed bool     `yaml:"closed"`
	IDList string   `yaml:"idList"`
	Labels []string `yaml:"labels"`
	Pos    string   `yaml:"pos,omitempty"` // "top", "bottom" or a positive float
}

func (cte CardToEdit) GetPos() interface{} {
	return getPos(cte.Pos)
}

// PRIVATE ---------------------------------------------------------------------------

func sanitize(name string) string {
	return strings.ReplaceAll(name, " ", "\\ ")
}

// getPos convert the given pos into appropriate type supported by Trello: either a string or a float
func getPos(in string) interface{} {
	if in == "top" || in == "bottom" {
		return in
	}
	pos, err := strconv.ParseFloat(in, 64)
	if err != nil {
		log.Debug().
			Str("pos", in).
			Err(err).
			Msg("could not parse pos to float")
		return in
	}
	return pos
}

// toTCliID converts a Trello entity into a unique ID understandable by tcli
// it's using the name, for the user experience in the completion, and the short link
// instead of the id to prevent having long lines in the completion.
func toTCliID(name, uid string) string {
	return fmt.Sprintf("%s[%s]", name, uid)
}
