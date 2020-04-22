package index

import (
	"reflect"
	"testing"
)

func TestIndex_AddToken(t *testing.T) {
	in := NewIndex()
	in.Data["plenti"] = []FileWithFreq{{"source1.txt", 1}}
	in.Data["differ"] = []FileWithFreq{{"source1.txt", 1}}
	in.Data["token"] = []FileWithFreq{{"source1.txt", 1}}

	expected := NewIndex()
	expected.Data["plenti"] = []FileWithFreq{{"source1.txt", 2}}
	expected.Data["differ"] = []FileWithFreq{{"source1.txt", 1}, {"source2.txt", 1}}
	expected.Data["token"] = []FileWithFreq{{"source1.txt", 1}}
	expected.Data["good"] = []FileWithFreq{{"source1.txt", 1}}

	if err := in.AddToken("plenty", "source1.txt"); err != nil {
		t.Errorf("error occured\nexpected: %v", expected)
	}
	if err := in.AddToken("differ", "source2.txt"); err != nil {
		t.Errorf("error occured\nexpected: %v", expected)
	}
	if err := in.AddToken("good", "source1.txt"); err != nil {
		t.Errorf("error occured\nexpected: %v", expected)
	}

	if !reflect.DeepEqual(in.Data, expected.Data) {
		t.Errorf("\nactual:   %v \nexpected: %v", in, expected)
	}
}
