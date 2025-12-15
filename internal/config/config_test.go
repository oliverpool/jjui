package config

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
)

func TestLoad_Colours(t *testing.T) {
	content := `
[ui.colors]
"text" = "white"
"selected" = { fg = "blue", bg = "black" }
`
	config := &Config{}
	err := config.Load(content)
	assert.NoError(t, err)
	assert.Len(t, config.UI.Colors, 2)
	assert.Equal(t, "white", config.UI.Colors["text"].Fg)
	assert.Equal(t, "blue", config.UI.Colors["selected"].Fg)
	assert.Equal(t, "black", config.UI.Colors["selected"].Bg)
}

func TestLoad_Theme_Simple(t *testing.T) {
	content := `
[ui]
theme = "my-theme"
`
	config := &Config{}
	err := config.Load(content)
	assert.NoError(t, err)
	assert.Equal(t, "my-theme", config.UI.Theme.Light)
	assert.Equal(t, "my-theme", config.UI.Theme.Dark)
}

func TestLoad_Theme_Nested(t *testing.T) {
	content := `
[ui.theme]
dark = "dark-theme"
light = "light-theme"
`
	config := &Config{}
	err := config.Load(content)
	assert.NoError(t, err)
	assert.Equal(t, "dark-theme", config.UI.Theme.Dark)
	assert.Equal(t, "light-theme", config.UI.Theme.Light)
}

func TestLoad_AutoRefreshInterval(t *testing.T) {
	content := `
[ui]
auto_refresh_interval = 5000
`
	config := &Config{}
	err := config.Load(content)
	assert.NoError(t, err)
	assert.Equal(t, 5000, config.UI.AutoRefreshInterval)
}

func TestLoad_Colors_StringAndObject(t *testing.T) {
	content := `
[ui.colors]
simple = "red"
complex = { fg = "blue", bg = "white", bold = true }
`
	config := &Config{}
	err := config.Load(content)
	assert.NoError(t, err)
	assert.Len(t, config.UI.Colors, 2)

	assert.Equal(t, "red", config.UI.Colors["simple"].Fg)
	assert.Equal(t, "", config.UI.Colors["simple"].Bg)
	assert.False(t, config.UI.Colors["simple"].Bold)

	assert.Equal(t, "blue", config.UI.Colors["complex"].Fg)
	assert.Equal(t, "white", config.UI.Colors["complex"].Bg)
	assert.True(t, config.UI.Colors["complex"].Bold)
}

func TestDefault_IsUpToDate(t *testing.T) {
	var formatted bytes.Buffer
	assert.NoError(t, toml.NewEncoder(&formatted).Encode(loadDefaultConfig()))

	raw, err := configFS.ReadFile("default/config.toml")
	assert.NoError(t, err)

	expected := strings.NewReplacer(
		// the toml library does not support round-tripping with comments
		// https://github.com/BurntSushi/toml/issues/213
		`template = ""`, `# template = 'builtin_log_compact' # overrides jj's templates.log`,
		`revset = ""`, `# revset = "zzzzzzz"               # overrides jj's revsets.log`,
	).Replace(formatted.String())

	if !assert.Equal(t, expected, string(raw), "The internal/config/default/config.toml does not seem to be up to date\nRun the test with JJUI_TEST_UPDATE=config to update it") && os.Getenv("JJUI_TEST_UPDATE") == "config" {
		err = os.WriteFile("default/config.toml", []byte(expected), 0o666)
		assert.NoError(t, err)
	}
}
