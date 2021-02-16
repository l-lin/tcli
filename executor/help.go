package executor

import (
	"bytes"
	"fmt"
	"github.com/cheynewallace/tabby"
	"github.com/l-lin/tcli/trello"
	"io"
	"text/tabwriter"
)

type help struct {
	stdout       io.Writer
	currentBoard *trello.Board
	currentList  *trello.List
}

func (h help) Execute(_ string) (*trello.Board, *trello.List) {
	var buffer bytes.Buffer
	w := tabwriter.NewWriter(&buffer, 0, 0, 4, ' ', 0)
	t := tabby.NewCustom(w)
	for _, executorFactory := range Factories {
		t.AddLine(executorFactory.Cmd, executorFactory.Description)
	}
	t.Print()
	fmt.Fprintf(h.stdout, "%s\n", buffer.String())
	return h.currentBoard, h.currentList
}
