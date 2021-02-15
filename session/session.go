package session

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/l-lin/tcli/executor"
	"github.com/l-lin/tcli/renderer"
	"github.com/l-lin/tcli/trello"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

const maxCardDescriptionLength = 20

var exitCommands = []string{"quit", "q", "exit"}

func NewSession(tr trello.Repository, r renderer.Renderer) *Session {
	return &Session{tr: tr, r: r}
}

type Session struct {
	tr           trello.Repository
	r            renderer.Renderer
	CurrentBoard *trello.Board
	CurrentList  *trello.List
}

func (s *Session) Executor(input string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return
	}
	for _, exitCommand := range exitCommands {
		if input == exitCommand {
			os.Exit(0)
			return
		}
	}

	cmd, found := getCmd(input)
	if !found {
		return
	}
	arg := getArg(input)
	log.Debug().
		Str("cmd", cmd).
		Str("arg", arg).
		Msg("executing command")
	if e := executor.New(cmd, s.tr, s.r, s.CurrentBoard, s.CurrentList); e != nil {
		s.CurrentBoard, s.CurrentList = e.Execute(arg)
	} else {
		log.Error().
			Str("cmd", cmd).
			Msg("executor not found")
	}
}

func (s *Session) Completer(d prompt.Document) []prompt.Suggest {
	cmd, found := getCmd(d.TextBeforeCursor())
	if !found {
		return []prompt.Suggest{}
	}
	if s.CurrentList != nil {
		return completerAtCardsLevel(cmd, s.CurrentList.ID, s.tr)(d)
	}
	if s.CurrentBoard != nil {
		return completerAtListsLevel(cmd, s.CurrentBoard.ID, s.tr)(d)
	}
	return completerAtBoardsLevel(cmd, s.tr)(d)
}

func (s *Session) LivePrefix() (string, bool) {
	builder := strings.Builder{}
	if s.CurrentBoard != nil {
		builder.WriteString(fmt.Sprintf("/%s", s.CurrentBoard.Name))
	}
	if s.CurrentList != nil {
		builder.WriteString(fmt.Sprintf("/%s", s.CurrentList.Name))
	}
	builder.WriteString("> ")
	return builder.String(), true
}

func getCmd(s string) (string, bool) {
	args := strings.Split(s, " ")
	if len(args) == 0 {
		return "", false
	}
	if args[0] == "" {
		return "", false
	}
	return args[0], true
}

func getArg(s string) string {
	args := strings.Split(s, " ")
	if len(args) < 2 {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(s, args[0]))
}
