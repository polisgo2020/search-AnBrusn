package index

import (
	"reflect"
	"testing"
)

func TestAppendIfMissing(t *testing.T) {
	int := []string{"a", "b", "c"}
	expected := []string{"a", "b", "c", "d"}
	actual := AppendIfMissing(int, "c")
	actual = AppendIfMissing(actual, "d")
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nactual:   %v \nexpected: %v", actual, expected)
	}
}
