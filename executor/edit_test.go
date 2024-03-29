package executor

import (
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/l-lin/tcli/renderer"
	"github.com/l-lin/tcli/trello"
	"gopkg.in/yaml.v2"
	"io"
	"strings"
	"testing"
)

const yamlFileType = "yaml"

var acceptStdin = func() io.ReadCloser { return mockReadWriterCloser{strings.NewReader("y\n")} }
var refuseStdin = func() io.ReadCloser { return mockReadWriterCloser{strings.NewReader("N\n")} }

func TestEdit_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type given struct {
		args                  []string
		buildTrelloRepository func() trello.Repository
		buildEditor           func() Editor
		stdin                 io.ReadCloser
	}
	type expected struct {
		stdout string
		stderr string
	}

	board1 := trello.Board{ID: "board 1", Name: "board"}
	list1 := trello.List{ID: "list 1", Name: "list"}
	list2 := trello.List{ID: "list 2", Name: "list name 2"}
	list3 := trello.List{ID: "list 3", Name: "list name 3"}
	lists := trello.Lists{list1, list2, list3}
	card1 := trello.Card{ID: "card 1", Name: "card", Desc: "card description", Closed: true, IDBoard: board1.ID, IDList: list1.ID, Pos: float64(123)}
	createdCard1 := trello.Card{ID: "card 1", Name: "created card", Desc: "created card description", Closed: false, IDBoard: board1.ID, IDList: list1.ID, Pos: card1.Pos}
	updatedCard1 := trello.Card{ID: "card 1", Name: "updated card", Desc: "updated card description", Closed: true, IDBoard: board1.ID, IDList: list1.ID, Pos: card1.Pos}
	cte1 := trello.NewCardToEdit(card1)
	labels := trello.Labels{
		{ID: "label 1", Name: "label name 1", Color: "red"},
		{ID: "label 2", Name: "label name 2", Color: "sky"},
		{ID: "label 3", Name: "", Color: "black"},
	}
	comment := trello.Comment{
		ID: "comment",
		Data: trello.CommentData{
			Card: trello.CommentDataCard{ID: card1.ID},
			Text: "comment content",
		},
	}
	createComment := trello.CreateComment{
		IDCard: card1.ID,
		Text:   "edited comment",
	}
	updateComment := trello.UpdateComment{
		ID:     comment.ID,
		IDCard: card1.ID,
		Text:   "edited comment",
	}

	editRenderer := renderer.NewEditInYaml()
	var tests = map[string]struct {
		given    given
		expected expected
	}{
		// CARD
		"edit /board/list/card - card creation": {
			given: given{
				args: []string{"/board/list/card"},
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
						Return(nil, errors.New("not found"))
					tr.EXPECT().
						FindLists(board1.ID).
						Return(lists, nil)
					tr.EXPECT().
						FindLabels(board1.ID).
						Return(labels, nil)
					tr.EXPECT().
						CreateCard(trello.NewCreateCard(createdCard1)).
						Return(&createdCard1, nil)
					return tr
				},
				buildEditor: func() Editor {
					ctc1 := trello.NewCardToCreate(trello.Card{
						Name:   card1.Name,
						IDList: card1.IDList,
					}, []string{})
					in, _ := editRenderer.MarshalCardToCreate(ctc1, nil, nil)
					out, _ := yaml.Marshal(trello.NewCardToCreate(createdCard1, []string{}))
					e := NewMockEditor(ctrl)
					e.EXPECT().
						Edit(in, yamlFileType).
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
		"edit /board/list/card - card edition": {
			given: given{
				args: []string{"/board/list/card"},
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
						FindLists(board1.ID).
						Return(lists, nil)
					tr.EXPECT().
						FindLabels(board1.ID).
						Return(labels, nil)
					tr.EXPECT().
						UpdateCard(trello.NewUpdateCard(updatedCard1)).
						Return(&updatedCard1, nil)
					return tr
				},
				buildEditor: func() Editor {
					in, _ := editRenderer.MarshalCardToEdit(cte1, nil, nil)
					out, _ := yaml.Marshal(trello.NewCardToEdit(updatedCard1))
					e := NewMockEditor(ctrl)
					e.EXPECT().
						Edit(in, yamlFileType).
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
				args: []string{"/board/list/card"},
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
						FindLists(board1.ID).
						Return(lists, nil)
					tr.EXPECT().
						FindLabels(board1.ID).
						Return(labels, nil)
					tr.EXPECT().
						UpdateCard(trello.NewUpdateCard(updatedCard1)).
						Return(&updatedCard1, nil).
						Times(0)
					return tr
				},
				buildEditor: func() Editor {
					in, _ := editRenderer.MarshalCardToEdit(cte1, nil, nil)
					e := NewMockEditor(ctrl)
					e.EXPECT().
						Edit(in, yamlFileType).
						Return(in, nil)
					return e
				},
				stdin: refuseStdin(),
			},
			expected: expected{},
		},
		// COMMENTS
		"edit /board/list/card/comment - comment creation": {
			given: given{
				args: []string{"/board/list/card/comment"},
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
						FindComment(card1.ID, comment.ID).
						Return(nil, commentNotFoundError(comment.ID))
					tr.EXPECT().
						CreateComment(createComment).
						Return(nil, nil)
					return tr
				},
				buildEditor: func() Editor {
					in := []byte(comment.ID)
					out := []byte(createComment.Text)
					e := NewMockEditor(ctrl)
					e.EXPECT().
						Edit(in, markdownFileType).
						Return(out, nil)
					return e
				},
				stdin: acceptStdin(),
			},
			expected: expected{},
		},
		"edit /board/list/card/comment - comment edition": {
			given: given{
				args: []string{"/board/list/card/comment"},
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
						FindComment(card1.ID, comment.ID).
						Return(&comment, nil)
					tr.EXPECT().
						UpdateComment(updateComment).
						Return(&comment, nil)
					return tr
				},
				buildEditor: func() Editor {
					in := []byte(comment.Data.Text)
					out := []byte(updateComment.Text)
					e := NewMockEditor(ctrl)
					e.EXPECT().
						Edit(in, markdownFileType).
						Return(out, nil)
					return e
				},
				stdin: acceptStdin(),
			},
			expected: expected{},
		},
		// ERRORS
		"edit": {
			given: given{
				args: []string{},
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
				buildEditor: func() Editor {
					return nil
				},
			},
			expected: expected{},
		},
		"edit /../..": {
			given: given{
				args: []string{"/../.."},
				buildTrelloRepository: func() trello.Repository {
					return nil
				},
				buildEditor: func() Editor {
					return nil
				},
			},
			expected: expected{
				stderr: "invalid path\n",
			},
		},
		"edit /": {
			given: given{
				args: []string{"/"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					return tr
				},
				buildEditor: func() Editor {
					return nil
				},
			},
			expected: expected{
				stderr: "nothing to edit\n",
			},
		},
		"no board found": {
			given: given{
				args: []string{"/board"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(nil, errors.New("not found"))
					return tr
				},
				buildEditor: func() Editor {
					return nil
				},
			},
			expected: expected{
				stderr: "no board found with name 'board'\n",
			},
		},
		"no list name": {
			given: given{
				args: []string{"/board"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					return tr
				},
				buildEditor: func() Editor {
					return nil
				},
			},
			expected: expected{
				stderr: "board edition not implemented yet\n",
			},
		},
		"no list found": {
			given: given{
				args: []string{"/board/list"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindList(board1.ID, list1.Name).
						Return(nil, errors.New("not found"))
					return tr
				},
				buildEditor: func() Editor {
					return nil
				},
			},
			expected: expected{
				stderr: "no list found with name 'list'\n",
			},
		},
		"no card name": {
			given: given{
				args: []string{"/board/list"},
				buildTrelloRepository: func() trello.Repository {
					tr := trello.NewMockRepository(ctrl)
					tr.EXPECT().
						FindBoard(board1.Name).
						Return(&board1, nil)
					tr.EXPECT().
						FindList(board1.ID, list1.Name).
						Return(&list1, nil)
					return tr
				},
				buildEditor: func() Editor {
					return nil
				},
			},
			expected: expected{
				stderr: "list edition not implemented yet\n",
			},
		},
		"edit /board/list/card - error when updating card": {
			given: given{
				args: []string{"/board/list/card"},
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
						FindLists(board1.ID).
						Return(lists, nil)
					tr.EXPECT().
						FindLabels(board1.ID).
						Return(labels, nil)
					tr.EXPECT().
						UpdateCard(trello.NewUpdateCard(card1)).
						Return(nil, errors.New("unexpected error"))
					return tr
				},
				buildEditor: func() Editor {
					in, _ := editRenderer.MarshalCardToEdit(cte1, nil, nil)
					e := NewMockEditor(ctrl)
					e.EXPECT().
						Edit(in, yamlFileType).
						Return(in, nil)
					return e
				},
				stdin: acceptStdin(),
			},
			expected: expected{
				stderr: "could not edit card 'card': unexpected error\n",
			},
		},
		"edit /board/list/card/comment - error when creating comment": {
			given: given{
				args: []string{"/board/list/card/comment"},
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
						FindComment(card1.ID, comment.ID).
						Return(nil, commentNotFoundError(comment.ID))
					tr.EXPECT().
						CreateComment(createComment).
						Return(nil, errors.New("unexpected error"))
					return tr
				},
				buildEditor: func() Editor {
					in := []byte(comment.ID)
					out := []byte(createComment.Text)
					e := NewMockEditor(ctrl)
					e.EXPECT().
						Edit(in, markdownFileType).
						Return(out, nil)
					return e
				},
				stdin: acceptStdin(),
			},
			expected: expected{
				stderr: "could not create comment 'comment': unexpected error\n",
			},
		},
		"edit /board/list/card/comment - error when updating comment": {
			given: given{
				args: []string{"/board/list/card/comment"},
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
						FindComment(card1.ID, comment.ID).
						Return(&comment, nil)
					tr.EXPECT().
						UpdateComment(updateComment).
						Return(nil, errors.New("unexpected error"))
					return tr
				},
				buildEditor: func() Editor {
					in := []byte(comment.Data.Text)
					out := []byte(createComment.Text)
					e := NewMockEditor(ctrl)
					e.EXPECT().
						Edit(in, markdownFileType).
						Return(out, nil)
					return e
				},
				stdin: acceptStdin(),
			},
			expected: expected{
				stderr: "could not update comment 'comment': unexpected error\n",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			stdoutBuf := bytes.Buffer{}
			stderrBuf := bytes.Buffer{}
			e := edit{
				executor: executor{
					tr:      tt.given.buildTrelloRepository(),
					session: &trello.Session{},
					stdout:  &stdoutBuf,
					stderr:  &stderrBuf,
				},
				editor:       tt.given.buildEditor(),
				stdin:        tt.given.stdin,
				editRenderer: editRenderer,
			}
			e.Execute(tt.given.args)
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
