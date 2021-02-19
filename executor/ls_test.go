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
		stdout string
		stderr string
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
						FindBoards().
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
			expected: expected{stdout: "boards content\n"},
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
						FindLists("board 1").
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
			expected: expected{stdout: "lists content\n"},
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
						FindCards("list 1").
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
			expected: expected{stdout: "cards content\n"},
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
			expected: expected{stdout: "card 1 content\n"},
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
		"cannot find boards": {
			given: given{
				arg: "",
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoards().
						Return(nil, errors.New("unexpected error"))
					return tr
				},
				buildRenderer: func() renderer.Renderer {
					return nil
				},
			},
			expected: expected{
				stderr: "could not fetch boards: unexpected error\n",
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
						FindLists("board 1").
						Return(nil, errors.New("unexpected error"))
					return tr
				},
				buildRenderer: func() renderer.Renderer {
					return nil
				},
			},
			expected: expected{
				stderr: "could not fetch lists for board 'board': unexpected error\n",
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
						FindCards("list 1").
						Return(nil, errors.New("unexpected error"))
					return tr
				},
				buildRenderer: func() renderer.Renderer {
					return nil
				},
			},
			expected: expected{
				stderr: "could not fetch cards for list 'list': unexpected error\n",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			stdoutBuf := bytes.Buffer{}
			stderrBuf := bytes.Buffer{}
			l := ls{
				executor{
					tr:     tt.given.buildTrelloRepository(),
					r:      tt.given.buildRenderer(),
					stdout: &stdoutBuf,
					stderr: &stderrBuf,
				},
			}
			l.Execute(tt.given.arg)

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
