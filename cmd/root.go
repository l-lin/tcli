package cmd

import (
	"encoding/json"
	"github.com/c-bata/go-prompt"
	"github.com/c-bata/go-prompt/completer"
	"github.com/l-lin/tcli/ioc"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
)

var container *ioc.Container

// NewRootCmd creates a root command that represents the base command when called without any subcommands
func NewRootCmd(version, buildDate string) *cobra.Command {
	c := &cobra.Command{
		Use:              "tcli",
		Short:            "Start Trello interactive CLI",
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
	return c.Flags()
}

func runRootCmd(_ *cobra.Command, _ []string) {
	s := container.Session
	p := prompt.New(
		s.Executor,
		s.Completer,
		prompt.OptionTitle("Trello interactive CLI"),
		prompt.OptionLivePrefix(s.LivePrefix),
		prompt.OptionAddKeyBind(prompt.KeyBind{
			Key: prompt.ControlC,
			Fn: func(_ *prompt.Buffer) {
				os.Exit(0)
			},
		}),
		prompt.OptionCompletionWordSeparator(completer.FilePathCompletionSeparator),
	)
	p.Run()
}

func initializeIocContainer(c *cobra.Command, _ []string) {
	fp := flagParser{Command: c}
	inputs := ioc.Inputs{
		Viper: viper.New(),
		Debug: fp.GetDebug(),
		File:  fp.GetConfigFile(),
	}
	container = ioc.Boostrap(inputs)
}
