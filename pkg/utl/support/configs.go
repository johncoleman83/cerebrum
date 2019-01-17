package support

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Directory Paths
var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
	circlepath = "/go/src/github.com/johncoleman83/cerebrum/"

	// hack for circle tests
	configsdir = "configs/application/"
)

func configsDirectoryFullPath() string {
	if strings.Contains(basepath, "_test/_obj_test") {
		return circlepath + configsdir
	}
	tail := strings.LastIndex(basepath, "pkg/utl/support")
	return basepath[:tail] + configsdir
}

// expectedFiles is simply a list of expected files for error checking
func expectedFiles() map[string]bool {
	return map[string]bool{
		"conf.development.yaml": true,
		"conf.testing.yaml":     true,
		"conf.staging.yaml":     true,
		"conf.production.yaml":  true,
	}
}

// checks if path has expected name format
func isExpectedConfigPath(cfgPath string) error {
	fileName := cfgPath[strings.LastIndex(cfgPath, "/")+1:]
	files := expectedFiles()
	if val, status := files[fileName]; !(val && status) {
		return fmt.Errorf("filename must be recognized")
	}
	if _, err := os.Stat(cfgPath); err != nil {
		return fmt.Errorf("error finding the path, %s", cfgPath)
	}
	log.Printf("config file: %s", cfgPath)
	return nil
}

// TestingConfigPath returns the path of the testing configuration yaml
func TestingConfigPath() string {
	return configsDirectoryFullPath() + "conf.testing.yaml"
}

// DevelopmentConfigPath returns the path of the development configuration yaml
func DevelopmentConfigPath() string {
	return configsDirectoryFullPath() + "conf.development.yaml"
}

// ExtractPathFromFlags returns path string from flags or default path
func ExtractPathFromFlags() (string, error) {
	cfgPath := flag.String("config", DevelopmentConfigPath(), "Path to config file")
	flag.Parse()

	if err := isExpectedConfigPath(*cfgPath); err != nil {
		return "", err
	}
	return *cfgPath, nil
}
