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
		args                  []string
		buildTrelloRepository func() trello.Repository
		buildRenderer         func() renderer.Renderer
	}
	type expected struct {
		stdout string
		stderr string
	}
	board := trello.Board{ID: "board 1", Name: "board"}
	list := trello.List{ID: "list 1", Name: "list"}
	card1 := trello.Card{ID: "card 1", Name: "card"}
	card2 := trello.Card{ID: "card 2", Name: "another-card"}

	var tests = map[string]struct {
		given    given
		expected expected
	}{
		"no arg": {
			given: given{
				args:                  []string{},
				buildTrelloRepository: func() trello.Repository { return nil },
				buildRenderer:         func() renderer.Renderer { return nil },
			},
			expected: expected{stdout: ""},
		},
		"first arg as empty string": {
			given: given{
				args:                  []string{""},
				buildTrelloRepository: func() trello.Repository { return nil },
				buildRenderer:         func() renderer.Renderer { return nil },
			},
			expected: expected{stdout: ""},
		},
		"show board info": {
			given: given{
				args: []string{"board"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil)
					return tr
				},
				buildRenderer: func() renderer.Renderer {
					r := renderer.NewMockRenderer(ctrl)
					r.EXPECT().
						RenderBoard(board).
						Return("board content")
					return r
				},
			},
			expected: expected{stdout: "board content\n"},
		},
		"show list info": {
			given: given{
				args: []string{"board/list"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil)
					tr.EXPECT().
						FindList(board.ID, list.Name).
						Return(&list, nil)
					return tr
				},
				buildRenderer: func() renderer.Renderer {
					r := renderer.NewMockRenderer(ctrl)
					r.EXPECT().
						RenderList(list).
						Return("list content")
					return r
				},
			},
			expected: expected{stdout: "list content\n"},
		},
		"show card info": {
			given: given{
				args: []string{"board/list/card"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil)
					tr.EXPECT().
						FindList(board.ID, list.Name).
						Return(&list, nil)
					tr.EXPECT().
						FindCard(list.ID, card1.Name).
						Return(&card1, nil)
					return tr
				},
				buildRenderer: func() renderer.Renderer {
					r := renderer.NewMockRenderer(ctrl)
					r.EXPECT().
						RenderCard(card1).
						Return("card content")
					return r
				},
			},
			expected: expected{stdout: "card content\n"},
		},
		"two cards": {
			given: given{
				args: []string{"board/list/card", "board/list/another-card"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil).
						Times(2)
					tr.EXPECT().
						FindList(board.ID, list.Name).
						Return(&list, nil).
						Times(2)
					tr.EXPECT().
						FindCard(list.ID, card1.Name).
						Return(&card1, nil)
					tr.EXPECT().
						FindCard(list.ID, card2.Name).
						Return(&card2, nil)
					return tr
				},
				buildRenderer: func() renderer.Renderer {
					r := renderer.NewMockRenderer(ctrl)
					r.EXPECT().
						RenderCard(card1).
						Return("card content")
					r.EXPECT().
						RenderCard(card2).
						Return("card 2 content")
					return r
				},
			},
			expected: expected{stdout: "card content\ncard 2 content\n"},
		},
		// ERRORS
		"invalid path": {
			given: given{
				args: []string{"/../.."},
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
				buildRenderer: func() renderer.Renderer {
					return nil
				},
			},
			expected: expected{
				stderr: "invalid path\n",
			},
		},
		"unknown-board": {
			given: given{
				args: []string{"unknown-board"},
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
				stderr: "no board found with name 'unknown-board'\n",
			},
		},
		"unknown-list": {
			given: given{
				args: []string{"board/unknown-list"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil)
					tr.EXPECT().
						FindList(board.ID, "unknown-list").
						Return(nil, errors.New("not found"))
					return tr
				},
				buildRenderer: func() renderer.Renderer {
					return nil
				},
			},
			expected: expected{
				stderr: "no list found with name 'unknown-list'\n",
			},
		},
		"unknown-card": {
			given: given{
				args: []string{"board/list/unknown-card"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board.Name).
						Return(&board, nil)
					tr.EXPECT().
						FindList(board.ID, list.Name).
						Return(&list, nil)
					tr.EXPECT().
						FindCard(list.ID, "unknown-card").
						Return(nil, errors.New("not found"))
					return tr
				},
				buildRenderer: func() renderer.Renderer {
					return nil
				},
			},
			expected: expected{
				stderr: "no card found with name 'unknown-card'\n",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			stdoutBuf := bytes.Buffer{}
			stderrBuf := bytes.Buffer{}
			c := cat{
				executor{
					tr:     tt.given.buildTrelloRepository(),
					r:      tt.given.buildRenderer(),
					stdout: &stdoutBuf,
					stderr: &stderrBuf,
				},
			}
			c.Execute(tt.given.args)

			actualStdout := stdoutBuf.String()
			if actualStdout != tt.expected.stdout {
				t.Errorf("expected stdout %v, actual stdout %v", tt.expected.stdout, actualStdout)
			}
			actualStderr := stderrBuf.String()
			if actualStderr != tt.expected.stderr {
				t.Errorf("expected stderr %v, actual stderr %v", tt.expected.stderr, actualStderr)
			}
		})
	}
}
