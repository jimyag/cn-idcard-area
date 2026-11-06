package idcardarea_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/jimyag/cn-idcard-area"
)

func TestLookup(t *testing.T) {
	for prefix, want := range map[string]idcardarea.Area{
		"440305": {
			Province: "广东省",
			City:     "深圳市",
			District: "南山区",
		},
		"110103": {
			Province: "北京市",
			District: "崇文区",
		},
		"130103": {
			Province: "河北省",
			City:     "石家庄市",
			District: "桥东区",
		},
		"120110": {
			Province: "天津市",
			District: "东丽区",
		},
		"419001": {
			Province: "河南省",
			District: "济源市",
		},
		"441900": {
			Province: "广东省",
			City:     "东莞市",
		},
		"440000": {Province: "广东省"},
	} {
		prefix, want := prefix, want
		t.Run(prefix, func(t *testing.T) {
			t.Parallel()
			if got, ok := idcardarea.Lookup(prefix); !ok || got != want {
				t.Fatalf("Lookup(%q) = (%+v, %t), want (%+v, true)", prefix, got, ok, want)
			}
		})
	}
}

func TestLookupMissing(t *testing.T) {
	for _, prefix := range []string{"999999", "000000", "", "44030", "44030x", "４４", "111111111111111111"} {
		if got, ok := idcardarea.Lookup(prefix); ok || got != (idcardarea.Area{}) {
			t.Errorf("Lookup(%q) = (%+v, %t), want zero Area and false", prefix, got, ok)
		}
	}
}

func ExampleLookup() {
	area, ok := idcardarea.Lookup("110103")
	if !ok {
		return
	}
	if err := json.NewEncoder(os.Stdout).Encode(area); err != nil {
		panic(err)
	}
	// Output: {"province":"北京市","city":"","district":"崇文区"}
}
