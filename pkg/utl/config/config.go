package config

import (
	"flag"
	"fmt"
	"log"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// checks if path has expected name format
func isExpectedConfigPath(cfgPath string) (error) {
	expectedPaths := map[string]bool{
    "./configs/conf.development.yaml": true,
    "./configs/conf.testing.yaml": true,
    "./configs/conf.staging.yaml": true,
		"./configs/conf.production.yaml": true,
	}
	if val, status := expectedPaths[cfgPath]; !(val && status) {
    return fmt.Errorf("error with path name: you must follow the syntax of './configs/conf.ENVIRONMENT.yaml'")
	}
	if _, errPath := os.Stat(cfgPath); errPath != nil {
		return fmt.Errorf("error finding the path, %s", cfgPath)
	} else {
		log.Printf("config file: %s", cfgPath)
		return nil
	}
}

// LoadConfig returns Configuration struct
func LoadConfig() (*Configuration, error) {
	cfgPath := flag.String("config", "./configs/conf.development.yaml", "Path to config file")
	flag.Parse()

	if errName := isExpectedConfigPath(*cfgPath); errName != nil {
		return nil, errName
	}
	bytes, errRead := ioutil.ReadFile(*cfgPath)
	if errRead != nil {
		return nil, fmt.Errorf("error reading config file, %s", errRead)
	}
	var cfg = new(Configuration)
	if errYaml := yaml.Unmarshal(bytes, cfg); errYaml != nil {
		return nil, fmt.Errorf("unable to decode config yaml into struct, %v", errYaml)
	}
	return cfg, nil
}

// Configuration holds data necessery for configuring application
type Configuration struct {
	Server *Server      `yaml:"server,omitempty"`
	DB     *Database    `yaml:"database,omitempty"`
	JWT    *JWT         `yaml:"jwt,omitempty"`
	App    *Application `yaml:"application,omitempty"`
}

// Database holds data necessery for database configuration
type Database struct {
	Dialect	 string `yaml:"dialect,omitempty"`
  User 		 string `yaml:"user,omitempty"`
  Password string `yaml:"password,omitempty"`
  Name 		 string `yaml:"name,omitempty"`
  Protocol string `yaml:"protocol,omitempty"`
  Host 		 string `yaml:"host,omitempty"`
  Port		 string `yaml:"port,omitempty"`
	Params 	 string `yaml:"params,omitempty"`
}

// Server holds data necessery for server configuration
type Server struct {
	Port         string `yaml:"port,omitempty"`
	Debug        bool   `yaml:"debug,omitempty"`
	ReadTimeout  int    `yaml:"read_timeout_seconds,omitempty"`
	WriteTimeout int    `yaml:"write_timeout_seconds,omitempty"`
}

// JWT holds data necessery for JWT configuration
type JWT struct {
	Secret           string `yaml:"secret,omitempty"`
	Duration         int    `yaml:"duration_minutes,omitempty"`
	RefreshDuration  int    `yaml:"refresh_duration_minutes,omitempty"`
	MaxRefresh       int    `yaml:"max_refresh_minutes,omitempty"`
	SigningAlgorithm string `yaml:"signing_algorithm,omitempty"`
}

// Application holds application configuration details
type Application struct {
	MinPasswordStr int    `yaml:"min_password_strength,omitempty"`
	SwaggerUIPath  string `yaml:"swagger_ui_path,omitempty"`
}
