# Inverted Index

github.com/polisgo2020/search-AnBrusn implements inverted index to perform full-text search.

## Usage

### Build index

Save inverted index into database (connection options from environment):

    go run main.go build -dir ~/path/to/text/files/
    
Save into file:

    go run main.go build -dir ~/path/to/text/files/ -index ~/path/to/output/file/

### Search over the index

Search in index from database:

    go run main.go search
    
Search in index from file:
    
    go run main.go search -index ~/path/to/output/file/
    
Search using web-interface (server address from environment):
    
    go run main.go search -http
    go run main.go search -index ~/path/to/output/file/ -http

Search results are ranged by amount of found tokens.

## Project dependencies

-   [`Stopwords`](github.com/bbalet/stopwords)
-   [`Snowball`](github.com/kljensen/snowball)
-   [`Cli`](github.com/urfave/cli/v2)
-   [`Zerolog`](github.com/rs/zerolog)
-   [`Env`](github.com/caarlos0/env)
-   [`Go-pg`](github.com/go-pg/pg)
