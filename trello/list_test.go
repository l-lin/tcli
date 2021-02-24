package trello

import (
	"fmt"
	"testing"
)

func TestFindLists(t *testing.T) {
	nbLists := 10
	lists := make(Lists, nbLists)
	for i := 0; i < nbLists; i++ {
		lists[i] = List{
			ID:   fmt.Sprintf("id%d", i),
			Name: fmt.Sprintf("name%d", i),
		}
	}
	lastListsIndex := len(lists) - 1
	lists[lastListsIndex].Name = fmt.Sprintf("name with space %d", lastListsIndex)

	var tests = map[string]struct {
		given    string
		expected *List
	}{
		"find by TCliID": {
			given:    "name6[id6]",
			expected: &lists[6],
		},
		"find by ID": {
			given:    "id8",
			expected: &lists[8],
		},
		"find by Name": {
			given:    "name3",
			expected: &lists[3],
		},
		"list not found": {
			given:    "unknown-list",
			expected: nil,
		},
		"find by name - list with space in its name": {
			given:    lists[lastListsIndex].Name,
			expected: &lists[lastListsIndex],
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

func TestList_SanitizedName(t *testing.T) {
	var tests = map[string]struct {
		given    List
		expected string
	}{
		"no special character": {
			given:    List{Name: "someList"},
			expected: "someList",
		},
		"containing spaces": {
			given:    List{Name: "some list"},
			expected: `some\ list`,
		},
		"containing unicodes": {
			given:    List{Name: "🎉 some list"},
			expected: `🎉\ some\ list`,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.given.SanitizedName()
			if actual != tt.expected {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}
