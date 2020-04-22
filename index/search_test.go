package index

import (
	"reflect"
	"testing"
)

func TestGetTokensFromInput(t *testing.T) {
	in := "plenty, different  tokens."
	expected := []string{"plenti", "differ", "token"}
	actual, err := GetTokensFromInput(in)

	if err != nil {
		t.Errorf("error occured\nexpected: %v", expected)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nactual:   %v \nexpected: %v", actual, expected)
	}
}

func TestIndex_FindInIndex(t *testing.T) {
	in := NewIndex()
	in.Data["plenti"] = []FileWithFreq{{"source1.txt", 2}}
	in.Data["differ"] = []FileWithFreq{{"source1.txt", 1}, {"source2.txt", 2}}
	in.Data["token"] = []FileWithFreq{{"source1.txt", 1}}
	in.Data["good"] = []FileWithFreq{{"source1.txt", 1}}

	expected := []FileWithFreq{{"source2.txt", 2}, {"source1.txt", 1}}
	actual, err := in.FindInIndex("different")
	if err != nil {
		t.Errorf("error occured\nexpected: %v", expected)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nactual:   %v \nexpected: %v", actual, expected)
	}

	var expectedNotFound []FileWithFreq
	actual, err = in.FindInIndex("NoSuchToken")
	if err != nil {
		t.Errorf("error occured\nexpected: %v", expectedNotFound)
	}
	if !reflect.DeepEqual(actual, expectedNotFound) {
		t.Errorf("\nactual:   %v \nexpected: %v", actual, expectedNotFound)
	}
}

func TestGetIntersection(t *testing.T) {
	in1 := []FileWithFreq{{"file1", 4}, {"file2", 3}, {"file3", 2}}
	in2 := []FileWithFreq{{"file1", 2}, {"file3", 1}, {"file5", 3}}

	expected := []FileWithFreq{{"file1", 6}, {"file3", 3}}
	actual := getIntersection(in1, in2)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("\nactual:   %v \nexpected: %v", actual, expected)
	}
}
