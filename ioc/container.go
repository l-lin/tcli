package ioc

import (
	"github.com/charmbracelet/glamour"
	"github.com/l-lin/tcli/conf"
	"github.com/l-lin/tcli/renderer"
	"github.com/l-lin/tcli/session"
	"github.com/l-lin/tcli/trello"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Container IoC used to bootstrap the beans
type Container struct {
	Inputs
	*conf.Conf
	TrelloRepository trello.Repository
	Renderer         renderer.Renderer
	*session.Session
}

// Bootstrap the beans from the given user inputs
func Boostrap(inputs Inputs) *Container {
	container := &Container{
		Inputs: inputs,
	}
	container.registerConf()
	container.registerTrelloRepository()
	container.registerRenderer()
	container.registerSession()

	container.setLogLevel()
	return container
}

func (c *Container) registerTrelloRepository() {
	var tr trello.Repository
	tr = trello.NewHttpRepository(*c.Conf, c.Debug)

	if c.Inputs.NoCache {
		c.TrelloRepository = tr
	} else {
		var cacheTr trello.Repository
		cacheTr = trello.NewCacheInMemory(tr)

		c.TrelloRepository = cacheTr
	}
}

func (c *Container) registerConf() {
	var cr conf.Repository
	cr = conf.NewFileRepository(c.File, c.Viper)
	var cp conf.Provider
	cp = conf.NewProvider(cr)
	if err := cp.Init(); err != nil {
		log.Fatal().
			Stack().
			Err(err).
			Msg("could not initialize config")
	}
	c.Conf = cp.Get()
}

func (c *Container) setLogLevel() {
	if c.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().
			Str("configFile", c.Viper.ConfigFileUsed()).
			Msg("using config file")
		c.Viper.Debug()
	}
}

func (c *Container) registerRenderer() {
	var cdr renderer.Description
	var err error
	if cdr, err = glamour.NewTermRenderer(glamour.WithAutoStyle()); err != nil {
		log.Fatal().
			Stack().
			Err(err).
			Msg("could not create description renderer")
	}

	var lr renderer.Labels
	lr = renderer.TermEnvLabel{}

	var r renderer.Renderer
	r = renderer.NewInTableRenderer(lr, cdr)
	c.Renderer = r
}

func (c *Container) registerSession() {
	c.Session = session.NewSession(*c.Conf, c.TrelloRepository, c.Renderer)
}
