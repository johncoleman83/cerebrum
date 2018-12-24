package config

import (
	"fmt"
	"os"
)

// Load Configuration from Environment
func LoadEnvironment() (error) {
	envs := map[string]bool{
    "dev": true,
    "stg": true,
		"prod": true,
	}
	// set ENVIRONMENT_NAME to dev if unset
	if val, status := envs[os.Getenv("ENVIRONMENT_NAME")]; !(val && status) {
    os.Setenv("ENVIRONMENT_NAME", "dev")
	}
	return nil
}
