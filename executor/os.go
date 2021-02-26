package executor

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type OS struct {
	executor
	stdin io.Reader
}

func NewOS(stdin io.Reader, stdout, stderr io.Writer) Executor {
	return &OS{
		executor: executor{
			stdout: stdout,
			stderr: stderr,
		},
		stdin: stdin,
	}
}

func (o OS) Execute(args []string) {
	s := strings.Join(args, " ")
	cmd := exec.Command("/bin/sh", "-c", s)
	cmd.Stdin = o.stdin
	cmd.Stdout = o.stdout
	cmd.Stderr = o.stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprint(o.stderr, err.Error())
	}
}
