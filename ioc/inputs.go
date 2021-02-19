package ioc

import "github.com/spf13/viper"

// Inputs are values provided from the user inputs (e.g. flags or arguments)
type Inputs struct {
	Viper   *viper.Viper
	Debug   bool
	File    string
	NoCache bool
}
