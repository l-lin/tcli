package trello

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFindCards(t *testing.T) {
	nbCards := 10
	cards := make(Cards, nbCards)
	for i := 0; i < nbCards; i++ {
		cards[i] = Card{
			ID:        fmt.Sprintf("id%d", i),
			Name:      fmt.Sprintf("name%d", i),
			ShortLink: fmt.Sprintf("shortLink%d", i),
		}
	}
	lastCardsIndex := len(cards) - 1
	cards[lastCardsIndex].Name = fmt.Sprintf("name with space %d", lastCardsIndex)

	var tests = map[string]struct {
		given    string
		expected *Card
	}{
		"find by TCliID": {
			given:    "name5[shortLink5]",
			expected: &cards[5],
		},
		"find by ID": {
			given:    "id8",
			expected: &cards[8],
		},
		"find by ShortLink": {
			given:    "shortLink1",
			expected: &cards[1],
		},
		"find by Name": {
			given:    "name3",
			expected: &cards[3],
		},
		"card not found": {
			given:    "unknown-card",
			expected: nil,
		},
		"find by TCliID - card with space in its name": {
			given:    fmt.Sprintf("%s[%s]", cards[lastCardsIndex].Name, cards[lastCardsIndex].ShortLink),
			expected: &cards[lastCardsIndex],
		},
		"find by Name - card with space in its name": {
			given:    cards[lastCardsIndex].Name,
			expected: &cards[lastCardsIndex],
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := FindCard(cards, tt.given)
			if tt.expected != nil && actual == nil || tt.expected == nil && actual != nil {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
				t.FailNow()
			}
			if tt.expected != nil {
				if actual.ID != tt.expected.ID || actual.Name != tt.expected.Name || actual.ShortLink != tt.expected.ShortLink {
					t.Errorf("expected %v, actual %v", tt.expected, actual)
				}
			}
		})
	}
}

func TestCards_SortedByPos(t *testing.T) {
	var tests = map[string]struct {
		given    Cards
		expected Cards
	}{
		"3 cards": {
			given: Cards{
				{ID: "card 1", Pos: 10},
				{ID: "card 2", Pos: 1},
				{ID: "card 3", Pos: 11},
			},
			expected: Cards{
				{ID: "card 2", Pos: 1},
				{ID: "card 1", Pos: 10},
				{ID: "card 3", Pos: 11},
			},
		},
		"no card": {
			given:    Cards{},
			expected: Cards{},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.given.SortedByPos()
			if !reflect.DeepEqual(tt.expected, actual) {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}

func TestNewCreateCard(t *testing.T) {
	// GIVEN
	c := Card{
		Name:   "card name",
		Desc:   "card description",
		IDList: "list id",
		Pos:    123,
	}

	// WHEN
	actual := NewCreateCard(c)

	// THEN
	if c.Name != actual.Name {
		t.Errorf("expected %v, actual %v", c.Name, actual.Name)
	}
	if c.Desc != actual.Desc {
		t.Errorf("expected %v, actual %v", c.Desc, actual.Desc)
	}
	if c.IDList != actual.IDList {
		t.Errorf("expected %v, actual %v", c.IDList, actual.IDList)
	}
	if c.Closed != actual.Closed {
		t.Errorf("expected %v, actual %v", c.Closed, actual.Closed)
	}
	if c.Pos != actual.Pos {
		t.Errorf("expected %v, actual %v", c.Pos, actual.Pos)
	}
}

func TestNewCardToCreate(t *testing.T) {
	// GIVEN
	c := Card{
		Name:   "card name",
		Desc:   "card description",
		IDList: "list id",
		Pos:    123,
	}

	// WHEN
	actual := NewCardToCreate(c, []string{})

	// THEN
	if c.Name != actual.Name {
		t.Errorf("expected %v, actual %v", c.Name, actual.Name)
	}
	if c.Desc != actual.Desc {
		t.Errorf("expected %v, actual %v", c.Desc, actual.Desc)
	}
	if c.IDList != actual.IDList {
		t.Errorf("expected %v, actual %v", c.IDList, actual.IDList)
	}
	expectedPos := "123.00"
	if actual.Pos != expectedPos {
		t.Errorf("expected %v, actual %v", expectedPos, actual.Pos)
	}
}

func TestNewUpdateCard(t *testing.T) {
	// GIVEN
	c := Card{
		ID:      "card id",
		Name:    "card name",
		Desc:    "card description",
		IDBoard: "board id",
		IDList:  "list id",
		Closed:  true,
		Pos:     123,
		Labels:  Labels{{ID: "label id 1"}, {ID: "label id 2"}, {ID: "label id 3"}},
	}

	// WHEN
	actual := NewUpdateCard(c)

	// THEN
	if c.ID != actual.ID {
		t.Errorf("expected %v, actual %v", c.ID, actual.ID)
	}
	if c.Name != actual.Name {
		t.Errorf("expected %v, actual %v", c.Name, actual.Name)
	}
	if c.Desc != actual.Desc {
		t.Errorf("expected %v, actual %v", c.Desc, actual.Desc)
	}
	if c.IDBoard != actual.IDBoard {
		t.Errorf("expected %v, actual %v", c.IDBoard, actual.IDBoard)
	}
	if c.IDList != actual.IDList {
		t.Errorf("expected %v, actual %v", c.IDList, actual.IDList)
	}
	if c.Closed != actual.Closed {
		t.Errorf("expected %v, actual %v", c.Closed, actual.Closed)
	}
	if c.Pos != actual.Pos {
		t.Errorf("expected %v, actual %v", c.Pos, actual.Pos)
	}
	expectedIDLabels := "label id 1,label id 2,label id 3"
	if actual.IDLabels != expectedIDLabels {
		t.Errorf("expected %v, actual %v", expectedIDLabels, actual.IDLabels)
	}
}

func TestNewCardToEdit(t *testing.T) {
	// GIVEN
	c := Card{
		Name:   "card name",
		Desc:   "card description",
		IDList: "list id",
		Closed: true,
		Pos:    123,
		Labels: Labels{
			{ID: "id red", Color: "red", Name: "name red"},
			{ID: "id green", Color: "green", Name: "name green"},
			{ID: "id black", Color: "black"},
		},
	}

	// WHEN
	actual := NewCardToEdit(c)

	// THEN
	if c.Name != actual.Name {
		t.Errorf("expected %v, actual %v", c.Name, actual.Name)
	}
	if c.Desc != actual.Desc {
		t.Errorf("expected %v, actual %v", c.Desc, actual.Desc)
	}
	if c.IDList != actual.IDList {
		t.Errorf("expected %v, actual %v", c.IDList, actual.IDList)
	}
	if c.Closed != actual.Closed {
		t.Errorf("expected %v, actual %v", c.Closed, actual.Closed)
	}
	expectedPos := "123.00"
	if actual.Pos != expectedPos {
		t.Errorf("expected %v, actual %v", expectedPos, actual.Pos)
	}
	expectedIDsLabel := []string{"red [name red]", "green [name green]", "black"}
	if !reflect.DeepEqual(expectedIDsLabel, actual.Labels) {
		t.Errorf("expected %v, actual %v", expectedIDsLabel, actual.Labels)
	}
}
