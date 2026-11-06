// Package idcardarea looks up mainland Chinese ID address codes offline.
package idcardarea

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/jimyag/cn-idcard-area/internal/addresscode"
)

// Area contains the administrative divisions associated with an address code.
// A level that does not exist is represented by an empty string.
type Area struct {
	Province string `json:"province"`
	City     string `json:"city"`
	District string `json:"district"`
}

// Lookup finds the area for a six-digit address code or a 15- or 18-character
// ID number. It checks the input format but does not verify the birth date or
// checksum. It returns the latest record by start year, including retired
// codes under their historical names. Invalid or unknown inputs return the
// zero Area and false, as does unavailable area data. See DataError to check
// whether the embedded data loaded successfully. Lookup is safe for concurrent use.
func Lookup(input string) (Area, bool) {
	if len(input) != 6 {
		input = addresscode.Prefix(input)
	}
	record, ok := areaIndex[input]
	return record.area, ok
}

// DataError reports the error loading the embedded area data, or nil on success.
func DataError() error {
	return areaDataErr
}

//go:embed data/areas.json
var areaJSON []byte

var areaIndex, areaDataErr = loadIndex(areaJSON)

type areaNode struct {
	Code     int        `json:"code"`
	Name     string     `json:"name"`
	Start    int        `json:"start"`
	Children []areaNode `json:"children"`
}

type areaRecord struct {
	area  Area
	start int
}

func loadIndex(data []byte) (map[string]areaRecord, error) {
	var nodes []areaNode
	if err := json.Unmarshal(data, &nodes); err != nil {
		return nil, fmt.Errorf("idcardarea: decode embedded area data: %w", err)
	}
	if len(nodes) == 0 {
		return nil, errors.New("idcardarea: empty embedded area data")
	}
	index := make(map[string]areaRecord)
	if err := indexNodes(index, 0, Area{}, nodes); err != nil {
		return nil, err
	}
	return index, nil
}

func indexNodes(index map[string]areaRecord, parentCode int, parent Area, nodes []areaNode) error {
	for _, node := range nodes {
		code := node.Code
		if code < 100_000 || code > 999_999 || node.Name == "" || node.Start <= 0 {
			return fmt.Errorf("idcardarea: invalid embedded record for code %d", code)
		}
		provinceCode, cityCode := code/10_000*10_000, code/100*100
		area := parent
		switch {
		case code == provinceCode && parentCode == 0:
			area.Province = node.Name
		case code == cityCode && code != provinceCode && parentCode == provinceCode:
			area.City = node.Name
		case code != cityCode && (parentCode == cityCode || parentCode == provinceCode):
			// Counties directly under a province have no prefecture-level city.
			area.District = node.Name
		default:
			return fmt.Errorf("idcardarea: invalid parent %d for code %d", parentCode, code)
		}
		key := strconv.Itoa(code)
		if node.Start > index[key].start {
			index[key] = areaRecord{
				area:  area,
				start: node.Start,
			}
		}
		if err := indexNodes(index, code, area, node.Children); err != nil {
			return err
		}
	}
	return nil
}
