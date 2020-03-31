package index

import (
	"reflect"
	"testing"
)

func TestIndex_AddToken(t *testing.T) {
	in := make(Index)
	in["plenti"] = []FileWithFreq{{"source1.txt", 1}}
	in["differ"] = []FileWithFreq{{"source1.txt", 1}}
	in["token"] = []FileWithFreq{{"source1.txt", 1}}

	expected := make(Index)
	expected["plenti"] = []FileWithFreq{{"source1.txt", 2}}
	expected["differ"] = []FileWithFreq{{"source1.txt", 1}, {"source2.txt", 1}}
	expected["token"] = []FileWithFreq{{"source1.txt", 1}}
	expected["good"] = []FileWithFreq{{"source1.txt", 1}}

	if err := in.AddToken("plenty", "source1.txt"); err != nil {
		t.Errorf("error occured\nexpected: %v", expected)
	}
	if err := in.AddToken("differ", "source2.txt"); err != nil {
		t.Errorf("error occured\nexpected: %v", expected)
	}
	if err := in.AddToken("good", "source1.txt"); err != nil {
		t.Errorf("error occured\nexpected: %v", expected)
	}

	if !reflect.DeepEqual(in, expected) {
		t.Errorf("\nactual:   %v \nexpected: %v", in, expected)
	}
}
