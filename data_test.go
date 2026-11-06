package idcardarea

import "testing"

func TestEmbeddedData(t *testing.T) {
	if err := DataError(); err != nil {
		t.Fatal(err)
	}
}

func TestLoadIndexRejectsInvalidData(t *testing.T) {
	for name, data := range map[string]string{
		"malformed JSON": `[`,
		"empty table":    `[]`,
		"invalid record": `[{"code":110000,"name":"Beijing","start":1981},{"code":120000,"name":"Tianjin"}]`,
		"invalid parent": `[{"code":110000,"name":"Beijing","start":1981,"children":[{"code":120101,"name":"Heping","start":1981}]}]`,
	} {
		t.Run(name, func(t *testing.T) {
			index, err := loadIndex([]byte(data))
			if index != nil || err == nil {
				t.Fatalf("loadIndex returned %d records and error %v; want nil index and an error", len(index), err)
			}
		})
	}
}
