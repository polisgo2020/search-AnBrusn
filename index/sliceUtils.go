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
