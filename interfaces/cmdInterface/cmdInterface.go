package cmdInterface

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/polisgo2020/search-AnBrusn/index"
	"github.com/rs/zerolog/log"
)

type CmdInterface struct {
	in         *os.File
	out        *os.File
	searchFunc func(userInput string) ([]index.FileWithFreq, error)
}

func New(in *os.File, out *os.File, searchFunc func(userInput string) ([]index.FileWithFreq, error)) (*CmdInterface, error) {
	log.Info().Msg("create command line user interface")
	if in == nil || out == nil {
		return nil, errors.New("invalid in or out")
	}
	return &CmdInterface{
		in:         in,
		out:        out,
		searchFunc: searchFunc,
	}, nil
}

// Run reads user input from in and searches over inverted index.
func (c *CmdInterface) Run() error {
	for {
		scanner := bufio.NewScanner(c.in)
		scanner.Scan()
		userInput := scanner.Text()
		log.Debug().Str("text", userInput).Msg("new search request")
		searchResults, err := c.searchFunc(userInput)
		if err != nil {
			return err
		}
		if len(searchResults) == 0 {
			fmt.Fprintln(c.out, "No results")
		} else {
			for _, el := range searchResults {
				fmt.Fprintf(c.out, "%s (%d words were found)\n", el.Filename, el.Freq)
			}
		}
	}
	return nil
}
