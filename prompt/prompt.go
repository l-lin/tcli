package prompt

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/l-lin/tcli/completer"
	"github.com/l-lin/tcli/conf"
	"github.com/l-lin/tcli/executor"
	"github.com/l-lin/tcli/renderer"
	"github.com/l-lin/tcli/trello"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"strings"
)

func NewPrompt(conf conf.Conf, tr trello.Repository, r renderer.Renderer) *Prompt {
	return &Prompt{
		conf:    conf,
		tr:      tr,
		r:       r,
		Session: &trello.Session{},
		stdout:  os.Stdout,
		stderr:  os.Stderr,
	}
}

// Prompt of the terminal to navigate seamlessly in interactive mode
type Prompt struct {
	conf    conf.Conf
	tr      trello.Repository
	r       renderer.Renderer
	Session *trello.Session
	stdout  io.Writer
	stderr  io.Writer
}

func (p *Prompt) Executor(in string) {
	input := strings.TrimSpace(in)
	if input == "" {
		return
	}

	cmd, found := getCmd(input)
	if !found {
		return
	}
	args, err := getArgs(input)
	if err != nil {
		fmt.Fprintf(p.stderr, err.Error())
		return
	}
	log.Debug().
		Str("cmd", cmd).
		Strs("args", args).
		Msg("executing command")
	if e := executor.New(p.conf, cmd, p.tr, p.r, p.Session); e != nil {
		e.Execute(args)
	} else {
		fmt.Fprintf(p.stderr, "command not found: %s\n", cmd)
	}
}

func (p *Prompt) Completer(d prompt.Document) []prompt.Suggest {
	c := completer.New(p.tr, p.Session)
	input := d.TextBeforeCursor()
	cmd, _ := getCmd(input)
	args, err := getArgs(input)
	if err != nil {
		return []prompt.Suggest{}
	}
	return c.Complete(cmd, args)
}

func (p *Prompt) LivePrefix() (string, bool) {
	builder := strings.Builder{}
	builder.WriteString("/")
	if p.Session.CurrentBoard != nil {
		builder.WriteString(fmt.Sprintf("%s", p.Session.CurrentBoard.Name))
	}
	if p.Session.CurrentList != nil {
		builder.WriteString(fmt.Sprintf("/%s", p.Session.CurrentList.Name))
	}
	builder.WriteString("> ")
	return builder.String(), true
}

func getCmd(s string) (string, bool) {
	args := strings.Split(s, " ")
	if args[0] == "" {
		return "", false
	}
	return args[0], true
}

// getArgs from the input
// shamelessly taken from https://github.com/chriswalz/bit/blob/f9bb2b246db444bb3f9f6d0d3656090d34a1905a/cmd/util.go#L508-L571
func getArgs(input string) ([]string, error) {
	var args []string
	state := "start"
	current := ""
	quote := "\""
	escapeNext := true
	runes := []rune(input)
	for i := 0; i < len(runes); i++ {
		c := runes[i]

		if state == "quotes" {
			if string(c) != quote {
				current += string(c)
			} else {
				args = append(args, current)
				current = ""
				state = "start"
			}
			continue
		}

		if escapeNext {
			current += string(c)
			escapeNext = false
			continue
		}

		if c == '\\' {
			escapeNext = true
			continue
		}

		if c == '"' || c == '\'' {
			state = "quotes"
			quote = string(c)
			continue
		}

		if state == "arg" {
			if c == ' ' || c == '\t' {
				args = append(args, current)
				current = ""
				state = "start"
			} else {
				current += string(c)
			}
			continue
		}

		if c != ' ' && c != '\t' {
			state = "arg"
			current += string(c)
		}
	}

	if state == "quotes" {
		return []string{}, fmt.Errorf("unclosed quote in command line '%s'", input)
	}

	if current != "" {
		args = append(args, current)
	}

	if args == nil || len(args) < 2 {
		return []string{}, nil
	}
	// ensure the arguments after the first one restarts the completion
	lastChar := input[len(input)-1]
	if lastChar == ' ' || lastChar == '\t' {
		args = append(args, "")
	}
	return args[1:], nil
}
