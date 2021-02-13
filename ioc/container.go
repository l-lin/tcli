package ioc

import (
	"github.com/l-lin/tcli/conf"
	"github.com/l-lin/tcli/trello"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Container IoC used to bootstrap the beans
type Container struct {
	Inputs
	*conf.Conf
	TrelloRepository trello.Repository
}

// Bootstrap the beans from the given user inputs
func Boostrap(inputs Inputs) *Container {
	container := &Container{
		Inputs: inputs,
	}
	container.registerConf()
	container.registerTrelloRepository()

	container.setLogLevel()
	return container
}

func (c *Container) registerTrelloRepository() {
	var tr trello.Repository
	tr = trello.NewHttpRepository(*c.Conf, c.Debug)
	c.TrelloRepository = tr
}

func (c *Container) registerConf() {
	var cr conf.Repository
	cr = conf.NewFileRepository(c.File, c.Viper)
	var cp conf.Provider
	cp = conf.NewProvider(cr, c.Inputs.TrelloDevKey, c.Inputs.TrelloAppName)
	if err := cp.Init(); err != nil {
		log.Fatal().
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
