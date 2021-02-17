package executor

import (
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/l-lin/tcli/trello"
	"gopkg.in/yaml.v2"
	"io"
	"strings"
	"testing"
)

var acceptStdin = func() io.ReadCloser { return mockReadWriterCloser{strings.NewReader("y\n")} }
var refuseStdin = func() io.ReadCloser { return mockReadWriterCloser{strings.NewReader("N\n")} }

func TestEdit_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type given struct {
		arg                   string
		buildTrelloRepository func() trello.Repository
		buildEditor           func() Editor
		currentBoard          *trello.Board
		currentList           *trello.List
		stdin                 io.ReadCloser
	}
	type expected struct {
		stdout string
		stderr string
	}
	var tests = map[string]struct {
		given    given
		expected expected
	}{
		// CARD
		"edit /board/list/card - happy path": {
			given: given{
				arg: "/board/list/card",
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("board").
						Return(&trello.Board{ID: "board 1", Name: "board"}, nil)
					tr.EXPECT().
						FindList("board 1", "list").
						Return(&trello.List{ID: "list 1", Name: "list"}, nil)
					card := trello.Card{ID: "card 1", Name: "card", Description: "card description", Closed: true}
					tr.EXPECT().
						FindCard("list 1", "card").
						Return(&card, nil)
					updatedCard := trello.Card{ID: "card 1", Name: "updated card", Description: "updated card description", Closed: true}
					tr.EXPECT().
						UpdateCard(trello.NewUpdateCard(updatedCard)).
						Return(&updatedCard, nil)
					return tr
				},
				buildEditor: func() Editor {
					in, _ := yaml.Marshal(newCardToEdit(trello.NewUpdateCard(trello.Card{ID: "card 1", Name: "card", Description: "card description", Closed: true})))
					out, _ := yaml.Marshal(newCardToEdit(trello.NewUpdateCard(trello.Card{ID: "card 1", Name: "updated card", Description: "updated card description", Closed: true})))
					e := NewMockEditor(ctrl)
					e.EXPECT().
						Edit(in).
						Return(out, nil)
					return e
				},
				stdin: acceptStdin(),
			},
			expected: expected{
				stdout: "",
				stderr: "",
			},
		},
		"edit /board/list/card - user refused to update card": {
			given: given{
				arg: "/board/list/card",
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("board").
						Return(&trello.Board{ID: "board 1", Name: "board"}, nil)
					tr.EXPECT().
						FindList("board 1", "list").
						Return(&trello.List{ID: "list 1", Name: "list"}, nil)
					card := trello.Card{ID: "card 1", Name: "card", Description: "card description", Closed: true}
					tr.EXPECT().
						FindCard("list 1", "card").
						Return(&card, nil)
					tr.EXPECT().
						UpdateCard(trello.NewUpdateCard(card)).
						Times(0)
					return tr
				},
				buildEditor: func() Editor {
					in, _ := yaml.Marshal(newCardToEdit(trello.NewUpdateCard(trello.Card{ID: "card 1", Name: "card", Description: "card description", Closed: true})))
					e := NewMockEditor(ctrl)
					e.EXPECT().
						Edit(in).
						Return(in, nil)
					return e
				},
				stdin: refuseStdin(),
			},
			expected: expected{
				stdout: "card 'card 1' not updated\n",
				stderr: "",
			},
		},
		"edit /board/list/card - error when updating card": {
			given: given{
				arg: "/board/list/card",
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard("board").
						Return(&trello.Board{ID: "board 1", Name: "board"}, nil)
					tr.EXPECT().
						FindList("board 1", "list").
						Return(&trello.List{ID: "list 1", Name: "list"}, nil)
					card := trello.Card{ID: "card 1", Name: "card", Description: "card description", Closed: true}
					tr.EXPECT().
						FindCard("list 1", "card").
						Return(&card, nil)
					tr.EXPECT().
						UpdateCard(trello.NewUpdateCard(card)).
						Return(nil, errors.New("unexpected error"))
					return tr
				},
				buildEditor: func() Editor {
					in, _ := yaml.Marshal(newCardToEdit(trello.NewUpdateCard(trello.Card{ID: "card 1", Name: "card", Description: "card description", Closed: true})))
					e := NewMockEditor(ctrl)
					e.EXPECT().
						Edit(in).
						Return(in, nil)
					return e
				},
				stdin: acceptStdin(),
			},
			expected: expected{
				stdout: "",
				stderr: "could not update card 'card': unexpected error\n",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			stdoutBuf := bytes.Buffer{}
			stderrBuf := bytes.Buffer{}
			e := edit{
				executor: executor{
					tr:           tt.given.buildTrelloRepository(),
					currentBoard: tt.given.currentBoard,
					currentList:  tt.given.currentList,
					stdout:       &stdoutBuf,
					stderr:       &stderrBuf,
				},
				editor: tt.given.buildEditor(),
				stdin:  tt.given.stdin,
			}
			e.Execute(tt.given.arg)
			actualStderr := stderrBuf.String()
			if actualStderr != tt.expected.stderr {
				t.Errorf("expected:\n%v\nactual:\n%v", tt.expected.stderr, actualStderr)
			}
		})
	}
}

type mockReadWriterCloser struct {
	io.Reader
}

func (m mockReadWriterCloser) Close() error {
	return nil
}
