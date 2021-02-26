package executor

import (
	"errors"
	"fmt"
	"github.com/l-lin/tcli/conf"
	"github.com/l-lin/tcli/renderer"
	"github.com/l-lin/tcli/trello"
	"io"
)

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

func New(conf conf.Conf, cmd string, tr trello.Repository, r renderer.Renderer, session *trello.Session, stdout, stderr io.Writer) Executor {
	for _, factory := range Factories {
		if factory.Cmd == cmd {
			return factory.Create(conf, tr, r, session, stdout, stderr)
		}
	}
	return nil
}

// ERRORS -------------------------------------------------------------------

var invalidPathError = errors.New("invalid path")

type boardNotFoundError string

func (b boardNotFoundError) Error() string {
	return fmt.Sprintf("no board found with name '%s'", string(b))
}

type listNotFoundError string

func (l listNotFoundError) Error() string {
	return fmt.Sprintf("no list found with name '%s'", string(l))
}

type cardNotFoundError string

func (c cardNotFoundError) Error() string {
	return fmt.Sprintf("no card found with name '%s'", string(c))
}

type commentNotFoundError string

func (c commentNotFoundError) Error() string {
	return fmt.Sprintf("no comment found with id '%s'", string(c))
}

// STEP EXECUTORS -------------------------------------------------------------------

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

func (se *stepExecutor) resolvePath(currentSession *trello.Session, arg string) *stepExecutor {
	pathResolver := trello.NewPathResolver(currentSession)
	se.p, se.err = pathResolver.Resolve(arg)
	return se
}

func (se *stepExecutor) then() *boardStepExecutor {
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

func (bse *boardStepExecutor) findBoard() *boardStepExecutor {
	if bse.err != nil || bse.isFinished {
		return bse
	}
	var err error
	if bse.p.BoardName == "" {
		bse.err = invalidPathError
	} else if bse.session.Board, err = bse.tr.FindBoard(bse.p.BoardName); err != nil || bse.session.Board == nil {
		bse.err = boardNotFoundError(bse.p.BoardName)
	}
	return bse
}

func (bse *boardStepExecutor) doOnBoard(action func(*trello.Board)) *boardStepExecutor {
	if bse.err != nil || bse.isFinished || bse.p.ListName != "" {
		return bse
	}
	if bse.session.Board != nil {
		action(bse.session.Board)
		bse.isFinished = true
	}
	return bse
}

func (bse *boardStepExecutor) then() *listStepExecutor {
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

func (lse *listStepExecutor) findList() *listStepExecutor {
	if lse.err != nil || lse.isFinished {
		return lse
	}
	var err error
	if lse.p.ListName == "" {
		lse.err = invalidPathError
	} else if lse.session.List, err = lse.tr.FindList(lse.session.Board.ID, lse.p.ListName); err != nil || lse.session.List == nil {
		lse.err = listNotFoundError(lse.p.ListName)
	}
	return lse
}

func (lse *listStepExecutor) doOnList(action func(*trello.List)) *listStepExecutor {
	if lse.isFinished || lse.p.CardName != "" {
		return lse
	}
	if lse.session.List != nil {
		action(lse.session.List)
		lse.isFinished = true
	}
	return lse
}

func (lse *listStepExecutor) then() *cardStepExecutor {
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

func (cse *cardStepExecutor) doOnCard(action func(*trello.Card)) *cardStepExecutor {
	if cse.err != nil || cse.isFinished || cse.p.CommentID != "" {
		return cse
	}

	if cse.session.Card != nil {
		action(cse.session.Card)
		cse.isFinished = true
	}
	return cse
}

func (cse *cardStepExecutor) doOnCardName(action func(cardName string, session *trello.Session)) *cardStepExecutor {
	if cse.isFinished || cse.p.CommentID != "" {
		return cse
	}

	var cErr cardNotFoundError
	if cse.err != nil {
		if !errors.As(cse.err, &cErr) {
			return cse
		}
		// reset error as we are already handling this case afterward
		cse.err = nil
	}
	if cse.p.CardName != "" {
		action(cse.p.CardName, cse.session)
		cse.isFinished = true
	}
	return cse
}

func (cse *cardStepExecutor) findCard() *cardStepExecutor {
	if cse.err != nil || cse.isFinished {
		return cse
	}
	var err error
	if cse.p.CardName == "" {
		cse.err = invalidPathError
	} else if cse.session.Card, err = cse.tr.FindCard(cse.session.List.ID, cse.p.CardName); err != nil || cse.session.Card == nil {
		cse.err = cardNotFoundError(cse.p.CardName)
	}
	return cse
}

func (cse *cardStepExecutor) then() *commentStepExecutor {
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

func (cse *commentStepExecutor) findComment() *commentStepExecutor {
	if cse.err != nil || cse.isFinished {
		return cse
	}
	var err error
	if cse.p.CommentID == "" {
		cse.err = invalidPathError
	} else if cse.comment, err = cse.tr.FindComment(cse.session.Card.ID, cse.p.CommentID); err != nil || cse.comment == nil {
		cse.err = commentNotFoundError(cse.p.CommentID)
	}
	return cse
}

func (cse *commentStepExecutor) doOnComment(action func(comment *trello.Comment)) *commentStepExecutor {
	if cse.err != nil || cse.isFinished {
		return cse
	}
	if cse.comment != nil {
		action(cse.comment)
		cse.isFinished = true
	}
	return cse
}

func (cse *commentStepExecutor) doOnCommentText(action func(commentText string, session *trello.Session)) *commentStepExecutor {
	if cse.isFinished {
		return cse
	}

	var e commentNotFoundError
	if cse.err != nil {
		if !errors.As(cse.err, &e) {
			return cse
		}
		// reset error as we are already handling this case afterward
		cse.err = nil
	}
	if cse.p.CommentID != "" {
		// commentID contains either the comment ID or the text of the comment to create
		action(cse.p.CommentID, cse.session)
		cse.isFinished = true
	}
	return cse
}
