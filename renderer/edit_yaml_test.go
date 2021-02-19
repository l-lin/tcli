package renderer

import (
	"github.com/l-lin/tcli/trello"
	"testing"
)

func TestEditInPrettyYaml_MarshalCardToCreate(t *testing.T) {
	type given struct {
		ctc        trello.CardToCreate
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
				ctc: trello.CardToCreate{
					Name:   "card",
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
# available board lists:
# list 1: list name 1
# list 2: list name 2
# list 3: list name 3
idList: "list 1"
# the position of the card in its list: "top", "bottom" or a positive float
pos: "bottom"
desc: |-
  `,
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			e := EditInPrettyYaml{}
			actual, actualErr := e.MarshalCardToCreate(tt.given.ctc, tt.given.boardLists)
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
					Pos:    "123",
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
# the position of the card in its list: "top", "bottom" or a positive float
pos: 123
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
			actual, actualErr := e.MarshalCardToEdit(tt.given.cte, tt.given.boardLists)
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
	e := EditInPrettyYaml{}

	var tests = map[string]struct {
		given []byte
		when  func(in []byte) (interface{}, error)
		then  func(actual interface{}, err error)
	}{
		"created card": {
			given: []byte(`name: "card"
# available board lists:
# list 1: list name 1
# list 2: list name 2
# list 3: list name 3
idList: "list 1"
# the position of the card in its list: "top", "bottom" or a positive float
pos: "top"
desc: |-
  # card description

  > some context

  foobar
`),
			when: func(in []byte) (interface{}, error) {
				var ctc trello.CardToCreate
				err := e.Unmarshal(in, &ctc)
				return ctc, err
			},
			then: func(actual interface{}, err error) {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
					t.FailNow()
				}
				expected := trello.CardToCreate{
					Name: "card",
					Desc: `# card description

> some context

foobar`,
					IDList: "list 1",
					Pos:    "top",
				}
				if actual != expected {
					t.Errorf("expected %v, actual %v", expected, actual)
				}
			},
		},
		"edited card": {
			given: []byte(`name: "card"
# whether the card should be archived (closed: true)
closed: false
# available board lists:
# list 1: list name 1
# list 2: list name 2
# list 3: list name 3
idList: "list 1"
# the position of the card in its list: "top", "bottom" or a positive float
pos: "top"
desc: |-
  # card description

  > some context

  foobar
`),
			when: func(in []byte) (interface{}, error) {
				var ctc trello.CardToEdit
				err := e.Unmarshal(in, &ctc)
				return ctc, err
			},
			then: func(actual interface{}, err error) {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
					t.FailNow()
				}
				expected := trello.CardToEdit{
					Name: "card",
					Desc: `# card description

> some context

foobar`,
					Closed: false,
					IDList: "list 1",
					Pos:    "top",
				}
				if actual != expected {
					t.Errorf("expected %v, actual %v", expected, actual)
				}
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.then(tt.when(tt.given))
		})
	}
}
