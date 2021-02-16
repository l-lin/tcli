package executor

import (
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/l-lin/tcli/renderer"
	"github.com/l-lin/tcli/trello"
	"testing"
)

func TestLs_Execute(t *testing.T) {
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
				arg: "",
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						GetBoards().
						Return(trello.Boards{
							{ID: "board 1", Name: "board"},
							{ID: "board 2", Name: "another board"},
						}, nil)
					return tr
				},
				buildRenderer: func() renderer.Renderer {
					r := renderer.NewMockRenderer(ctrl)
					r.EXPECT().
						RenderBoards(trello.Boards{
							{ID: "board 1", Name: "board"},
							{ID: "board 2", Name: "another board"},
						}).
						Return("boards content")
					return r
				},
			},
			expected: expected{output: "boards content"},
		},
		"show lists": {
			given: given{
				arg: "board",
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("board").
						Return(&trello.Board{ID: "board 1", Name: "board"}, nil)
					tr.EXPECT().
						GetLists("board 1").
						Return(trello.Lists{
							{ID: "list 1", Name: "list"},
							{ID: "list 2", Name: "another list"},
						}, nil)
					return tr
				},
				buildRenderer: func() renderer.Renderer {
					r := renderer.NewMockRenderer(ctrl)
					r.EXPECT().
						RenderLists(trello.Lists{
							{ID: "list 1", Name: "list"},
							{ID: "list 2", Name: "another list"},
						}).
						Return("lists content")
					return r
				},
			},
			expected: expected{output: "lists content"},
		},
		"show cards": {
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
					tr.EXPECT().
						GetCards("list 1").
						Return(trello.Cards{
							{ID: "card 1", Name: "card"},
							{ID: "card 2", Name: "another card"},
						}, nil)
					return tr
				},
				buildRenderer: func() renderer.Renderer {
					r := renderer.NewMockRenderer(ctrl)
					r.EXPECT().
						RenderCards(trello.Cards{
							{ID: "card 1", Name: "card"},
							{ID: "card 2", Name: "another card"},
						}).
						Return("cards content")
					return r
				},
			},
			expected: expected{output: "cards content"},
		},
		"show single card": {
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
						RenderCards(trello.Cards{{ID: "card 1", Name: "card"}}).
						Return("card 1 content")
					return r
				},
			},
			expected: expected{output: "card 1 content"},
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
				errOutput: "no board found with name 'unknown-board'",
			},
		},
		"unknown list": {
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
				errOutput: "no list found with name 'unknown-list'",
			},
		},
		"unknown card": {
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
				errOutput: "no card found with name 'unknown-card'",
			},
		},
		"cannot find boards": {
			given: given{
				arg: "",
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						GetBoards().
						Return(nil, errors.New("unexpected error"))
					return tr
				},
				buildRenderer: func() renderer.Renderer {
					return nil
				},
			},
			expected: expected{
				errOutput: "could not fetch boards: unexpected error",
			},
		},
		"cannot find lists": {
			given: given{
				arg: "board",
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("board").
						Return(&trello.Board{ID: "board 1", Name: "board"}, nil)
					tr.EXPECT().
						GetLists("board 1").
						Return(nil, errors.New("unexpected error"))
					return tr
				},
				buildRenderer: func() renderer.Renderer {
					return nil
				},
			},
			expected: expected{
				errOutput: "could not fetch lists for board 'board': unexpected error",
			},
		},
		"cannot find cards": {
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
					tr.EXPECT().
						GetCards("list 1").
						Return(nil, errors.New("unexpected error"))
					return tr
				},
				buildRenderer: func() renderer.Renderer {
					return nil
				},
			},
			expected: expected{
				errOutput: "could not fetch cards for list 'list': unexpected error",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			outBuf := bytes.Buffer{}
			errBuf := bytes.Buffer{}
			l := ls{
				executor{
					tr:        tt.given.buildTrelloRepository(),
					r:         tt.given.buildRenderer(),
					output:    &outBuf,
					errOutput: &errBuf,
				},
			}
			l.Execute(tt.given.arg)

			actualOutput := outBuf.String()
			if actualOutput != tt.expected.output {
				t.Errorf("expected output %v, actual output %v", tt.expected.output, actualOutput)
			}
			actualErrOutput := errBuf.String()
			if actualErrOutput != tt.expected.errOutput {
				t.Errorf("expected errOutput %v, actual errOutput %v", tt.expected.errOutput, actualErrOutput)
			}
		})
	}
}
