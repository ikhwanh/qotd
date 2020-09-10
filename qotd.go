package main

import (
	"database/sql"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gookit/color"
	"github.com/ikhwanh/qotd/cfg"
	"github.com/urfave/cli/v2"

	_ "github.com/mattn/go-sqlite3"
)

var (
	qotd    cfg.Qotd
	gNumber int
)

const (
	totalAyat = 6236
)

func main() {

	config, err := cfg.New(cfg.DefaultPath())
	gNumber = 10

	configPath := cfg.Lookup()
	if configPath == "" {
		// No config file found, loading defaults
	} else {
		config.SetURL("file://" + configPath)
		// Loading config file from path
		err = config.Load()
		if err != nil {
			log.Fatal("Error loading user config file")
		}
	}

	setupConfig(config)

	qotd = config.Qotds[config.Cursor]

	app := &cli.App{
		Name:    "qotd",
		Version: "v0.0.1",
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "Ikhwan Hafidhi",
				Email: "ikhwanh.dev@gmail.com",
			},
		},
		Copyright: "(c) 2020 Ikhwan Hafidhi",
		Usage:     "show surah of quran when new terminal window opened",
		Action:    show,
		Commands: []*cli.Command{
			&cli.Command{
				Name:        "generate",
				Aliases:     []string{"g"},
				Description: "generate random quran's surah and cache",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "Number", Aliases: []string{"n"}, Value: 10, Destination: &gNumber},
				},
				Action: func(c *cli.Context) error {
					qotds, err := generate()
					if err != nil {
						return err
					}

					config.Qotds = qotds
					if err != nil {
						return err
					}

					return nil
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	// save new config
	err = config.Save()
	if err != nil {
		log.Fatal(err)
	}

}

func setupConfig(c *cfg.Config) {
	if c.IsRefreshNeeded() {
		qotds, err := generate()

		if err != nil {
			log.Fatal(err)
		}

		c.Cursor = 0
		c.DayLastUpdated = time.Now().Day()
		c.Qotds = qotds
	}

	if c.IsNewDay() {
		// new day has come so move cursor
		c.Cursor = c.Cursor + 1
		c.SetNewDay()
	}
}

func show(c *cli.Context) error {
	color.Bold.Printf("-------------------Quran Of The Day-------------------\n")
	color.Bold.Printf("Quran Surah: ")
	color.Normal.Printf("%d:%d\n", qotd.SurahIndex, qotd.Ayat)
	color.Bold.Printf("Nama: ")
	color.Normal.Printf("%s\n", qotd.SurahName)
	color.Bold.Printf("Terjemahan: \n")
	color.Normal.Printf("%s\n\n", qotd.Translation)

	return nil
}

func generate() ([]cfg.Qotd, error) {
	db, err := sql.Open("sqlite3", cfg.DataPath())

	if err != nil {
		return nil, err
	}
	defer db.Close()

	rand.Seed(time.Now().UnixNano())

	stmt, err := db.Prepare("select surah.name, surah.no, ayat.no, ayat.indo from ayat join surah on ayat.surah_no == surah.no limit 1 offset ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	qotds := make([]cfg.Qotd, gNumber)
	for i := 0; i < gNumber; i++ {
		rand := rand.Intn(totalAyat)

		var surahName string
		var surahNo int
		var ayatNo int
		var translation string
		err = stmt.QueryRow(rand).Scan(&surahName, &surahNo, &ayatNo, &translation)
		qotds[i] = cfg.Qotd{
			Ayat:        ayatNo,
			SurahIndex:  surahNo,
			SurahName:   surahName,
			Translation: translation,
		}

		if err != nil {
			return nil, err
		}
	}

	return qotds, nil

}
