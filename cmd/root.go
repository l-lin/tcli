package cmd

import (
	"encoding/json"
	"github.com/l-lin/tcli/ioc"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var container *ioc.Container

// NewRootCmd creates a root command that represents the base command when called without any subcommands
func NewRootCmd(version, buildDate string) *cobra.Command {
	c := &cobra.Command{
		Use:              "tcli",
		Short:            "Trello interactive CLI",
		Run:              runRootCmd,
		PersistentPreRun: initializeIocContainer,
	}
	initRootCmd(c, version, buildDate)
	return c
}

func initRootCmd(c *cobra.Command, version, buildDate string) *pflag.FlagSet {
	c.Version = func(version, buildDate string) string {
		res, err := json.Marshal(map[string]string{"version": version, "build_date": buildDate})
		if err != nil {
			log.Fatal().Err(err).Msg("could not marshal version json")
		}
		return string(res)
	}(version, buildDate)

	c.SetVersionTemplate(`{{printf "%s" .Version}}`)
	c.PersistentFlags().String("config", "", "config file (default will look at $PWD/.tcli.yml then at $HOME/.tcli.yml)")
	c.PersistentFlags().Bool("debug", false, "debug mode")
	c.PersistentFlags().String("trello-dev-key", "", "override Trello developer key")
	c.PersistentFlags().String("trello-app-name", "", "override Trello app name")
	return c.Flags()
}

func runRootCmd(c *cobra.Command, _ []string) {
	log.Info().Msg("Hello, world")
}

func initializeIocContainer(c *cobra.Command, _ []string) {
	fp := flagParser{Command: c}
	inputs := ioc.Inputs{
		Viper:         viper.New(),
		Debug:         fp.GetDebug(),
		File:          fp.GetConfigFile(),
		TrelloDevKey:  fp.GetTrelloDevKey(),
		TrelloAppName: fp.GetTrelloAppName(),
	}
	container = ioc.Boostrap(inputs)
}
