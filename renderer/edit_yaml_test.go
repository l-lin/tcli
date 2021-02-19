package renderer

import (
	"github.com/l-lin/tcli/trello"
	"testing"
)

func TestEditInYaml_RenderCardToEdit(t *testing.T) {
	type given struct {
		cte        trello.CardToEdit
		boardLists trello.Lists
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
				cte: trello.CardToEdit{
					Name: "card",
					Desc: `# card description

> some context

foobar`,
					Closed: false,
					IDList: "list 1",
				},
				boardLists: trello.Lists{
					{ID: "list 1", Name: "list name 1"},
					{ID: "list 2", Name: "list name 2"},
					{ID: "list 3", Name: "list name 3"},
				},
			},
			expected: expected{
				hasError: false,
				content: `name: "card"
# whether the card should be archived (closed: true)
closed: false
# available board lists:
# list 1: list name 1
# list 2: list name 2
# list 3: list name 3
idList: "list 1"
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
			actual, actualErr := e.Marshal(tt.given.cte, tt.given.boardLists)
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

func TestEditInPrettyYaml_Unmarshal(t *testing.T) {
	type expected struct {
		hasError   bool
		editedCard trello.CardToEdit
	}
	var tests = map[string]struct {
		given    []byte
		expected expected
	}{
		"edited card": {
			given: []byte(`name: "card"
# whether the card should be archived (closed: true)
closed: false
# available board lists:
# list 1: list name 1
# list 2: list name 2
# list 3: list name 3
idList: "list 1"
desc: |-
  # card description

  > some context

  foobar
`),
			expected: expected{
				hasError: false,
				editedCard: trello.CardToEdit{
					Name: "card",
					Desc: `# card description

> some context

foobar`,
					Closed: false,
					IDList: "list 1",
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			e := EditInPrettyYaml{}
			var actual trello.CardToEdit
			actualErr := e.Unmarshal(tt.given, &actual)
			if tt.expected.hasError != (actualErr != nil) {
				t.Errorf("expected error %v, actual %v", tt.expected.hasError, actualErr)
				t.FailNow()
			}
			if actual != tt.expected.editedCard {
				t.Errorf("expected %v, actual %v", tt.expected.editedCard, actual)
			}
		})
	}
}
