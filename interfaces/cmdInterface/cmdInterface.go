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
	in            *os.File
	out           *os.File
	invertedIndex *index.Index
}

func New(in *os.File, out *os.File, invertedIndex *index.Index) (*CmdInterface, error) {
	log.Debug().Msg("create command line user interface")
	if in == nil || out == nil || invertedIndex == nil {
		return nil, errors.New("invalid in, out or index object")
	}
	return &CmdInterface{
		in:            in,
		out:           out,
		invertedIndex: invertedIndex,
	}, nil
}

// Run reads user input from in and searches over inverted index.
func (c *CmdInterface) Run() error {
	for {
		scanner := bufio.NewScanner(c.in)
		scanner.Scan()
		userInput := scanner.Text()
		log.Debug().Str("text", userInput).Msg("new search request")
		searchResults, err := c.invertedIndex.FindInIndex(userInput)
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