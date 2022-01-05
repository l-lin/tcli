package conf

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/rs/zerolog/log"
	"io"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
)

const (
	defaultTrelloApiKey  = "85e0a48198cd720949ad9d829c464b2e"
	defaultTrelloAppName = "Trello CLI"
)

// creator is used as a builder to create a new Conf by asking the user the needed information
type creator struct {
	*Conf
	AppName string
	Err     error
	stdin   io.ReadCloser
	stdout  io.WriteCloser
}

func (creator *creator) askTrelloAccessToken() *creator {
	if creator.Err != nil {
		return creator
	}
	v := url.Values{}
	v.Set("key", creator.ApiKey)
	v.Set("name", creator.AppName)
	v.Set("response_type", "token")
	v.Set("expires", "never")
	v.Set("scope", "read,write")
	u := fmt.Sprintf("%s/authorize?%v", creator.BaseURL, v.Encode())
	log.Info().
		Str("url", u).
		Msg("Please copy the API access token from Trello")
	openBrowser(u)
	prompt := promptui.Prompt{
		Label:    "Paste the API token here",
		Validate: validateNotEmpty,
		Stdin:    creator.stdin,
		Stdout:   creator.stdout,
	}
	creator.AccessToken, creator.Err = prompt.Run()
	return creator
}

func (creator *creator) askTrelloApiKey() *creator {
	prompt := promptui.Prompt{
		Label:    "Trello API key",
		Default:  defaultTrelloApiKey,
		Validate: validateNotEmpty,
		Stdin:    creator.stdin,
		Stdout:   creator.stdout,
	}
	creator.ApiKey, creator.Err = prompt.Run()
	return creator
}

func (creator *creator) askTrelloAppName() *creator {
	prompt := promptui.Prompt{
		Label:    "Trello app name",
		Default:  defaultTrelloAppName,
		Validate: validateNotEmpty,
		Stdin:    creator.stdin,
		Stdout:   creator.stdout,
	}
	creator.AppName, creator.Err = prompt.Run()
	return creator
}

func (creator *creator) askFormat() *creator {
	prompt := promptui.Select{
		Label:  "Format to use when editing a card (yaml or toml)",
		Items:  allFormats,
		Stdin:  creator.stdin,
		Stdout: creator.stdout,
	}
	_, creator.Format, creator.Err = prompt.Run()
	return creator
}

func (creator creator) create() (*Conf, error) {
	return creator.Conf, creator.Err
}

func validateNotEmpty(s string) error {
	if strings.Trim(s, " ") == "" || strings.Trim(s, "\t") == "" {
		return errors.New("cannot be empty")
	}
	return nil
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("could not open browser")
	}
}
