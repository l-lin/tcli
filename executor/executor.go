package executor

import (
	"errors"
	"fmt"
	"github.com/l-lin/tcli/conf"
	"github.com/l-lin/tcli/renderer"
	"github.com/l-lin/tcli/trello"
	"io"
)

var errInvalidPath = errors.New("invalid path")

type Executor interface {
	Execute(args []string)
}

type executor struct {
	tr      trello.Repository
	r       renderer.Renderer
	session *trello.Session
	stdout  io.Writer
	stderr  io.Writer
}

func (e executor) getCardFromArg(arg string) (*trello.Card, error) {
	pathResolver := trello.NewPathResolver(e.session)
	p, err := pathResolver.Resolve(arg)
	if err != nil {
		return nil, err
	}

	if p.BoardName == "" || p.ListName == "" || p.CardName == "" {
		return nil, fmt.Errorf("invalid path")
	}

	var list *trello.List
	if list, err = e.getList(p.BoardName, p.ListName); err != nil {
		return nil, err
	}

	var card *trello.Card
	if card, err = e.tr.FindCard(list.ID, p.CardName); err != nil || card == nil {
		return nil, fmt.Errorf("no card found with name '%s'", p.CardName)
	}
	return card, nil
}

func (e executor) getListFromArg(arg string) (*trello.List, error) {
	pathResolver := trello.NewPathResolver(e.session)
	p, err := pathResolver.Resolve(arg)
	if err != nil {
		return nil, err
	}

	if p.BoardName == "" || p.ListName == "" {
		return nil, fmt.Errorf("invalid path")
	}

	return e.getList(p.BoardName, p.ListName)
}

func (e executor) getListAndCardNameFromArg(arg string) (*trello.List, string, error) {
	pathResolver := trello.NewPathResolver(e.session)
	p, err := pathResolver.Resolve(arg)
	if err != nil {
		return nil, "", err
	}
	if p.BoardName == "" || p.ListName == "" {
		return nil, "", fmt.Errorf("invalid path")
	}
	list, err := e.getList(p.BoardName, p.ListName)
	return list, p.CardName, err
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

func New(conf conf.Conf, cmd string, tr trello.Repository, r renderer.Renderer, session *trello.Session) Executor {
	for _, factory := range Factories {
		if factory.Cmd == cmd {
			return factory.Create(conf, tr, r, session)
		}
	}
	return nil
}

// start step builder pattern to have fluent way to process the executions
func start(tr trello.Repository) *stepExecutor {
	return &stepExecutor{tr: tr, session: &trello.Session{}}
}

type stepExecutor struct {
	err        error
	isFinished bool
	tr         trello.Repository
	session    *trello.Session
	p          trello.Path
}

func (se *stepExecutor) resolvePath(currentSession *trello.Session, arg string) *boardStepExecutor {
	pathResolver := trello.NewPathResolver(currentSession)
	se.p, se.err = pathResolver.Resolve(arg)
	return &boardStepExecutor{stepExecutor: *se}
}

// BOARD STEP EXECUTOR -------------------------------------------------------------------

type boardStepExecutor struct {
	stepExecutor
}

func (bse *boardStepExecutor) doOnEmptyBoardName(action func()) *boardStepExecutor {
	if bse.err != nil {
		return bse
	}
	if bse.p.BoardName == "" {
		action()
		bse.isFinished = true
	}
	return bse
}

func (bse *boardStepExecutor) thenFindBoard() *listStepExecutor {
	if bse.err != nil || bse.isFinished {
		return &listStepExecutor{stepExecutor: bse.stepExecutor}
	}
	var err error
	if bse.p.BoardName == "" {
		bse.err = errInvalidPath
	} else if bse.session.Board, err = bse.tr.FindBoard(bse.p.BoardName); err != nil || bse.session.Board == nil {
		bse.err = fmt.Errorf("no board found with name '%s'", bse.p.BoardName)
	}
	return &listStepExecutor{stepExecutor: bse.stepExecutor}
}

// LIST STEP EXECUTOR -------------------------------------------------------------------

type listStepExecutor struct {
	stepExecutor
}

func (lse *listStepExecutor) doOnEmptyListName(action func(session *trello.Session)) *listStepExecutor {
	if lse.err != nil || lse.isFinished {
		return lse
	}
	if lse.p.ListName == "" {
		action(lse.session)
		lse.isFinished = true
	}
	return lse
}

func (lse *listStepExecutor) thenFindList() *cardStepExecutor {
	if lse.err != nil || lse.isFinished {
		return &cardStepExecutor{stepExecutor: lse.stepExecutor}
	}
	var err error
	if lse.p.ListName == "" {
		lse.err = errInvalidPath
	} else if lse.session.List, err = lse.tr.FindList(lse.session.Board.ID, lse.p.ListName); err != nil || lse.session.List == nil {
		lse.err = fmt.Errorf("no list found with name '%s'", lse.p.ListName)
	}
	return &cardStepExecutor{stepExecutor: lse.stepExecutor}
}

// CARD STEP EXECUTOR -------------------------------------------------------------------

type cardStepExecutor struct {
	stepExecutor
}

func (cse *cardStepExecutor) doOnEmptyCardName(action func(session *trello.Session)) *cardStepExecutor {
	if cse.err != nil || cse.isFinished {
		return cse
	}
	if cse.p.CardName == "" {
		action(cse.session)
		cse.isFinished = true
	}
	return cse
}

func (cse *cardStepExecutor) doOnCardName(action func(cardName string, session *trello.Session)) *cardStepExecutor {
	if cse.err != nil || cse.isFinished {
		return cse
	}

	if cse.p.CommentID != "" {
		return cse
	}

	action(cse.p.CardName, cse.session)
	cse.isFinished = true
	return cse
}

func (cse *cardStepExecutor) thenFindCard() *commentStepExecutor {
	if cse.err != nil || cse.isFinished {
		return &commentStepExecutor{stepExecutor: cse.stepExecutor}
	}
	var err error
	if cse.p.CardName == "" {
		cse.err = errInvalidPath
	} else if cse.session.Card, err = cse.tr.FindCard(cse.session.List.ID, cse.p.CardName); err != nil || cse.session.Card == nil {
		cse.err = fmt.Errorf("no card found with name '%s'", cse.p.CardName)
	}
	return &commentStepExecutor{stepExecutor: cse.stepExecutor}
}

// COMMENT STEP EXECUTOR -------------------------------------------------------------------

type commentStepExecutor struct {
	stepExecutor
	comment *trello.Comment
}

func (cse *commentStepExecutor) doOnEmptyCommentID(action func(session *trello.Session)) *commentStepExecutor {
	if cse.err != nil || cse.isFinished {
		return cse
	}
	if cse.p.CommentID == "" {
		action(cse.session)
		cse.isFinished = true
	}
	return cse
}

func (cse *commentStepExecutor) thenFindComment() *commentStepExecutor {
	if cse.err != nil || cse.isFinished {
		return cse
	}
	var err error
	if cse.p.CommentID == "" {
		cse.err = errInvalidPath
	} else if cse.comment, err = cse.tr.FindComment(cse.session.Card.ID, cse.p.CommentID); err != nil || cse.comment == nil {
		cse.err = fmt.Errorf("no comment found with id '%s'", cse.p.CommentID)
	}
	return cse
}

func (cse *commentStepExecutor) andDoOnComment(action func(comment *trello.Comment)) error {
	if cse.err != nil || cse.isFinished {
		return cse.err
	}
	action(cse.comment)
	return cse.err
}
