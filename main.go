package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unicode"

	"github.com/polisgo2020/search-AnBrusn/config"
	"github.com/polisgo2020/search-AnBrusn/index"
	"github.com/polisgo2020/search-AnBrusn/interfaces/cmdInterface"
	"github.com/polisgo2020/search-AnBrusn/interfaces/webInterface"
	"github.com/polisgo2020/search-AnBrusn/models/database"
	"github.com/polisgo2020/search-AnBrusn/models/file"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var cfg config.Config

// readFile extracts tokens from file and adds them in inverted index.
func readFile(ctx context.Context, path string, filename string,
	dataChan chan<- [2]string, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	data, err := ioutil.ReadFile(path)
	if err != nil {
		errChan <- fmt.Errorf("can not read file %w", err)
		return
	}
	words := strings.FieldsFunc(string(data), func(r rune) bool {
		return !unicode.IsLetter(r)
	})
	for _, word := range words {
		select {
		case <-ctx.Done():
			return
		default:
			dataChan <- [2]string{word, filename}
		}
	}
}

// createFromDirectory extracts tokens from files in directory and adds them in inverted index.
func createFromDirectory(dirname string, dataChan chan [2]string, errChan chan error, listener func(ctx context.Context)) error {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return  fmt.Errorf("can not access input directory %w", err)
	}

	defer close(errChan)
	defer close(dataChan)

	ctx, cancel := context.WithCancel(context.Background())

	go listener(ctx)

	wgForFilesReading := &sync.WaitGroup{}
	wgForFilesReading.Add(len(files))
	go func() {
		wgForFilesReading.Wait()
		cancel()
	}()

	for _, currentFile := range files {
		path := filepath.Join(dirname, currentFile.Name())
		go readFile(ctx, path, currentFile.Name(), dataChan, errChan, wgForFilesReading)
	}

	for {
		select {
		case err := <-errChan:
			cancel()
			return err
		case <-ctx.Done():
			return nil
		}
	}
}

// searchInIndex runs web or cmd interface and perform search over file or database
func searchInIndex(c *cli.Context) error {
	var invertedIndex index.Index
	var err error
	var searchFunc func(userInput string) ([]index.FileWithFreq, error)
	if c.String("index") != "" {
		log.Debug().Str("index", c.String("index")).Msg("searching in index")
		invertedIndex, err = file.ReadIndex(c.String("index"))
		if err != nil {
			return fmt.Errorf("can not read index %w", err)
		}
		searchFunc = invertedIndex.FindInIndex
	} else {
		log.Debug().Msg("searching in database")
		rep, err := database.New(cfg.PgSQL)
		if err != nil {
			return fmt.Errorf("can not connect to database %w", err)
		}
		defer rep.CloseConnection()
		searchFunc = rep.FindInIndex
	}
	if c.Bool("http") == true {
		srv := &http.Server{Addr: cfg.Server}
		w, err := webInterface.New(srv, searchFunc)
		if err != nil {
			return err
		}
		return w.Run()
	} else {
		c, err := cmdInterface.New(os.Stdin, os.Stdout, searchFunc)
		if err != nil {
			return err
		}
		return c.Run()
	}
}

// createIndexFromFiles creates index ans saves it in file or database
func createIndexFromFiles(c *cli.Context) error {
	log.Debug().Str("directory", c.String("dir")).Msg("building index")
	if c.String("index") != "" {
		i := index.NewIndex()
		if err := createFromDirectory(c.String("dir"), i.DataChan, i.ErrChan, i.Listener); err != nil {
			return fmt.Errorf("can not create index fron directory %w", err)
		}
		log.Debug().Str("index", c.String("index")).Msg("into file")
		return file.WriteInvertedIndex(c.String("index"), i)
	} else {
		log.Debug().Msg("into database")
		rep, err := database.New(cfg.PgSQL)
		if err != nil {
			return fmt.Errorf("can not connect to database %w", err)
		}
		defer rep.CloseConnection()
		if err := createFromDirectory(c.String("dir"), rep.DataChan, rep.ErrChan, rep.Listener); err != nil {
			return fmt.Errorf("can not create index fron directory %w", err)
		}
	}
	return nil
}

func main() {
	var err error
	cfg, err = config.Load()
	if err != nil {
		log.Err(err).Msg("error while getting system env")
	}
	logLevel, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Err(err).Msg("error while getting loglevel from system env")
	}
	zerolog.SetGlobalLevel(logLevel)
	app := &cli.App{
		Name:  "Index",
		Usage: "Create index from directory and search in index",
		Commands: []*cli.Command{
			{
				Name:  "build",
				Usage: "Create index from directory",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "dir",
						Usage:    "Path to directory with input files",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "index",
						Usage:    "Path to output index file",
						Required: false,
					},
				},
				Action: createIndexFromFiles,
			},
			{
				Name:  "search",
				Usage: "Search in index",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "index",
						Usage:    "Path to index file",
						Required: false,
					},
					&cli.BoolFlag{
						Name:  "http",
						Usage: "Input from http",
					},
				},
				Action: searchInIndex,
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Err(err)
	}
}
