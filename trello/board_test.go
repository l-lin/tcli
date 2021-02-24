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
			ID:        fmt.Sprintf("id%d", i),
			Name:      fmt.Sprintf("name%d", i),
			ShortLink: fmt.Sprintf("shortLink%d", i),
		}
	}
	lastBoardIndex := len(boards) - 1
	boards[lastBoardIndex].Name = fmt.Sprintf("name with space %d", lastBoardIndex)

	var tests = map[string]struct {
		given    string
		expected *Board
	}{
		"find by TCliID": {
			given:    "name5[shortLink5]",
			expected: &boards[5],
		},
		"find by ID": {
			given:    "id8",
			expected: &boards[8],
		},
		"find by ShortLink": {
			given:    "shortLink1",
			expected: &boards[1],
		},
		"find by Name": {
			given:    "name3",
			expected: &boards[3],
		},
		"board not found": {
			given:    "unknown board",
			expected: nil,
		},
		"find by TCliID - board with space in its name": {
			given:    fmt.Sprintf("%s[%s]", boards[lastBoardIndex].Name, boards[lastBoardIndex].ShortLink),
			expected: &boards[lastBoardIndex],
		},
		"find by Name - board with space in its name": {
			given:    boards[lastBoardIndex].Name,
			expected: &boards[lastBoardIndex],
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
