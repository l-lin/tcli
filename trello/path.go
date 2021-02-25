package trello

import (
	"errors"
	"github.com/rs/zerolog/log"
	"path"
	"path/filepath"
	"strings"
)

var invalidPathErr = errors.New("invalid path")

type Path struct {
	BoardName string
	ListName  string
	CardName  string
	CommentID string
}

func NewPathResolver(session *Session) PathResolver {
	boardName := ""
	if session.Board != nil {
		boardName = session.Board.Name
	}
	listName := ""
	if session.List != nil {
		listName = session.List.Name
	}
	cardName := ""
	if session.Card != nil {
		cardName = session.Card.Name
	}
	return PathResolver{
		Path: Path{
			BoardName: boardName,
			ListName:  listName,
			CardName:  cardName,
		},
	}
}

type PathResolver struct {
	Path
}

func (pr *PathResolver) Resolve(relativePath string) (p Path, err error) {
	if relativePath == "" {
		p = pr.Path
		pr.logResolved(p)
		return
	}
	var resolvedPath string
	if path.IsAbs(relativePath) {
		resolvedPath = strings.Trim(relativePath, "/")
	} else {
		resolvedPath = strings.Trim(filepath.Join(pr.BoardName, pr.ListName, pr.CardName, relativePath), "/")
	}
	if isInvalid(resolvedPath) {
		err = invalidPathErr
		return
	}
	if isTopLevel(resolvedPath) {
		pr.logResolved(p)
		return
	}
	paths := strings.Split(resolvedPath, "/")
	if len(paths) > 4 {
		err = invalidPathErr
		return
	}
	// only 4 levels: boards > lists > cards > comments
	result := make([]string, 4)
	copy(result, paths[0:])
	p.BoardName = result[0]
	p.ListName = result[1]
	p.CardName = result[2]
	p.CommentID = result[3]
	pr.logResolved(p)
	return
}

func (pr *PathResolver) logResolved(p Path) {
	log.Debug().
		Interface("currentPath", pr).
		Interface("resolvedPath", p).
		Msg("resolved path")
}

func isInvalid(resolvedPath string) bool {
	return strings.HasPrefix(resolvedPath, "..")
}

func isTopLevel(resolvedPath string) bool {
	return resolvedPath == "."
}
