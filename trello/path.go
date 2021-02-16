package trello

import (
	"errors"
	"github.com/rs/zerolog/log"
	"path"
	"path/filepath"
	"strings"
)

var invalidPathErr = errors.New("invalid path")

func NewPathResolver(currentBoard *Board, currentList *List) PathResolver {
	boardName := ""
	if currentBoard != nil {
		boardName = currentBoard.Name
	}
	listName := ""
	if currentList != nil {
		listName = currentList.Name
	}
	return PathResolver{
		currentBoardName: boardName,
		currentListName:  listName,
	}
}

type PathResolver struct {
	currentBoardName string
	currentListName  string
}

func (r *PathResolver) Resolve(relativePath string) (boardName, listName, cardName string, err error) {
	if relativePath == "" {
		boardName = r.currentBoardName
		listName = r.currentListName
		r.logResolved(boardName, listName, cardName)
		return
	}
	var resolvedPath string
	if path.IsAbs(relativePath) {
		resolvedPath = strings.Trim(relativePath, "/")
	} else {
		resolvedPath = strings.Trim(filepath.Join(r.currentBoardName, r.currentListName, relativePath), "/")
	}
	if isInvalid(resolvedPath) {
		err = invalidPathErr
		return
	}
	if isTopLevel(resolvedPath) {
		r.logResolved(boardName, listName, cardName)
		return
	}
	paths := strings.Split(resolvedPath, "/")
	// only 3 levels: boards > lists > cards
	switch len(paths) {
	case 1:
		boardName = paths[0]
	case 2:
		boardName = paths[0]
		listName = paths[1]
	case 3:
		boardName = paths[0]
		listName = paths[1]
		cardName = paths[2]
	default:
		err = invalidPathErr
	}
	r.logResolved(boardName, listName, cardName)
	return
}

func (r *PathResolver) logResolved(boardName string, listName string, cardName string) {
	log.Debug().
		Str("currentBoardName", r.currentBoardName).
		Str("currentListName", r.currentListName).
		Str("boardName", boardName).
		Str("listName", listName).
		Str("cardName", cardName).
		Msg("resolved path")
}

func isInvalid(resolvedPath string) bool {
	return strings.HasPrefix(resolvedPath, "..")
}

func isTopLevel(resolvedPath string) bool {
	return resolvedPath == "."
}
