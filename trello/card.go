package trello

import (
	"sort"
	"strconv"
)

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

func (c Cards) SortedByPos() Cards {
	sort.Slice(c, func(i, j int) bool {
		return c[i].Pos < c[j].Pos
	})
	return c
}

type Card struct {
	ID        string  `json:"id"        toml:"id"`
	Name      string  `json:"name"      toml:"name"`
	Desc      string  `json:"desc"      toml:"desc"`
	IDBoard   string  `json:"idBoard"   toml:"idBoard"`
	IDList    string  `json:"idList"    toml:"idList"`
	Closed    bool    `json:"closed"    toml:"closed"`
	ShortLink string  `json:"shortLink" toml:"shortLink"`
	ShortURL  string  `json:"shortUrl"  toml:"shortUrl"`
	Pos       float64 `json:"pos"       toml:"pos"`
	Labels    `json:"labels" toml:"labels"`
}

func (c Card) TCliID() string {
	return toTCliID(c.SanitizedName(), c.ShortLink)
}

func (c Card) SanitizedName() string {
	return sanitize(c.Name)
}

// CARD CREATION ---------------------------------------------------------------------------------------

// CreateCard represents the resources used to create a new card
// See https://developer.atlassian.com/cloud/trello/rest/api-group-cards/#api-cards-post for more info
type CreateCard struct {
	Name     string      `json:"name"               toml:"name"`
	Desc     string      `json:"desc"               toml:"desc"`
	IDList   string      `json:"idList"             toml:"idList"`
	IDLabels string      `json:"idLabels,omitempty" toml:"idLabels,omitempty"`
	Closed   bool        `json:"closed,omitempty"   toml:"closed,omitempty"`
	Pos      interface{} `json:"pos,omitempty"      toml:"pos,omitempty"` // "top", "bottom" or a positive float
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

func NewCardToCreate(card Card, defaultLabels []string) CardToCreate {
	return CardToCreate{
		Name:   card.Name,
		Desc:   card.Desc,
		IDList: card.IDList,
		Pos:    strconv.FormatFloat(card.Pos, 'f', 2, 64),
		Labels: defaultLabels,
	}
}

// CardToCreate is the representation used in the card creation
// it's different from the other card representation because we do not want to expose everything to the user
// like for instance, the card ID as the user
type CardToCreate struct {
	Name   string   `yaml:"name"          toml:"name"`
	Desc   string   `yaml:"desc"          toml:"desc"`
	IDList string   `yaml:"idList"        toml:"idList"`
	Pos    string   `yaml:"pos,omitempty" toml:"pos,omitempty"` // "top", "bottom" or a positive float
	Labels []string `yaml:"labels"        toml:"labels"`
}

func (ctc CardToCreate) GetPos() interface{} {
	return getPos(ctc.Pos)
}

// CARD UPDATE ---------------------------------------------------------------------------------------

// UpdateCard represents the resources used to update a card
// See https://developer.atlassian.com/cloud/trello/rest/api-group-cards/#api-cards-id-put for more info
type UpdateCard struct {
	ID       string      `json:"id"                 toml:"id"`
	Name     string      `json:"name"               toml:"name"`
	Desc     string      `json:"desc"               toml:"desc"`
	IDBoard  string      `json:"idBoard"            toml:"idBoard"`
	IDList   string      `json:"idList"             toml:"idList"`
	IDLabels string      `json:"idLabels,omitempty" toml:"idLabels,omitempty"`
	Closed   bool        `json:"closed,omitempty"   toml:"closed,omitempty"`
	Pos      interface{} `json:"pos,omitempty"      toml:"pos,omitempty"` // "top", "bottom" or a positive float
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
	Name   string   `yaml:"name"          toml:"name"`
	Desc   string   `yaml:"desc"          toml:"desc"`
	Closed bool     `yaml:"closed"        toml:"closed"`
	IDList string   `yaml:"idList"        toml:"idList"`
	Labels []string `yaml:"labels"        toml:"labels"`
	Pos    string   `yaml:"pos,omitempty" toml:"pos,omitempty"` // "top", "bottom" or a positive float
}

func (cte CardToEdit) GetPos() interface{} {
	return getPos(cte.Pos)
}
