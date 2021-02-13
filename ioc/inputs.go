package ioc

import "github.com/spf13/viper"

// Inputs are values provided from the user inputs (e.g. flags or arguments)
type Inputs struct {
	// mandatory fields
	Viper *viper.Viper
	Debug bool
	File  string
	// key used to identify the app
	TrelloDevKey string
	// name of the app registered in the Trello account
	TrelloAppName string
}
