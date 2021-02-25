package trello

import (
	"fmt"
	"reflect"
	"testing"
)

func TestComments_SortedByDateDesc(t *testing.T) {
	var tests = map[string]struct {
		given    Comments
		expected Comments
	}{
		"3 comments with well formatted dates": {
			given: Comments{
				{ID: "comment 1", Date: "2014-11-12T11:45:26.371Z"},
				{ID: "comment 2", Date: "2014-12-12T11:45:26.371Z"},
				{ID: "comment 3", Date: "2013-11-12T11:45:26.371Z"},
			},
			expected: Comments{
				{ID: "comment 2", Date: "2014-12-12T11:45:26.371Z"},
				{ID: "comment 1", Date: "2014-11-12T11:45:26.371Z"},
				{ID: "comment 3", Date: "2013-11-12T11:45:26.371Z"},
			},
		},
		"3 comments with well formatted dates and one with wrong date": {
			given: Comments{
				{ID: "comment 1", Date: "2014-11-12T11:45:26.371Z"},
				{ID: "comment 2", Date: "un-parsable"},
				{ID: "comment 3", Date: "2013-11-12T11:45:26.371Z"},
				{ID: "comment 4", Date: "2014-12-12T11:45:26.371Z"},
			},
			expected: Comments{
				{ID: "comment 4", Date: "2014-12-12T11:45:26.371Z"},
				{ID: "comment 1", Date: "2014-11-12T11:45:26.371Z"},
				{ID: "comment 3", Date: "2013-11-12T11:45:26.371Z"},
				{ID: "comment 2", Date: "un-parsable"},
			},
		},
		"3 comments with well formatted dates and one empty content": {
			given: Comments{
				{ID: "comment 1", Date: "2014-11-12T11:45:26.371Z"},
				{ID: "comment 2", Date: ""},
				{ID: "comment 3", Date: "2013-11-12T11:45:26.371Z"},
				{ID: "comment 4", Date: "2014-12-12T11:45:26.371Z"},
			},
			expected: Comments{
				{ID: "comment 4", Date: "2014-12-12T11:45:26.371Z"},
				{ID: "comment 1", Date: "2014-11-12T11:45:26.371Z"},
				{ID: "comment 3", Date: "2013-11-12T11:45:26.371Z"},
				{ID: "comment 2", Date: ""},
			},
		},
		"no comment": {
			given:    Comments{},
			expected: Comments{},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.given.SortedByDateDesc()
			if !reflect.DeepEqual(tt.expected, actual) {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}

func TestFindComment(t *testing.T) {
	nbComments := 10
	comments := make(Comments, nbComments)
	for i := 0; i < nbComments; i++ {
		comments[i] = Comment{
			ID: fmt.Sprintf("comment%d", i),
		}
	}
	var tests = map[string]struct {
		given    string
		expected *Comment
	}{
		"found": {
			given:    "comment8",
			expected: &Comment{ID: "comment8"},
		},
		"not found": {
			given:    "unknown",
			expected: nil,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := FindComment(comments, tt.given)
			if !reflect.DeepEqual(tt.expected, actual) {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}
