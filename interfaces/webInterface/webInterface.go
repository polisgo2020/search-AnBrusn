package webInterface

import (
	"errors"
	"html/template"
	"net/http"
	"path"

	"github.com/polisgo2020/search-AnBrusn/index"
	"github.com/rs/zerolog/log"
)

type WebInterface struct {
	srv           *http.Server
	invertedIndex *index.Index
}

func New(srv *http.Server, invertedIndex *index.Index) (*WebInterface, error) {
	log.Debug().Str("server", srv.Addr).Msg("create http user interface")
	if srv == nil || invertedIndex == nil {
		return nil, errors.New("invalid server or index object")
	}
	return &WebInterface{
		srv:           srv,
		invertedIndex: invertedIndex,
	}, nil
}

// Run reads user input with http user interface and searches over inverted index.
func (web *WebInterface) Run() error {
	tmplSearch, err := template.ParseFiles(path.Join("templates", "search.html"))
	if err != nil {
		return err
	}
	tmplResults, err := template.ParseFiles(path.Join("templates", "results.html"))
	if err != nil {
		return err
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmplSearch.Execute(w, nil)
	})
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		userInput := r.FormValue("userInput")
		log.Debug().Str("text", userInput).Msg("new search request")
		if err != nil {
			log.Err(err).Msg("error while reading index")
			http.Error(w, "reading index error", http.StatusInternalServerError)
			return
		}
		searchResults, err := web.invertedIndex.FindInIndex(userInput)
		if err != nil {
			log.Err(err).Msg("error while searching")
			http.Error(w, "searching error", http.StatusInternalServerError)
			return
		}
		tmplResults.Execute(w, searchResults)
	})
	if err := web.srv.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
