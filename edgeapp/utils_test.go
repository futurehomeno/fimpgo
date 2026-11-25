package edgeapp

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupLog(t *testing.T) {
	origFormatter, origLevel, origOut := log.StandardLogger().Formatter, log.GetLevel(), log.StandardLogger().Out
	t.Cleanup(func() {
		log.SetFormatter(origFormatter)
		log.SetLevel(origLevel)
		log.SetOutput(origOut)
	})

	tests := []struct {
		format, level string
		want          log.Level
		fmtType       string
	}{
		{"json", "warn", log.WarnLevel, "*logrus.JSONFormatter"},
		{"budzik", "error", log.ErrorLevel, "*formatters.BudzikFormatter"},
		{"", "bogus", log.DebugLevel, "*logrus.TextFormatter"},
	}
	for _, tt := range tests {
		SetupLog("", tt.level, tt.format)
		assert.Equal(t, tt.want, log.GetLevel())
		assert.Equal(t, tt.fmtType, fmt.Sprintf("%T", log.StandardLogger().Formatter))
	}

	file := filepath.Join(t.TempDir(), "app.log")
	SetupLog(file, "info", "json")
	log.Info("hello")
	b, err := os.ReadFile(file) //nolint:gosec
	require.NoError(t, err)
	assert.Contains(t, string(b), "hello")
}

func TestSetLogLevel(t *testing.T) {
	orig := log.GetLevel()
	t.Cleanup(func() { log.SetLevel(orig) })
	SetLogLevel("error")
	assert.Equal(t, log.ErrorLevel, log.GetLevel())
	SetLogLevel("bogus")
	assert.Equal(t, log.DebugLevel, log.GetLevel())
}
