package executor

import (
	"github.com/rs/zerolog/log"
	"os"
	"os/exec"
)

func NewOsEditor() Editor {
	// works on other platforms (i.e. Windows, MacOS)?
	return OsEditor{Command: "editor"}
}

type OsEditor struct {
	Command string
}

func (o OsEditor) Edit(in []byte) (out []byte, err error) {
	tmpFile, err := os.CreateTemp(os.TempDir(), "tcli-*.yaml")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile.Name())

	log.Debug().
		Str("tmpFile", tmpFile.Name()).
		Msg("writing content in temp file")

	// first write the content of the card in the temp file
	if err = os.WriteFile(tmpFile.Name(), in, 0644); err != nil {
		return nil, err
	}

	cmd := exec.Command(o.Command, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return nil, err
	}

	return os.ReadFile(tmpFile.Name())
}
