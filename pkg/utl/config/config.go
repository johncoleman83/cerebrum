// Package config is used for loading the environmental configurations
package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/johncoleman83/cerebrum/pkg/utl/support"

	yaml "gopkg.in/yaml.v2"
)

// Configuration holds data necessery for configuring application
type Configuration struct {
	Server *Server      `yaml:"server,omitempty"`
	DB     *Database    `yaml:"database,omitempty"`
	JWT    *JWT         `yaml:"jwt,omitempty"`
	App    *Application `yaml:"application,omitempty"`
}

// Database holds data necessery for database configuration
type Database struct {
	Dialect  string `yaml:"dialect,omitempty"`
	User     string `yaml:"user,omitempty"`
	Password string `yaml:"password,omitempty"`
	Name     string `yaml:"name,omitempty"`
	Protocol string `yaml:"protocol,omitempty"`
	Host     string `yaml:"host,omitempty"`
	Port     string `yaml:"port,omitempty"`
	Settings string `yaml:"settings,omitempty"`
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

// expectedFiles is a safeguard to ensure that the proper files
// are being used to load environmental config data
func expectedFiles() map[string]bool {
	return map[string]bool{
		"conf.development.yaml": true,
		"conf.testing.yaml":     true,
		"conf.staging.yaml":     true,
		"conf.production.yaml":  true,
	}
}

// isExpectedConfigPath checks that the input path has expected name format
func isExpectedConfigPath(cfgPath string) error {
	fileName := cfgPath[strings.LastIndex(cfgPath, "/")+1:]
	files := expectedFiles()
	if val, status := files[fileName]; !(val && status) {
		return fmt.Errorf("filename must be recognized")
	}
	if _, errPath := os.Stat(cfgPath); errPath != nil {
		return fmt.Errorf("error finding the path, %s", cfgPath)
	}
	log.Printf("config file: %s", cfgPath)
	return nil
}

// readFileAndBuildStruct reads the input file and builds a config struct
// that is serialized from all the data in the config rile
func readFileAndBuildStruct(cfgPath string) (*Configuration, error) {
	bytes, errRead := ioutil.ReadFile(cfgPath)
	if errRead != nil {
		return nil, fmt.Errorf("error reading config file, %v", errRead)
	}
	var cfg = new(Configuration)
	if errYaml := yaml.Unmarshal(bytes, cfg); errYaml != nil {
		return nil, fmt.Errorf("unable to decode config yaml into struct, %v", errYaml)
	}
	return cfg, nil
}

// LoadConfigFromFlags returns Configuration struct compiled from flags
// or it uses the default DevelopmentConfigPath() from the support package
func LoadConfigFromFlags() (*Configuration, error) {
	cfgPath := flag.String("config", support.DevelopmentConfigPath(), "Path to config file")
	flag.Parse()

	if errName := isExpectedConfigPath(*cfgPath); errName != nil {
		return nil, errName
	}
	return readFileAndBuildStruct(*cfgPath)
}

// LoadConfigFrom returns Configuration struct compiled from input path
func LoadConfigFrom(path string) (*Configuration, error) {
	return readFileAndBuildStruct(path)
}
