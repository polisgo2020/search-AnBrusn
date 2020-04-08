/*
Package index implements inverted index, search over the built index.

Usage

New token can be added with AddToken function that extract token from word and add it to inverted index.
Example:

	err := i.AddToken(" word", "sourceFile.txt")

To search over the index use FindInIndex function.

	searchResults, err := i.FindInIndex("this is search query")

Search results are ranged by amount of found tokens.
*/
package index

// AppendIfMissing add an element in slice if it is missing
func AppendIfMissing(slice []string, newElement string) []string {
	for _, el := range slice {
		if el == newElement {
			return slice
		}
	}
	return append(slice, newElement)
}
