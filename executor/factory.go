package executor

import (
	"github.com/l-lin/tcli/conf"
	"github.com/l-lin/tcli/renderer"
	"github.com/l-lin/tcli/trello"
	"io"
	"os"
)

var Factories = []Factory{
	{
		Cmd:         "help",
		Description: "display help",
		Create: func(_ conf.Conf, tr trello.Repository, r renderer.Renderer, session *trello.Session, stdout, stderr io.Writer) Executor {
			return &help{
				stdout: stdout,
			}
		},
	},
	{
		Cmd:         "exit",
		Description: "exit CLI",
		Create: func(_ conf.Conf, tr trello.Repository, r renderer.Renderer, session *trello.Session, stdout, stderr io.Writer) Executor {
			return &exit{}
		},
	},
	{
		Cmd:         "clear",
		Description: "clear the terminal screen & cache",
		Create: func(_ conf.Conf, tr trello.Repository, r renderer.Renderer, session *trello.Session, stdout, stderr io.Writer) Executor {
			return &clear{
				executor{
					tr:      tr,
					r:       r,
					session: session,
					stdout:  stdout,
					stderr:  stderr,
				},
			}
		},
	},
	{
		Cmd:         "cd",
		Description: "change level in the hierarchy",
		Create: func(_ conf.Conf, tr trello.Repository, r renderer.Renderer, session *trello.Session, stdout, stderr io.Writer) Executor {
			return &cd{executor{
				tr:      tr,
				r:       r,
				session: session,
				stdout:  stdout,
				stderr:  stderr,
			}}
		},
	},
	{
		Cmd:         "ls",
		Description: "list resource content",
		Create: func(_ conf.Conf, tr trello.Repository, r renderer.Renderer, session *trello.Session, stdout, stderr io.Writer) Executor {
			return &ls{executor{
				tr:      tr,
				r:       r,
				session: session,
				stdout:  stdout,
				stderr:  stderr,
			}}
		},
	},
	{
		Cmd:         "cat",
		Description: "show resource content info",
		Create: func(_ conf.Conf, tr trello.Repository, r renderer.Renderer, session *trello.Session, stdout, stderr io.Writer) Executor {
			return &cat{executor{
				tr:      tr,
				r:       r,
				session: session,
				stdout:  stdout,
				stderr:  stderr,
			}}
		},
	},
	{
		Cmd:         "edit",
		Description: "edit resource content",
		Create: func(conf conf.Conf, tr trello.Repository, r renderer.Renderer, session *trello.Session, stdout, stderr io.Writer) Executor {
			return &edit{
				executor: executor{
					tr:      tr,
					r:       r,
					session: session,
					stdout:  stdout,
					stderr:  stderr,
				},
				stdin:        os.Stdin,
				editor:       NewOsEditor(conf.Editor),
				editRenderer: renderer.NewEdit(conf.Format),
			}
		},
	},
	{
		Cmd:         "touch",
		Description: "create new resource",
		Create: func(conf conf.Conf, tr trello.Repository, r renderer.Renderer, session *trello.Session, stdout, stderr io.Writer) Executor {
			return &touch{
				executor: executor{
					tr:      tr,
					r:       r,
					session: session,
					stdout:  stdout,
					stderr:  stderr,
				},
			}
		},
	},
	{
		Cmd:         "rm",
		Description: "archive resource",
		Create: func(conf conf.Conf, tr trello.Repository, r renderer.Renderer, session *trello.Session, stdout, stderr io.Writer) Executor {
			return &rm{
				executor: executor{
					tr:      tr,
					r:       r,
					session: session,
					stdout:  stdout,
					stderr:  stderr,
				},
				stdin: os.Stdin,
			}
		},
	},
	{
		Cmd:         "mv",
		Description: "move resource",
		Create: func(conf conf.Conf, tr trello.Repository, r renderer.Renderer, session *trello.Session, stdout, stderr io.Writer) Executor {
			return &mv{
				executor: executor{
					tr:      tr,
					r:       r,
					session: session,
					stdout:  stdout,
					stderr:  stderr,
				},
			}
		},
	},
	{
		Cmd:         "cp",
		Description: "copy resource",
		Create: func(conf conf.Conf, tr trello.Repository, r renderer.Renderer, session *trello.Session, stdout, stderr io.Writer) Executor {
			return &cp{
				executor: executor{
					tr:      tr,
					r:       r,
					session: session,
					stdout:  stdout,
					stderr:  stderr,
				},
			}
		},
	},
}

type Factory struct {
	Cmd         string
	Description string
	Create      func(conf conf.Conf, tr trello.Repository, r renderer.Renderer, session *trello.Session, stdout, stderr io.Writer) Executor
}
