package completer

import (
	"github.com/c-bata/go-prompt"
	"github.com/l-lin/tcli/trello"
	"testing"
)

func TestBoardsToSuggestions(t *testing.T) {
	var tests = map[string]struct {
		given    trello.Boards
		expected []prompt.Suggest
	}{
		"two boards": {
			given: trello.Boards{
				{Name: "board 1", ShortLink: "shortLink 1"},
				{Name: "board 2", ShortLink: "shortLink 2"},
			},
			expected: []prompt.Suggest{
				{Text: "board 1 [shortLink 1]"},
				{Text: "board 2 [shortLink 2]"},
			},
		},
		"no board": {
			given:    trello.Boards{},
			expected: []prompt.Suggest{},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := boardsToSuggestions(tt.given)
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

func TestListsToSuggestions(t *testing.T) {
	var tests = map[string]struct {
		given    trello.Lists
		expected []prompt.Suggest
	}{
		"two lists": {
			given: trello.Lists{
				{Name: "list 1"},
				{Name: "list 2"},
			},
			expected: []prompt.Suggest{
				{Text: "list 1"},
				{Text: "list 2"},
			},
		},
		"no list": {
			given:    trello.Lists{},
			expected: []prompt.Suggest{},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := listsToSuggestions(tt.given)
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

func TestCardsToSuggestions(t *testing.T) {
	var tests = map[string]struct {
		given    trello.Cards
		expected []prompt.Suggest
	}{
		"two cards": {
			given: trello.Cards{
				{Name: "card 1", Desc: "card description 1", ShortLink: "shortLink 1"},
				{Name: "card 2", Desc: "card description 2", ShortLink: "shortLink 2"},
			},
			expected: []prompt.Suggest{
				{Text: "card 1 [shortLink 1]", Description: "card description 1"},
				{Text: "card 2 [shortLink 2]", Description: "card description 2"},
			},
		},
		"no card": {
			given:    trello.Cards{},
			expected: []prompt.Suggest{},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := cardsToSuggestions(tt.given)
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

func TestTruncateCardDescription(t *testing.T) {
	var tests = map[string]struct {
		given    string
		expected string
	}{
		"empty description": {
			given:    "",
			expected: "",
		},
		"long description": {
			given:    "long description that exceed the threshold",
			expected: "long description tha",
		},
		"short description": {
			given:    "short desc",
			expected: "short desc",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := truncateCardDescription(tt.given)
			if actual != tt.expected {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}
