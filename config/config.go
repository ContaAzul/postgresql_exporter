package config

import (
	"errors"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Sql struct {
	ConnectionName   string `yaml:"connection_name,omitempty"`
	DatabaseName     string `yaml:"database_name,omitempty"`
	DatabaseUser     string `yaml:"database_user,omitempty"`
	DatabasePassword string `yaml:"database_password,omitempty"`
}

type Database struct {
	URL  string `yaml:"url,omitempty"`
	Name string `yaml:"name,omitempty"`
	Sql  Sql    `yaml:"sql,omitempty"`
}

type Config struct {
	Databases []Database `yaml:"databases,omitempty"`
}

func Parse(path string) (Config, error) {
	var cfg Config
	bts, err := ioutil.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("failed to read config file '%s': %s", path, err)
	}
	if err := yaml.Unmarshal(bts, &cfg); err != nil {
		return cfg, fmt.Errorf("failed to unmarshall config file '%s': %s", path, err)
	}
	if err := validate(cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func validate(config Config) error {
	names := make(map[string]bool)
	for _, conf := range config.Databases {
		if conf.Name == "" {
			return errors.New("failed to validate configuration. Database name cannot be empty")
		}
		if conf.URL == "" && conf.Sql.ConnectionName == "" {
			return fmt.Errorf("failed to validate configuration. URL or sql field cannot be empty in the '%s' database", conf.Name)
		}
		if names[conf.Name] {
			return fmt.Errorf("failed to validate configuration. A database named '%s' has already been declared", conf.Name)
		}
		names[conf.Name] = true
	}
	return nil
}
