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
}

func NewPathResolver(session *Session) PathResolver {
	boardName := ""
	if session.CurrentBoard != nil {
		boardName = session.CurrentBoard.Name
	}
	listName := ""
	if session.CurrentList != nil {
		listName = session.CurrentList.Name
	}
	return PathResolver{
		Path: Path{
			BoardName: boardName,
			ListName:  listName,
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
		resolvedPath = strings.Trim(filepath.Join(pr.BoardName, pr.ListName, relativePath), "/")
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
	if len(paths) > 3 {
		err = invalidPathErr
		return
	}
	// only 3 levels: boards > lists > cards
	result := make([]string, 3)
	copy(result, paths[0:])
	p.BoardName = result[0]
	p.ListName = result[1]
	p.CardName = result[2]
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
