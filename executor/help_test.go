package executor

import (
	"bytes"
	"testing"
)

func TestHelp_Execute(t *testing.T) {
	// GIVEN
	buf := bytes.Buffer{}
	h := help{stdout: &buf}

	// WHEN
	h.Execute(nil)

	// THEN
	expected := `help     display help
exit     exit CLI
cd       change level in the hierarchy
ls       list resource content
cat      show resource content info
edit     edit resource content
touch    create new resource
rm       archive resource
mv       move resource

`
	actual := buf.String()
	if actual != expected {
		t.Errorf("expected:\n%v\nactual:\n%v", expected, actual)
	}
}
