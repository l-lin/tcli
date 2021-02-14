package executor

import (
	"github.com/l-lin/tcli/renderer"
	"github.com/l-lin/tcli/trello"
)

type Executor struct {
	tr trello.Repository
	r  renderer.Renderer
}
