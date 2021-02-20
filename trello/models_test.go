package trello

import (
	"fmt"
	"github.com/logrusorgru/aurora/v3"
	"testing"
)

func TestFindBoards(t *testing.T) {
	nbBoards := 10
	boards := make(Boards, nbBoards)
	for i := 0; i < nbBoards; i++ {
		boards[i] = Board{
			ID:        fmt.Sprintf("id %d", i),
			Name:      fmt.Sprintf("name %d", i),
			ShortLink: fmt.Sprintf("shortLink %d", i),
		}
	}

	var tests = map[string]struct {
		given    string
		expected *Board
	}{
		"find by TCliID": {
			given:    "name 5 [shortLink 5]",
			expected: &boards[5],
		},
		"find by ID": {
			given:    "id 8",
			expected: &boards[8],
		},
		"find by ShortLink": {
			given:    "shortLink 1",
			expected: &boards[1],
		},
		"find by Name": {
			given:    "name 3",
			expected: &boards[3],
		},
		"board not found": {
			given:    "unknown board",
			expected: nil,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := FindBoard(boards, tt.given)
			if tt.expected != nil && actual == nil || tt.expected == nil && actual != nil {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
				t.FailNow()
			}
			if tt.expected != nil {
				if *actual != *tt.expected {
					t.Errorf("expected %v, actual %v", tt.expected, actual)
				}
			}
		})
	}
}

func TestFindLists(t *testing.T) {
	nbLists := 10
	lists := make(Lists, nbLists)
	for i := 0; i < nbLists; i++ {
		lists[i] = List{
			ID:   fmt.Sprintf("id %d", i),
			Name: fmt.Sprintf("name %d", i),
		}
	}

	var tests = map[string]struct {
		given    string
		expected *List
	}{
		"find by ID": {
			given:    "id 8",
			expected: &lists[8],
		},
		"find by Name": {
			given:    "name 3",
			expected: &lists[3],
		},
		"list not found": {
			given:    "unknown list",
			expected: nil,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := FindList(lists, tt.given)
			if tt.expected != nil && actual == nil || tt.expected == nil && actual != nil {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
				t.FailNow()
			}
			if tt.expected != nil {
				if *actual != *tt.expected {
					t.Errorf("expected %v, actual %v", tt.expected, actual)
				}
			}
		})
	}
}

func TestFindCards(t *testing.T) {
	nbCards := 10
	cards := make(Cards, nbCards)
	for i := 0; i < nbCards; i++ {
		cards[i] = Card{
			ID:        fmt.Sprintf("id %d", i),
			Name:      fmt.Sprintf("name %d", i),
			ShortLink: fmt.Sprintf("shortLink %d", i),
		}
	}

	var tests = map[string]struct {
		given    string
		expected *Card
	}{
		"find by TCliID": {
			given:    "name 5 [shortLink 5]",
			expected: &cards[5],
		},
		"find by ID": {
			given:    "id 8",
			expected: &cards[8],
		},
		"find by ShortLink": {
			given:    "shortLink 1",
			expected: &cards[1],
		},
		"find by Name": {
			given:    "name 3",
			expected: &cards[3],
		},
		"card not found": {
			given:    "unknown card",
			expected: nil,
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

func TestLabels_String(t *testing.T) {
	var tests = map[string]struct {
		given    Labels
		expected string
	}{
		"3 labels": {
			given: Labels{
				{ID: "label 1", Name: "label name 1", Color: "red"},
				{ID: "label 2", Name: "label name 2", Color: "sky"},
				{ID: "label 3", Name: "label name 3", Color: "black"},
			},
			expected: "label 1,label 2,label 3",
		},
		"1 label": {
			given: Labels{
				{ID: "label 1", Name: "label name 1", Color: "red"},
			},
			expected: "label 1",
		},
		"no label": {
			given:    Labels{},
			expected: "",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.given.String()
			if actual != tt.expected {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}

func TestLabels_Slice(t *testing.T) {
	var tests = map[string]struct {
		given    Labels
		expected []string
	}{
		"3 labels": {
			given: Labels{
				{ID: "label 1", Name: "label name 1", Color: "red"},
				{ID: "label 2", Name: "label name 2", Color: "sky"},
				{ID: "label 3", Name: "label name 3", Color: "black"},
			},
			expected: []string{"label 1", "label 2", "label 3"},
		},
		"1 label": {
			given: Labels{
				{ID: "label 1", Name: "label name 1", Color: "red"},
			},
			expected: []string{"label 1"},
		},
		"no label": {
			given:    Labels{},
			expected: []string{},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.given.Slice()
			if len(actual) != len(tt.expected) {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
				t.FailNow()
			}
			for i := 0; i < len(actual); i++ {
				if actual[i] != tt.expected[i] {
					t.Errorf("%d: expected %v, actual %v", i, tt.expected[i], actual[i])
				}
			}
		})
	}
}

func TestLabel_Colorize(t *testing.T) {
	s := "content"
	var tests = map[string]struct {
		given    Label
		expected aurora.Value
	}{
		"lime": {
			given:    Label{Color: "lime"},
			expected: aurora.BgGreen(s),
		},
		"no color found": {
			given:    Label{Color: "unknown color"},
			expected: aurora.White(s),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.given.Colorize(s)
			if actual != tt.expected {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}

// CARD CREATION ---------------------------------------------------------------------------------------

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
	actual := NewCardToCreate(c)

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

func TestCardToCreate_GetPos(t *testing.T) {
	var tests = map[string]struct {
		given    CardToCreate
		expected interface{}
	}{
		"top": {
			given:    CardToCreate{Pos: "top"},
			expected: "top",
		},
		"bottom": {
			given:    CardToCreate{Pos: "top"},
			expected: "top",
		},
		"int number": {
			given:    CardToCreate{Pos: "1234"},
			expected: float64(1234),
		},
		"float number": {
			given:    CardToCreate{Pos: "1234.56"},
			expected: 1234.56,
		},
		"unknown value": {
			given:    CardToCreate{Pos: "unknown"},
			expected: "unknown",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.given.GetPos()
			if actual != tt.expected {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}

func TestCardToCreate_IDLabelsInString(t *testing.T) {
	var tests = map[string]struct {
		given    CardToCreate
		expected string
	}{
		"3 labels": {
			given:    CardToCreate{IDLabels: []string{"label 1", "label 2", "label 3"}},
			expected: "label 1,label 2,label 3",
		},
		"1 label": {
			given:    CardToCreate{IDLabels: []string{"label 1"}},
			expected: "label 1",
		},
		"no label": {
			given:    CardToCreate{IDLabels: []string{}},
			expected: "",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.given.IDLabelsInString()
			if actual != tt.expected {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}

// CARD UPDATE ---------------------------------------------------------------------------------------

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
		Labels: Labels{{ID: "label id 1"}, {ID: "label id 2"}, {ID: "label id 3"}},
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
	expectedIDLabels := []string{"label id 1", "label id 2", "label id 3"}
	if len(actual.IDLabels) != len(expectedIDLabels) {
		t.Errorf("expected %v, actual %v", expectedIDLabels, actual.IDLabels)
		t.FailNow()
	}
	for i := 0; i < len(expectedIDLabels); i++ {
		if expectedIDLabels[i] != actual.IDLabels[i] {
			t.Errorf("%d: expected %v, actual %v", i, expectedIDLabels[i], actual.IDLabels[i])
		}
	}
}

func TestCardToEdit_GetPos(t *testing.T) {
	var tests = map[string]struct {
		given    CardToEdit
		expected interface{}
	}{
		"top": {
			given:    CardToEdit{Pos: "top"},
			expected: "top",
		},
		"bottom": {
			given:    CardToEdit{Pos: "top"},
			expected: "top",
		},
		"int number": {
			given:    CardToEdit{Pos: "1234"},
			expected: float64(1234),
		},
		"float number": {
			given:    CardToEdit{Pos: "1234.56"},
			expected: 1234.56,
		},
		"unknown value": {
			given:    CardToEdit{Pos: "unknown"},
			expected: "unknown",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.given.GetPos()
			if actual != tt.expected {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}

func TestCardToEdit_IDLabelsInString(t *testing.T) {
	var tests = map[string]struct {
		given    CardToEdit
		expected string
	}{
		"3 labels": {
			given:    CardToEdit{IDLabels: []string{"label 1", "label 2", "label 3"}},
			expected: "label 1,label 2,label 3",
		},
		"1 label": {
			given:    CardToEdit{IDLabels: []string{"label 1"}},
			expected: "label 1",
		},
		"no label": {
			given:    CardToEdit{IDLabels: []string{}},
			expected: "",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.given.IDLabelsInString()
			if actual != tt.expected {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}
