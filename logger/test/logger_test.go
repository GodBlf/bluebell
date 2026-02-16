package test

import (
	"bluebell/logger"
	"bluebell/settings"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestInitLogger(t *testing.T) {
	oldGlobalLogger := zap.L()
	t.Cleanup(func() {
		zap.ReplaceGlobals(oldGlobalLogger)
	})

	allLogFile, err := os.CreateTemp("", "bluebell-all-*.log")
	require.NoError(t, err)
	require.NoError(t, allLogFile.Close())
	errorLogFile, err := os.CreateTemp("", "bluebell-error-*.log")
	require.NoError(t, err)
	require.NoError(t, errorLogFile.Close())
	allLogPath := filepath.Clean(allLogFile.Name())
	errorLogPath := filepath.Clean(errorLogFile.Name())
	t.Cleanup(func() {
		_ = os.Remove(allLogPath)
		_ = os.Remove(errorLogPath)
	})

	cfg := &settings.LoggerConfig{
		Level:        "warn",
		Filename:     allLogPath,
		ErrorName:    errorLogPath,
		MaxSize:      10,
		MaxBackups:   1,
		MaxAge:       1,
		Compress:     false,
		LogInConsole: false,
	}

	logger.InitLogger(cfg)

	l := logger.L(context.Background())
	require.NotNil(t, l)
	l.Info("ignore-info")
	l.Error("record-error")
	_ = l.Sync()

	allLogBytes, err := os.ReadFile(allLogPath)
	require.NoError(t, err)
	errorLogBytes, err := os.ReadFile(errorLogPath)
	require.NoError(t, err)

	allLogContent := string(allLogBytes)
	errorLogContent := string(errorLogBytes)

	assert.Equal(t, zap.WarnLevel, logger.AtomicLevel.Level())
	assert.Contains(t, allLogContent, "record-error")
	assert.NotContains(t, allLogContent, "ignore-info")
	assert.Contains(t, errorLogContent, "record-error")
	assert.NotContains(t, errorLogContent, "ignore-info")
}
