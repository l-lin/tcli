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
		stdout string
		stderr string
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
			expected: expected{stdout: ""},
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
			expected: expected{stdout: "board content\n"},
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
			expected: expected{stdout: "list content\n"},
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
			expected: expected{stdout: "card content\n"},
		},
		// ERRORS
		"invalid path": {
			given: given{
				arg: "/../..",
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
				stderr: "no board found with name 'unknown-board'\n",
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
				stderr: "no list found with name 'unknown-list'\n",
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
			c.Execute(tt.given.arg)

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
