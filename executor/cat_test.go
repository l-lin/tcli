package executor

import (
	"bytes"
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
	var tests = map[string]struct {
		given    given
		expected string
	}{
		"no arg": {
			given: given{
				arg:                   "",
				buildTrelloRepository: func() trello.Repository { return nil },
				buildRenderer:         func() renderer.Renderer { return nil },
			},
			expected: "",
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
			expected: "board content",
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
			expected: "list content",
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
			expected: "card content",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			buf := bytes.Buffer{}
			c := cat{
				executor{
					tr:     tt.given.buildTrelloRepository(),
					r:      tt.given.buildRenderer(),
					output: &buf,
				},
			}
			c.Execute(tt.given.arg)

			actual := buf.String()
			if actual != tt.expected {
				t.Errorf("expected %v, actual %v", tt.expected, actual)
			}
		})
	}
}
