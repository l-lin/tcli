package executor

import (
	"fmt"
	"github.com/l-lin/tcli/conf"
	"github.com/l-lin/tcli/renderer"
	"github.com/l-lin/tcli/trello"
	"io"
)

type Executor interface {
	Execute(args []string) (*trello.Board, *trello.List)
}

type executor struct {
	tr           trello.Repository
	r            renderer.Renderer
	currentBoard *trello.Board
	currentList  *trello.List
	stdout       io.Writer
	stderr       io.Writer
}

func (e executor) getCardFromArg(arg string) (*trello.Card, error) {
	pathResolver := trello.NewPathResolver(e.currentBoard, e.currentList)
	boardName, listName, cardName, err := pathResolver.Resolve(arg)
	if err != nil {
		return nil, err
	}

	if boardName == "" || listName == "" || cardName == "" {
		return nil, fmt.Errorf("invalid path")
	}

	var list *trello.List
	if list, err = e.getList(boardName, listName); err != nil {
		return nil, err
	}

	var card *trello.Card
	if card, err = e.tr.FindCard(list.ID, cardName); err != nil || card == nil {
		return nil, fmt.Errorf("no card found with name '%s'", cardName)
	}
	return card, nil
}

func (e executor) getListFromArg(arg string) (*trello.List, error) {
	pathResolver := trello.NewPathResolver(e.currentBoard, e.currentList)
	boardName, listName, _, err := pathResolver.Resolve(arg)
	if err != nil {
		return nil, err
	}

	if boardName == "" || listName == "" {
		return nil, fmt.Errorf("invalid path")
	}

	return e.getList(boardName, listName)
}

func (e executor) getListAndCardNameFromArg(arg string) (*trello.List, string, error) {
	pathResolver := trello.NewPathResolver(e.currentBoard, e.currentList)
	boardName, listName, cardName, err := pathResolver.Resolve(arg)
	if err != nil {
		return nil, "", err
	}
	if boardName == "" || listName == "" {
		return nil, "", fmt.Errorf("invalid path")
	}
	list, err := e.getList(boardName, listName)
	return list, cardName, err
}

func (e executor) getList(boardName, listName string) (list *trello.List, err error) {
	var board *trello.Board
	if board, err = e.tr.FindBoard(boardName); err != nil || board == nil {
		return nil, fmt.Errorf("no board found with name '%s'", boardName)
	}
	if list, err = e.tr.FindList(board.ID, listName); err != nil || list == nil {
		return nil, fmt.Errorf("no list found with name '%s'", listName)
	}
	return list, nil
}

func New(conf conf.Conf, cmd string, tr trello.Repository, r renderer.Renderer, currentBoard *trello.Board, currentList *trello.List) Executor {
	for _, factory := range Factories {
		if factory.Cmd == cmd {
			return factory.Create(conf, tr, r, currentBoard, currentList)
		}
	}
	return nil
}
