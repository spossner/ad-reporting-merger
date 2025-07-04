package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	cfg, err := LoadConfig()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	groups := cfg.GetGroups()
	assert.Len(t, groups, 2)

	// Check first group
	assert.Equal(t, "AdManager Reporting", groups[0].Prefix)
	assert.Equal(t, "raw.csv", groups[0].Output)

	// Check second group
	assert.Equal(t, "Revenue per AdUnit", groups[1].Prefix)
	assert.Equal(t, "raw-revenue.csv", groups[1].Output)
}

func TestGetWorkDir(t *testing.T) {
	cfg, err := LoadConfig()
	require.NoError(t, err)

	workDir := cfg.GetWorkDir()
	assert.Equal(t, "~/Downloads", workDir)
}

func TestGetWorkDirDefault(t *testing.T) {
	cfg := &Config{} // Empty config
	workDir := cfg.GetWorkDir()
	assert.Equal(t, "~/Downloads", workDir)
}