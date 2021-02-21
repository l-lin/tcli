package executor

import (
	"github.com/l-lin/tcli/trello"
	"os"
)

type exit struct{}

func (e exit) Execute(_ []string) (*trello.Board, *trello.List) {
	os.Exit(0)
	return nil, nil
}
