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
