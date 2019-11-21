package config

import (
	"io/ioutil"

	"github.com/apex/log"
	"gopkg.in/yaml.v2"
)

type Database struct {
	URL  string `yaml:"url,omitempty"`
	Name string `yaml:"name,omitempty"`
}

type Config struct {
	Databases []Database `yaml:"databases,omitempty"`
}

func Parse(path string) Config {
	var cfg Config
	bts, err := ioutil.ReadFile(path)
	if err != nil {
		log.WithError(err).Fatalf("failed to read config file: %s", path)
	}
	if err := yaml.Unmarshal(bts, &cfg); err != nil {
		log.WithError(err).Fatalf("failed to unmarshall config file: %s", path)
	}
	return validate(cfg)
}

func validate(config Config) Config {
	names := make(map[string]bool)
	for _, conf := range config.Databases {
		if names[conf.Name] {
			log.Fatalf("failed to validate configuration. A database named '%s' has already been declared'", conf.Name)
		}
		names[conf.Name] = true
	}
	return config
}
