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
		args                  []string
		buildTrelloRepository func() trello.Repository
		buildRenderer         func() renderer.Renderer
	}
	type expected struct {
		stdout string
		stderr string
	}
	board1 := trello.Board{ID: "board 1", Name: "board"}
	board2 := trello.Board{ID: "board 2", Name: "another board"}
	boards := trello.Boards{board1, board2}
	list1 := trello.List{ID: "list 1", Name: "list"}
	list2 := trello.List{ID: "list 2", Name: "another list"}
	list3 := trello.List{ID: "list 3", Name: "list 3"}
	lists1 := trello.Lists{list1, list2}
	lists2 := trello.Lists{list3}
	card1 := trello.Card{ID: "card 1", Name: "card"}
	card2 := trello.Card{ID: "card 2", Name: "another card"}
	cards := trello.Cards{card1, card2}
	comment1 := trello.Comment{ID: "comment"}
	comment2 := trello.Comment{ID: "another comment"}
	comments := trello.Comments{comment1, comment2}

	var tests = map[string]struct {
		given    given
		expected expected
	}{
		"no arg": {
			given: given{
				args: []string{},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoards().
						Return(boards, nil)
					return tr
				},
				buildRenderer: func() renderer.Renderer {
					r := renderer.NewMockRenderer(ctrl)
					r.EXPECT().
						RenderBoards(boards).
						Return("boards content")
					return r
				},
			},
			expected: expected{stdout: "boards content\n"},
		},
		"empty string as 1st arg": {
			given: given{
				args: []string{""},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoards().
						Return(boards, nil)
					return tr
				},
				buildRenderer: func() renderer.Renderer {
					r := renderer.NewMockRenderer(ctrl)
					r.EXPECT().
						RenderBoards(boards).
						Return("boards content")
					return r
				},
			},
			expected: expected{stdout: "boards content\n"},
		},
		"show lists": {
			given: given{
				args: []string{board1.Name},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindLists(board1.ID).
						Return(lists1, nil)
					return tr
				},
				buildRenderer: func() renderer.Renderer {
					r := renderer.NewMockRenderer(ctrl)
					r.EXPECT().
						RenderLists(lists1).
						Return("lists content")
					return r
				},
			},
			expected: expected{stdout: "lists content\n"},
		},
		"show cards": {
			given: given{
				args: []string{"board/list"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindList(board1.ID, list1.Name).
						Return(&list1, nil)
					tr.EXPECT().
						FindCards(list1.ID).
						Return(cards, nil)
					return tr
				},
				buildRenderer: func() renderer.Renderer {
					r := renderer.NewMockRenderer(ctrl)
					r.EXPECT().
						RenderCards(cards).
						Return("cards content")
					return r
				},
			},
			expected: expected{stdout: "cards content\n"},
		},
		"show 2 lists": {
			given: given{
				args: []string{board1.Name, board2.Name},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindBoard(board2.Name).
						Return(&board2, nil)
					tr.EXPECT().
						FindLists(board1.ID).
						Return(lists1, nil)
					tr.EXPECT().
						FindLists(board2.ID).
						Return(lists2, nil)
					return tr
				},
				buildRenderer: func() renderer.Renderer {
					r := renderer.NewMockRenderer(ctrl)
					r.EXPECT().
						RenderLists(lists1).
						Return("lists 1 content")
					r.EXPECT().
						RenderLists(lists2).
						Return("lists 2 content")
					return r
				},
			},
			expected: expected{stdout: "lists 1 content\nlists 2 content\n"},
		},
		"show comments": {
			given: given{
				args: []string{"board/list/card"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindList(board1.ID, list1.Name).
						Return(&list1, nil)
					tr.EXPECT().
						FindCard(list1.ID, card1.Name).
						Return(&card1, nil)
					tr.EXPECT().
						FindComments(card1.ID).
						Return(comments, nil)
					return tr
				},
				buildRenderer: func() renderer.Renderer {
					r := renderer.NewMockRenderer(ctrl)
					r.EXPECT().
						RenderComments(comments).
						Return("comments content")
					return r
				},
			},
			expected: expected{stdout: "comments content\n"},
		},
		"show single comment": {
			given: given{
				args: []string{"board/list/card/comment"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindList(board1.ID, list1.Name).
						Return(&list1, nil)
					tr.EXPECT().
						FindCard(list1.ID, card1.Name).
						Return(&card1, nil)
					tr.EXPECT().
						FindComment(card1.ID, comment1.ID).
						Return(&comment1, nil)
					return tr
				},
				buildRenderer: func() renderer.Renderer {
					r := renderer.NewMockRenderer(ctrl)
					r.EXPECT().
						RenderComment(comment1).
						Return("comment 1 content")
					return r
				},
			},
			expected: expected{stdout: "comment 1 content\n"},
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
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindList(board1.ID, "unknown-list").
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
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindList(board1.ID, list1.Name).
						Return(&list1, nil)
					tr.EXPECT().
						FindCard(list1.ID, "unknown-card").
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
				args: []string{""},
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
				args: []string{board1.Name},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindLists(board1.ID).
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
				args: []string{"board/list"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindList(board1.ID, list1.Name).
						Return(&list1, nil)
					tr.EXPECT().
						FindCards(list1.ID).
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
					tr:      tt.given.buildTrelloRepository(),
					r:       tt.given.buildRenderer(),
					session: &trello.Session{},
					stdout:  &stdoutBuf,
					stderr:  &stderrBuf,
				},
			}
			l.Execute(tt.given.args)

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
