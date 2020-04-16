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
	srv        *http.Server
	searchFunc func(userInput string) ([]index.FileWithFreq, error)
}

func New(srv *http.Server, searchFunc func(userInput string) ([]index.FileWithFreq, error)) (*WebInterface, error) {
	log.Info().Str("server", srv.Addr).Msg("create http user interface")
	if srv == nil {
		return nil, errors.New("invalid server")
	}
	return &WebInterface{
		srv:        srv,
		searchFunc: searchFunc,
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
		searchResults, err := web.searchFunc(userInput)
		if err != nil {
			log.Err(err).Msg("error while searching")
			http.Error(w, "searching error", http.StatusInternalServerError)
			return
		}
		tmplResults.Execute(w, searchResults)
	})
	return web.srv.ListenAndServe()
}
