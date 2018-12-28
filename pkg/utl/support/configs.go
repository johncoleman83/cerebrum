package support

import (
	"path/filepath"
	"runtime"
	"strings"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func configsDirectoryFullPath() string {
	tail := strings.LastIndex(basepath, "pkg/utl/support")
	return basepath[:tail] + "configs/"
}

// TestingConfigPath returns the path of the testing configuration yaml
func TestingConfigPath() string {
	return configsDirectoryFullPath() + "conf.testing.yaml"
}

// DevelopmentConfigPath returns the path of the development configuration yaml
func DevelopmentConfigPath() string {
	return configsDirectoryFullPath() + "conf.development.yaml"
}
