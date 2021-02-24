package executor

import (
	"github.com/l-lin/tcli/conf"
	"github.com/l-lin/tcli/renderer"
	"github.com/l-lin/tcli/trello"
	"os"
)

var Factories = []Factory{
	{
		Cmd:         "help",
		Description: "display help",
		Create: func(_ conf.Conf, tr trello.Repository, r renderer.Renderer, session *trello.Session) Executor {
			return &help{
				stdout: os.Stdout,
			}
		},
	},
	{
		Cmd:         "exit",
		Description: "exit CLI",
		Create: func(_ conf.Conf, tr trello.Repository, r renderer.Renderer, session *trello.Session) Executor {
			return &exit{}
		},
	},
	{
		Cmd:         "cd",
		Description: "change level in the hierarchy",
		Create: func(_ conf.Conf, tr trello.Repository, r renderer.Renderer, session *trello.Session) Executor {
			return &cd{executor{
				tr:      tr,
				r:       r,
				session: session,
				stdout:  os.Stdout,
				stderr:  os.Stderr,
			}}
		},
	},
	{
		Cmd:         "ls",
		Description: "list resource content",
		Create: func(_ conf.Conf, tr trello.Repository, r renderer.Renderer, session *trello.Session) Executor {
			return &ls{executor{
				tr:      tr,
				r:       r,
				session: session,
				stdout:  os.Stdout,
				stderr:  os.Stderr,
			}}
		},
	},
	{
		Cmd:         "cat",
		Description: "show resource content info",
		Create: func(_ conf.Conf, tr trello.Repository, r renderer.Renderer, session *trello.Session) Executor {
			return &cat{executor{
				tr:      tr,
				r:       r,
				session: session,
				stdout:  os.Stdout,
				stderr:  os.Stderr,
			}}
		},
	},
	{
		Cmd:         "edit",
		Description: "edit resource content",
		Create: func(conf conf.Conf, tr trello.Repository, r renderer.Renderer, session *trello.Session) Executor {
			return &edit{
				executor: executor{
					tr:      tr,
					r:       r,
					session: session,
					stdout:  os.Stdout,
					stderr:  os.Stderr,
				},
				stdin:        os.Stdin,
				editor:       NewOsEditor(conf.Editor),
				editRenderer: renderer.NewEditInPrettyYaml(),
			}
		},
	},
	{
		Cmd:         "touch",
		Description: "create new resource",
		Create: func(conf conf.Conf, tr trello.Repository, r renderer.Renderer, session *trello.Session) Executor {
			return &touch{
				executor: executor{
					tr:      tr,
					r:       r,
					session: session,
					stdout:  os.Stdout,
					stderr:  os.Stderr,
				},
			}
		},
	},
	{
		Cmd:         "rm",
		Description: "archive resource",
		Create: func(conf conf.Conf, tr trello.Repository, r renderer.Renderer, session *trello.Session) Executor {
			return &rm{
				executor: executor{
					tr:      tr,
					r:       r,
					session: session,
					stdout:  os.Stdout,
					stderr:  os.Stderr,
				},
				stdin: os.Stdin,
			}
		},
	},
	{
		Cmd:         "mv",
		Description: "move resource",
		Create: func(conf conf.Conf, tr trello.Repository, r renderer.Renderer, session *trello.Session) Executor {
			return &mv{
				executor: executor{
					tr:      tr,
					r:       r,
					session: session,
					stdout:  os.Stdout,
					stderr:  os.Stderr,
				},
			}
		},
	},
	{
		Cmd:         "cp",
		Description: "copy resource",
		Create: func(conf conf.Conf, tr trello.Repository, r renderer.Renderer, session *trello.Session) Executor {
			return &cp{
				executor: executor{
					tr:      tr,
					r:       r,
					session: session,
					stdout:  os.Stdout,
					stderr:  os.Stderr,
				},
			}
		},
	},
}

type Factory struct {
	Cmd         string
	Description string
	Create      func(conf conf.Conf, tr trello.Repository, r renderer.Renderer, session *trello.Session) Executor
}
