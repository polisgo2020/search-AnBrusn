# Inverted Index

github.com/polisgo2020/search-AnBrusn implements inverted index to perform full-text search.

## Usage

### Build index

    go run main.go build -dir ~/path/to/text/files/ -index ~/path/to/output/file/

### Search over the index

    go run main.go search -index ~/path/to/output/file/
    go run main.go search -index ~/path/to/output/file/ -http :8080

Search results are ranged by amount of found tokens.

## Project dependencies

-   [`Stopwords`](github.com/bbalet/stopwords)
-   [`Snowball`](github.com/kljensen/snowball)
-   [`Cli`](github.com/urfave/cli/v2)
-   [`Zerolog`](github.com/rs/zerolog)
