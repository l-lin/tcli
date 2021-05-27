package executor

import (
	"fmt"
	"os/exec"
)

type clear struct {
	executor
}

func (e *clear) Execute(_ []string) {
	e.tr.Refresh()
	cmd := exec.Command("/bin/sh", "-c", "clear")
	cmd.Stdout = e.stdout
	cmd.Stderr = e.stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprint(e.stderr, err.Error())
	}
}
