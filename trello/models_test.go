package trello

import (
	"fmt"
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
