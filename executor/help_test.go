package executor

import (
	"bytes"
	"testing"
)

func TestHelp_Execute(t *testing.T) {
	// GIVEN
	buf := bytes.Buffer{}
	h := help{out: &buf}

	// WHEN
	h.Execute("")

	// THEN
	expected := `help    display help
exit    exit CLI
cd      change level in the hierarchy
ls      list resource content
`
	actual := buf.String()
	if actual != expected {
		t.Errorf("expected:\n%v\nactual:\n%v", expected, actual)
	}
}
