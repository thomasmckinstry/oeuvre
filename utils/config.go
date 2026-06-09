package utils

import (
	"gopkg.in/yaml.v3"
	"os"
)

type ConfigModel struct {
	Theme struct {
		Focus   string `yaml:"Focus"`
		Unfocus string `yaml:"Unfocus"`
	} `yaml:"Theme"`

	Database string `yaml:"Database"`

	MediaOptions  []string `yaml:"MediaOptions"`
	StatusOptions []string `yaml:"StatusOptions"`
}

var Config ConfigModel

func ReadConfig(filepath string) {
	file, err := os.ReadFile(filepath)
	CheckError("Failed to read config: ", err)

	err = yaml.Unmarshal(file, &Config)
	CheckError("Failed to unmarshal config: ", err)
	SetTheme(Config.Theme.Focus, Config.Theme.Unfocus)
}
