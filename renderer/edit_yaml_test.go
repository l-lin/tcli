package renderer

import (
	"github.com/l-lin/tcli/trello"
	"testing"
)

func TestEditInPrettyYaml_MarshalCardToCreate(t *testing.T) {
	type given struct {
		ctc        trello.CardToCreate
		boardLists trello.Lists
		labels     trello.Labels
	}
	type expected struct {
		hasError bool
		content  string
	}
	var tests = map[string]struct {
		given    given
		expected expected
	}{
		"card with long description": {
			given: given{
				ctc: trello.CardToCreate{
					Name:   "card",
					IDList: "list 1",
				},
				boardLists: trello.Lists{
					{ID: "list 1", Name: "list name 1"},
					{ID: "list 2", Name: "list name 2"},
					{ID: "list 3", Name: "list name 3"},
				},
				labels: trello.Labels{
					{ID: "id red", Name: "name red", Color: "red"},
					{ID: "id sky", Name: "name sky", Color: "sky"},
					{ID: "id black", Name: "", Color: "black"},
				},
			},
			expected: expected{
				hasError: false,
				content: `name: "card"
# available lists:
# list 1: list name 1
# list 2: list name 2
# list 3: list name 3
idList: "list 1"
# the position of the card in its list: "top", "bottom" or a positive float
pos: "bottom"
# available labels (use color or ID):
# id red: red [name red]
# id sky: sky [name sky]
# id black: black
labels: 
  - 
desc: |-
  `,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			e := EditInPrettyYaml{}
			actual, actualErr := e.MarshalCardToCreate(tt.given.ctc, tt.given.boardLists, tt.given.labels)
			if tt.expected.hasError != (actualErr != nil) {
				t.Errorf("expected error %v, actual error %v", tt.expected.hasError, actualErr)
				t.FailNow()
			}
			if string(actual) != tt.expected.content {
				t.Errorf("expected:\n%v\nactual:\n%v", tt.expected.content, string(actual))
			}
		})
	}
}

func TestEditInYaml_RenderCardToEdit(t *testing.T) {
	type given struct {
		cte        trello.CardToEdit
		boardLists trello.Lists
		labels     trello.Labels
	}
	type expected struct {
		hasError bool
		content  string
	}
	labels := trello.Labels{
		{ID: "id red", Name: "name red", Color: "red"},
		{ID: "id sky", Name: "name sky", Color: "sky"},
		{ID: "id black", Name: "", Color: "black"},
	}
	var tests = map[string]struct {
		given    given
		expected expected
	}{
		"card with long description, 3 board lists, 3 board labels": {
			given: given{
				cte: trello.CardToEdit{
					Name: "card",
					Desc: `# card description

> some context

foobar`,
					Closed: false,
					IDList: "list 1",
					Pos:    "123",
					Labels: []string{labels[0].ToTCliColor(), labels[2].ToTCliColor()},
				},
				boardLists: trello.Lists{
					{ID: "list 1", Name: "list name 1"},
					{ID: "list 2", Name: "list name 2"},
					{ID: "list 3", Name: "list name 3"},
				},
				labels: labels,
			},
			expected: expected{
				hasError: false,
				content: `name: "card"
# whether the card should be archived (closed: true)
closed: false
# available lists:
# list 1: list name 1
# list 2: list name 2
# list 3: list name 3
idList: "list 1"
# the position of the card in its list: "top", "bottom" or a positive float
pos: 123
# available labels (use color or ID):
# id red: red [name red]
# id sky: sky [name sky]
# id black: black
labels:
  - red [name red]
  - black
desc: |-
  # card description

  > some context

  foobar
`,
			},
		},
		"card with long description, no board lists, 3 board labels": {
			given: given{
				cte: trello.CardToEdit{
					Name: "card",
					Desc: `# card description

> some context

foobar`,
					Closed: false,
					IDList: "list 1",
					Pos:    "123",
					Labels: []string{labels[0].ToTCliColor(), labels[2].ToTCliColor()},
				},
				boardLists: trello.Lists{},
				labels:     labels,
			},
			expected: expected{
				hasError: false,
				content: `name: "card"
# whether the card should be archived (closed: true)
closed: false
# available lists:
idList: "list 1"
# the position of the card in its list: "top", "bottom" or a positive float
pos: 123
# available labels (use color or ID):
# id red: red [name red]
# id sky: sky [name sky]
# id black: black
labels:
  - red [name red]
  - black
desc: |-
  # card description

  > some context

  foobar
`,
			},
		},
		"card with long description, 3 board lists, no board labels": {
			given: given{
				cte: trello.CardToEdit{
					Name: "card",
					Desc: `# card description

> some context

foobar`,
					Closed: false,
					IDList: "list 1",
					Pos:    "123",
					Labels: []string{},
				},
				boardLists: trello.Lists{
					{ID: "list 1", Name: "list name 1"},
					{ID: "list 2", Name: "list name 2"},
					{ID: "list 3", Name: "list name 3"},
				},
				labels: trello.Labels{},
			},
			expected: expected{
				hasError: false,
				content: `name: "card"
# whether the card should be archived (closed: true)
closed: false
# available lists:
# list 1: list name 1
# list 2: list name 2
# list 3: list name 3
idList: "list 1"
# the position of the card in its list: "top", "bottom" or a positive float
pos: 123
# available labels (use color or ID):
labels:
desc: |-
  # card description

  > some context

  foobar
`,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			e := EditInPrettyYaml{}
			actual, actualErr := e.MarshalCardToEdit(tt.given.cte, tt.given.boardLists, tt.given.labels)
			if tt.expected.hasError != (actualErr != nil) {
				t.Errorf("expected error %v, actual error %v", tt.expected.hasError, actualErr)
				t.FailNow()
			}
			if string(actual) != tt.expected.content {
				t.Errorf("expected:\n%v\nactual:\n%v", tt.expected.content, string(actual))
			}
		})
	}
}
