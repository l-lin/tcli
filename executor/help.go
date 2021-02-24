package executor

import (
	"bytes"
	"fmt"
	"github.com/cheynewallace/tabby"
	"io"
	"text/tabwriter"
)

type help struct {
	stdout io.Writer
}

func (h help) Execute(_ []string) {
	var buffer bytes.Buffer
	w := tabwriter.NewWriter(&buffer, 0, 0, 4, ' ', 0)
	t := tabby.NewCustom(w)
	for _, executorFactory := range Factories {
		t.AddLine(executorFactory.Cmd, executorFactory.Description)
	}
	t.Print()
	fmt.Fprintf(h.stdout, "%s\n", buffer.String())
}
