package poeditor_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/blacksails/poeditor"
)

func TestTermUnmarshalJSON(t *testing.T) {
	unmarshalTests := []poeditor.Term{
		{
			Term:        "test",
			Translation: poeditor.Singular("test"),
		},
		{
			Translation: poeditor.Plural{
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
		var t2 poeditor.Term
		if err := json.Unmarshal(tJSON, &t2); err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if !reflect.DeepEqual(t1, t2) {
			_, ok := t2.Translation.(poeditor.Singular)
			fmt.Println(ok)
			fmt.Println(reflect.TypeOf(t1.Translation), reflect.TypeOf(t2.Translation))
			t.Errorf("\nExpected %+v \nGot      %+v", t1, t2)
		}
	}
}
