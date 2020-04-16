package database

import (
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
	pg *pg.DB
}

func New(cnString string) (Repository, error) {
	pgOpt, err := pg.ParseURL(cnString)
	if err != nil {
		return Repository{}, err
	}
	pgdb := pg.Connect(pgOpt)
	return Repository{pg: pgdb,}, nil
}

func closeConnection(pgdb *pg.DB) {
	if err := pgdb.Close(); err != nil {
		log.Err(err)
	}
}

// FindInIndex implements search in database
func (rep Repository) FindInIndex(userInput string) ([]index.FileWithFreq, error) {
	inputTokens, err := index.GetTokensFromInput(userInput)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return res, nil
}

// WriteInvertedIndex saves index into database
func (rep Repository) WriteInvertedIndex(invertedIndex index.Index) error {
	defer closeConnection(rep.pg)
	if _, err := rep.pg.Exec("TRUNCATE tokens, files, occurrences RESTART IDENTITY;"); err != nil {
		return err
	}

	for token, files := range invertedIndex {
		t := Token{Token: token}
		if err := rep.pg.Insert(&t); err != nil {
			return err
		}
		for _, fileInfo := range files {
			var f File
			if err := rep.pg.Model(&f).Where("name = ?", fileInfo.Filename).Select(); err != nil {
				f = File{Name: fileInfo.Filename}
				if err = rep.pg.Insert(&f); err != nil {
					return err
				}
			}
			o := Occurrence{TokenID: t.ID, FileID: f.ID, Number: fileInfo.Freq}
			if err := rep.pg.Insert(&o); err != nil {
				return err
			}
		}
	}
	return nil
}
