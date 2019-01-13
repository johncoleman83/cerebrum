package support_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/johncoleman83/cerebrum/pkg/utl/support"
)

func TestTestingConfigPath(t *testing.T) {
	t.Run("test should get the proper testing config path", func(t *testing.T) {
		cfgPath := support.TestingConfigPath()
		fileName := cfgPath[strings.LastIndex(cfgPath, "/")+1:]
		assert.Equal(t, "conf.testing.yaml", fileName, "filename should be the correct test filename")
	})
}

func TestDevelopmentConfigPath(t *testing.T) {
	t.Run("test should get the proper development config path", func(t *testing.T) {
		cfgPath := support.DevelopmentConfigPath()
		fileName := cfgPath[strings.LastIndex(cfgPath, "/")+1:]
		assert.Equal(t, "conf.development.yaml", fileName, "filename should be the correct test filename")
	})
}

func TestExtractPathFromFlags(t *testing.T) {
	t.Run("test should return default path if no flags provided", func(t *testing.T) {
		cfgPath, err := support.ExtractPathFromFlags()
		assert.Nil(t, err, "there should not be an error when extracting the path from flags")

		fileName := cfgPath[strings.LastIndex(cfgPath, "/")+1:]
		assert.Equal(t, "conf.development.yaml", fileName, "filename should be the correct test filename")
	})
}
