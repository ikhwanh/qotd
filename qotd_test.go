package main

import (
	"testing"
	"time"

	"github.com/ikhwanh/qotd/cfg"
)

func TestSetupConfig(t *testing.T) {
	config := &cfg.Config{
		DayLastUpdated: 0,
		Cursor:         0,
	}

	setupConfig(config)

	day := time.Now().Day()

	if config.Cursor != 0 {
		t.Errorf("Cursor is %d; want 0", config.Cursor)
	}

	if config.DayLastUpdated != day {
		t.Errorf("DayLastUpdated is %d; want %d", config.DayLastUpdated, day)
	}

	// test new day has come
	config.DayLastUpdated = day - 1
	config.Cursor = 1
	config.Qotds = make([]cfg.Qotd, 10)

	setupConfig(config)

	if config.Cursor != 2 {
		t.Errorf("Cursor is %d; want 2", config.Cursor)
	}

	if config.DayLastUpdated != day {
		t.Errorf("DayLastUpdated is %d; want %d", config.DayLastUpdated, day)
	}

	// test cursor out of index
	config.Cursor = 10

	setupConfig(config)

	if config.Cursor != 0 {
		t.Errorf("Cursor is %d; cursor should renew to 0", config.Cursor)
	}

}
