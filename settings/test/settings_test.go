package test

import (
	"bluebell/settings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootPath(t *testing.T) {
	p := settings.RootPath
	t.Run("测试RootPath", func(t *testing.T) {
		assert.Equal(t, "D:/Projects/bluebell/settings/../", p, "RootPath should be D:/Projects/bluebell/")

	})
}

func TestInitConfig(t *testing.T) {
	settings.InitAppConfig()
	config := settings.GlobalConfig
	t.Run("测试AppConfig", func(t *testing.T) {
		assert.NotNil(t, config, "GlobalConfig should not be nil")
		assert.Equal(t, "bluebell", config.Name, "App name should be bluebell")
		assert.Equal(t, "dev", config.Mode, "App mode should be dev")
		assert.Equal(t, "v0.0.1", config.Version, "App version should be 0.1.0")
		assert.Equal(t, 8084, config.Port, "App port should be 8084")
	})

}
