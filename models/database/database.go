package database

import (
	"context"
	"fmt"
	"github.com/go-pg/pg"
	"github.com/polisgo2020/search-AnBrusn/index"
	"github.com/rs/zerolog/log"
)

// Token entry
type Token struct {
	ID    int    `sql:"id,pk"`
	Token string `sql:"token"`
}

// File entry
type File struct {
	ID   int    `sql:"id,pk"`
	Name string `sql:"name"`
}

// Occurrence entry
type Occurrence struct {
	ID      int `sql:"id,pk"`
	TokenID int `sql:"token_id"`
	FileID  int `sql:"file_id"`
	Number  int `sql:"num"`
}

type Repository struct {
	pg       *pg.DB
	ErrChan  chan error
	DataChan chan [2]string
}

// New connects to database
func New(cnString string) (Repository, error) {
	pgOpt, err := pg.ParseURL(cnString)
	if err != nil {
		return Repository{}, err
	}
	log.Debug().
		Str("database", pgOpt.Database).
		Str("user", pgOpt.User).
		Str("address", pgOpt.Addr).
		Msg("create database connection")
	pgdb := pg.Connect(pgOpt)
	return Repository{pg: pgdb, ErrChan: make(chan error), DataChan: make(chan [2]string)}, nil
}

func (rep Repository) CloseConnection() {
	if err := rep.pg.Close(); err != nil {
		log.Err(err)
	}
}

// FindInIndex implements search in database
func (rep Repository) FindInIndex(userInput string) ([]index.FileWithFreq, error) {
	inputTokens, err := index.GetTokensFromInput(userInput)
	if err != nil {
		return nil, fmt.Errorf("can not get tokens from user input %w", err)
	}
	var res []index.FileWithFreq
	err = rep.pg.Model(new(File)).
		ColumnExpr("file.name AS filename").
		ColumnExpr("sum(occurrences.num) AS freq").
		Join("JOIN occurrences ON file.id = occurrences.file_id").
		Join("JOIN tokens ON tokens.id = occurrences.token_id").
		WhereIn("tokens.token IN (?)", inputTokens).
		Group("file.name").
		Having("count(tokens.id)=(?)", len(inputTokens)).
		OrderExpr("sum(occurrences.num) desc").
		Select(&res)
	if err != nil {
		return nil, fmt.Errorf("can not select from database %w", err)
	}
	return res, nil
}

// AddToken adds occurrence of token into database
func (rep Repository) AddToken(word string, filename string) error {
	stemWord, err := index.GetTokenFromWord(word)
	if err != nil {
		return fmt.Errorf("can not extract token from word %w", err)
	}
	if len(stemWord) > 0 {
		var fileId int
		_, err := rep.pg.Model(&File{Name: filename}).
			Column("id").
			Where("name = ?", filename).
			Returning("id").
			SelectOrInsert(&fileId)
		if err != nil {
			return fmt.Errorf("can not insert file %w", err)
		}
		var tokenId int
		_, err = rep.pg.Model(&Token{Token: stemWord}).
			Column("id").
			Where("token = ?", stemWord).
			Returning("id").
			SelectOrInsert(&tokenId)
		if err != nil {
			return fmt.Errorf("can not insert token %w", err)
		}
		exists, err := rep.pg.Model(new(Occurrence)).
			Where("token_id = ?", tokenId).
			Where("file_id = ?", fileId).
			Exists()
		if err != nil {
			return fmt.Errorf("can not insert occurence %w", err)
		}
		if !exists {
			_, err = rep.pg.Model(&Occurrence{TokenID: tokenId, FileID: fileId, Number: 1}).Insert()
		} else {
			_, err = rep.pg.Model(new(Occurrence)).
				Set("num = num + 1").
				Where("token_id = ?", tokenId).
				Where("file_id = ?", fileId).
				Update()
		}
		if err != nil {
			return fmt.Errorf("cannot insert occurrence %w", err)
		}
	}
	return nil
}

// Listener listens channel of words and adds them in database
func (rep Repository) Listener(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case wordInfo := <-rep.DataChan:
			if err := rep.AddToken(wordInfo[0], wordInfo[1]); err != nil {
				rep.ErrChan <- err
			}
		}
	}
}
