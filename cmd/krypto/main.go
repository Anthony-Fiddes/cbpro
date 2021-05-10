package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/spf13/viper"
)

const (
	keyFieldName        = "KEY"
	keyDefault          = "YOUR ACCESS KEY"
	passphraseFieldName = "PASSPHRASE"
	passphraseDefault   = "YOUR PASSPHRASE"
	secretFieldName     = "SECRET"
	secretDefault       = "YOUR SECRET"
	configName          = "settings"
	configType          = "json"
	configPath          = "."
	configHomePath      = "$HOME/.krypto"
)

func readConfig() {
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(configHomePath)

	defaults := map[string]string{
		keyFieldName:        keyDefault,
		passphraseFieldName: passphraseDefault,
		secretFieldName:     secretDefault,
	}
	for f, d := range defaults {
		viper.SetDefault(f, d)
	}

	err := viper.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if ok {
			// Create the config file if it doesn't already exist since Viper
			// won't do that for you.
			fp := path.Join(configPath, configName+"."+configType)
			f, err := os.OpenFile(fp, os.O_CREATE, 0755)
			if err != nil {
				log.Fatal(fmt.Errorf("error while creating config file: %w", err))
			}
			f.Close()
			err = viper.WriteConfig()
			if err != nil {
				log.Fatal(fmt.Errorf("error while writing to config file: %w", err))
			}
		} else {
			log.Fatal(fmt.Errorf("error loading config file: %w", err))
		}
	}

	for f, d := range defaults {
		if viper.Get(f) == d {
			log.Fatal(errors.New("the config file has not been changed. Check settings.json"))
		}
	}
}

func main() {
	readConfig()
}
