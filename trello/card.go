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
