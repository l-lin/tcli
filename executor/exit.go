package executor

import (
	"os"
)

type exit struct{}

func (e exit) Execute(_ []string) {
	os.Exit(0)
}
