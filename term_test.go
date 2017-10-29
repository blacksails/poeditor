package poeditor_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/blacksails/poeditor"
)

func TestTranslationUnmarshalJSON(t *testing.T) {
	unmarshalTests := []poeditor.Translation{
		{
			Content: "test",
		},
		{
			Content: poeditor.Plural{
				One:   "1 test",
				Other: "2 tests",
			},
		},
	}

	for _, t1 := range unmarshalTests {
		tJSON, err := json.Marshal(t1)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		var t2 poeditor.Translation
		if err := json.Unmarshal(tJSON, &t2); err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if !reflect.DeepEqual(t1, t2) {
			t.Errorf("\nExpected %+v \nGot      %+v", t1, t2)
		}
	}
}
