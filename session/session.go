package session

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/l-lin/tcli/completer"
	"github.com/l-lin/tcli/executor"
	"github.com/l-lin/tcli/renderer"
	"github.com/l-lin/tcli/trello"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"strings"
)

func NewSession(tr trello.Repository, r renderer.Renderer) *Session {
	return &Session{tr: tr, r: r, output: os.Stdout, errOutput: os.Stderr}
}

type Session struct {
	tr           trello.Repository
	r            renderer.Renderer
	CurrentBoard *trello.Board
	CurrentList  *trello.List
	output       io.Writer
	errOutput    io.Writer
}

func (s *Session) Executor(in string) {
	input := strings.TrimSpace(in)
	if input == "" {
		return
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
		fmt.Fprintf(s.errOutput, "command not found: %s\n", cmd)
	}
}

func (s *Session) Completer(d prompt.Document) []prompt.Suggest {
	c := completer.New(s.tr, s.CurrentBoard, s.CurrentList)
	input := strings.TrimSpace(d.TextBeforeCursor())
	cmd, _ := getCmd(input)
	arg := getArg(input)
	return c.Complete(cmd, arg)
}

func (s *Session) LivePrefix() (string, bool) {
	builder := strings.Builder{}
	builder.WriteString("/")
	if s.CurrentBoard != nil {
		builder.WriteString(fmt.Sprintf("%s", s.CurrentBoard.Name))
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
