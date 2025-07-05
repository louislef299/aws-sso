package config

import (
	"fmt"
	"log"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type ConfigValue struct {
	Name        string
	Description string
}

var currentConfigValues []*ConfigValue

// BindConfigValue will bind the Viper config value to the provided pflag
func BindConfigValue(name string, flag *pflag.Flag) {
	err := viper.BindPFlag(name, flag)
	if err != nil {
		log.Println("WARNING: could not bind flag", name)
		return
	}
	AddConfigValue(fmt.Sprintf("<BOUND_FLAG>%s", name), flag.Usage)
}

func AddConfigValue(name, description string) {
	currentConfigValues = append(currentConfigValues, &ConfigValue{
		Name:        name,
		Description: description,
	})
}

func GetCurrentConfigValues() []*ConfigValue {
	return currentConfigValues
}

func (c *ConfigValue) String() string {
	return fmt.Sprintf("%s:\t%s", c.Name, c.Description)
}
