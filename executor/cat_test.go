package executor

import (
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/l-lin/tcli/renderer"
	"github.com/l-lin/tcli/trello"
	"testing"
)

func TestCat_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type given struct {
		arg                   string
		buildTrelloRepository func() trello.Repository
		buildRenderer         func() renderer.Renderer
	}
	type expected struct {
		output    string
		errOutput string
	}
	var tests = map[string]struct {
		given    given
		expected expected
	}{
		"no arg": {
			given: given{
				arg:                   "",
				buildTrelloRepository: func() trello.Repository { return nil },
				buildRenderer:         func() renderer.Renderer { return nil },
			},
			expected: expected{output: ""},
		},
		"show board info": {
			given: given{
				arg: "board",
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("board").
						Return(&trello.Board{ID: "board 1", Name: "board"}, nil)
					return tr
				},
				buildRenderer: func() renderer.Renderer {
					r := renderer.NewMockRenderer(ctrl)
					r.EXPECT().
						RenderBoard(trello.Board{ID: "board 1", Name: "board"}).
						Return("board content")
					return r
				},
			},
			expected: expected{output: "board content\n"},
		},
		"show list info": {
			given: given{
				arg: "board/list",
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("board").
						Return(&trello.Board{ID: "board 1", Name: "board"}, nil)
					tr.EXPECT().
						FindList("board 1", "list").
						Return(&trello.List{ID: "list 1", Name: "list"}, nil)
					return tr
				},
				buildRenderer: func() renderer.Renderer {
					r := renderer.NewMockRenderer(ctrl)
					r.EXPECT().
						RenderList(trello.List{ID: "list 1", Name: "list"}).
						Return("list content")
					return r
				},
			},
			expected: expected{output: "list content\n"},
		},
		"show card info": {
			given: given{
				arg: "board/list/card",
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("board").
						Return(&trello.Board{ID: "board 1", Name: "board"}, nil)
					tr.EXPECT().
						FindList("board 1", "list").
						Return(&trello.List{ID: "list 1", Name: "list"}, nil)
					tr.EXPECT().
						FindCard("list 1", "card").
						Return(&trello.Card{ID: "card 1", Name: "card"}, nil)
					return tr
				},
				buildRenderer: func() renderer.Renderer {
					r := renderer.NewMockRenderer(ctrl)
					r.EXPECT().
						RenderCard(trello.Card{ID: "card 1", Name: "card"}).
						Return("card content")
					return r
				},
			},
			expected: expected{output: "card content\n"},
		},
		// ERRORS
		"unknown-board": {
			given: given{
				arg: "unknown-board",
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("unknown-board").
						Return(nil, errors.New("not found"))
					return tr
				},
				buildRenderer: func() renderer.Renderer {
					return nil
				},
			},
			expected: expected{
				errOutput: "no board found with name 'unknown-board'\n",
			},
		},
		"unknown-list": {
			given: given{
				arg: "board/unknown-list",
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("board").
						Return(&trello.Board{ID: "board 1", Name: "board"}, nil)
					tr.EXPECT().
						FindList("board 1", "unknown-list").
						Return(nil, errors.New("not found"))
					return tr
				},
				buildRenderer: func() renderer.Renderer {
					return nil
				},
			},
			expected: expected{
				errOutput: "no list found with name 'unknown-list'\n",
			},
		},
		"unknown-card": {
			given: given{
				arg: "board/list/unknown-card",
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("board").
						Return(&trello.Board{ID: "board 1", Name: "board"}, nil)
					tr.EXPECT().
						FindList("board 1", "list").
						Return(&trello.List{ID: "list 1", Name: "list"}, nil)
					tr.EXPECT().
						FindCard("list 1", "unknown-card").
						Return(nil, errors.New("not found"))
					return tr
				},
				buildRenderer: func() renderer.Renderer {
					return nil
				},
			},
			expected: expected{
				errOutput: "no card found with name 'unknown-card'\n",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			outputBuf := bytes.Buffer{}
			errOutputBuf := bytes.Buffer{}
			c := cat{
				executor{
					tr:        tt.given.buildTrelloRepository(),
					r:         tt.given.buildRenderer(),
					output:    &outputBuf,
					errOutput: &errOutputBuf,
				},
			}
			c.Execute(tt.given.arg)

			actualOutput := outputBuf.String()
			if actualOutput != tt.expected.output {
				t.Errorf("expected output %v, actual output %v", tt.expected.output, actualOutput)
			}
			actualErrOutput := errOutputBuf.String()
			if actualErrOutput != tt.expected.errOutput {
				t.Errorf("expected errOutput %v, actual errOutput %v", tt.expected.errOutput, actualErrOutput)
			}
		})
	}
}
