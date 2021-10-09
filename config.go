package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// ConfigData struct
type ConfigData struct {
	RCCServerIP        string
	RCCServerPort      string
	RCCServerAuthToken string
}

// Config object
var Config ConfigData

func userDConfigDir() (string, error) {
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return dir, nil
	}
	// fall back to something sane on all platforms
	return os.UserConfigDir()
}

func init() {
	// Set default values for config if config does not exist
	Config = ConfigData{}

	configDir, err := userDConfigDir()
	if err != nil {
		panic(err)
	}

	configPath := fmt.Sprintf("%s/rcc/config.json", configDir)

	if _, err := os.Stat(fmt.Sprintf("%s/rcc", configDir)); err != nil {
		if err := os.MkdirAll(fmt.Sprintf("%s/rcc", configDir), 0755); err != nil {
			panic(err)
		}
	}

	err = GetConfig(configPath, &Config)
	if err != nil {
		panic(err)
	}

	if Config.RCCServerAuthToken == "" || Config.RCCServerIP == "" || Config.RCCServerPort == "" {
		fmt.Printf("Please update your config file at '%s'\n", configPath)
		os.Exit(0)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// GetConfig get a config
func GetConfig(configFileName string, configPointer interface{}) error {
	if fileExists(configFileName) { // Get existing configuration from configFileName
		b, err := ioutil.ReadFile(configFileName)
		if err != nil {
			return err
		}

		err = json.Unmarshal(b, configPointer)
		if err != nil {
			fmt.Println("Failed to unmarshal configuration file")
			return err
		}

		return nil
	}

	// If configFileName doesn't exist, create a new config file
	b, err := json.MarshalIndent(configPointer, "", " ")
	if err != nil {
		fmt.Println("Failed to marshal configuration file")
		return err
	}

	// Sevae default config into to file.
	err = ioutil.WriteFile(configFileName, b, 0755)
	if err != nil {
		fmt.Println("Failed to write configuration file")
		return err
	}

	return nil
}
