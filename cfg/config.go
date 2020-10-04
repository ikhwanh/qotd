package cfg

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"

	gap "github.com/muesli/go-app-paths"
)

const appName = "qotd"

var cfgFileName = "qotd.conf"
var dataDbName = "quran.db"

// Qotd is quran of the day
type Qotd struct {
	SurahName   string
	SurahIndex  int
	Ayat        int
	Translation string
}

// Config contains an entire configuration
type Config struct {
	DayLastUpdated int
	Cursor         int
	Qotds          []Qotd
	backend        ConfigBackend
	url            *url.URL
}

// ConfigBackend is the interface implemented by the configuration backends.
type ConfigBackend interface {
	Load(*url.URL) (*Config, error)
	Save(*Config) error
}

// Save the current configuration.
func (c *Config) Save() error {
	return c.backend.Save(c)
}

// Load the configuration.
func (c *Config) Load() error {
	config, err := c.backend.Load(c.url)
	if err != nil {
		return err
	}
	c.DayLastUpdated = config.DayLastUpdated
	c.Cursor = config.Cursor
	c.Qotds = config.Qotds
	return nil
}

// Backend currently being used.
func (c *Config) Backend() ConfigBackend {
	return c.backend
}

// SetURL updates the configuration URL.
//
// Next time the config is loaded or saved
// the new URL will be used.
func (c *Config) SetURL(u string) error {
	url, err := url.Parse(u)
	if err != nil {
		return err
	}

	c.url = url

	return nil
}

// URL currently being used.
func (c *Config) URL() *url.URL {
	return c.url
}

// New returns a new Config struct.
//
// The URL will be matched against all the supported
// backends and the first backend that can handle the
// URL scheme will be loaded.
//
// A UNIX style path is also accepted, and will be handled
// by the default FileBackend.
func New(url string) (*Config, error) {
	config := &Config{}
	var backend ConfigBackend

	if url == "" {
		return nil, fmt.Errorf("Empty URL provided but not supported")
	}

	err := config.SetURL(url)
	if err != nil {
		return nil, err
	}

	switch config.url.Scheme {
	case "", "file":
		backend = NewFileBackend()
	default:
		return nil, fmt.Errorf("Configuration backend '%s' not supported", config.url.Scheme)
	}

	config.backend = backend

	return config, nil
}

// DefaultPath returns default config path.
//
// The path returned is OS dependant. If there's an error
// while trying to figure out the OS dependant path, "beehive.conf"
// in the current working dir is returned.
func DefaultPath() string {
	userScope := gap.NewScope(gap.User, appName)
	path, err := userScope.ConfigPath(cfgFileName)
	if err != nil {
		return cfgFileName
	}

	return path
}

// DataPath return default quran.db path
func DataPath() string {
	if exist(dataDbName) {
		return dataDbName
	}

	userScope := gap.NewScope(gap.User, appName)
	path, err := userScope.DataPath(dataDbName)

	if err != nil {
		return dataDbName
	}

	return path
}

// Lookup tries to find the config file.
//
// If a config file is found in the current working directory, that's returned.
// Otherwise we try to locate it following an OS dependant:
//
// Unix:
//   - ~/.config/app/filename.conf
// macOS:
//   - ~/Library/Preferences/app/filename.conf
// Windows:
//   - %LOCALAPPDATA%/app/Config/filename.conf
//
// If no valid config file is found, an empty string is returned.
func Lookup() string {
	paths := []string{}
	defaultPath := DefaultPath()
	if exist(defaultPath) {
		paths = append(paths, defaultPath)
	}

	// Prepend .conf to the search path if exists, takes priority
	// over the rest.
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current working directory")
		cwd = "."
	}
	cwdCfg := filepath.Join(cwd, cfgFileName)
	if exist(cwdCfg) {
		paths = append([]string{cwdCfg}, paths...)
	}
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

// IsNewDay is check now is new day
func (c *Config) IsNewDay() bool {
	if c.DayLastUpdated != time.Now().Day() {
		return true
	}
	return false
}

// SetNewDay is set DayLastUpdated to now day
func (c *Config) SetNewDay() {
	c.DayLastUpdated = time.Now().Day()
}

// IsRefreshNeeded is check qotds should be re-cache
func (c *Config) IsRefreshNeeded() bool {
	if len(c.Qotds) == 0 || c.Cursor > len(c.Qotds)-1 {
		return true
	}
	return false
}

func exist(file string) bool {
	_, err := os.Stat(file)
	if err == nil {
		return true
	}

	return false
}
