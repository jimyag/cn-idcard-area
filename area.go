// Package idcardarea looks up mainland Chinese ID address codes offline.
package idcardarea

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
)

// Area contains the administrative divisions associated with an address code.
// A level that does not exist is represented by an empty string.
type Area struct {
	Province string `json:"province"`
	City     string `json:"city"`
	District string `json:"district"`
}

// Lookup finds the area for a six-digit address code (the first six digits of
// an ID number). It returns the latest record by start year, including retired
// codes under their historical names. For an invalid or unknown code, it
// returns the zero Area and false. Lookup is safe for concurrent use.
func Lookup(prefix string) (Area, bool) {
	if len(prefix) != 6 {
		return Area{}, false
	}
	record, ok := areaIndex()[prefix]
	return record.area, ok
}

//go:embed data/areas.json
var areaJSON []byte

var areaIndex = sync.OnceValue(loadIndex)

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

func loadIndex() map[string]areaRecord {
	var nodes []areaNode
	if err := json.Unmarshal(areaJSON, &nodes); err != nil {
		panic(fmt.Errorf("idcardarea: decode embedded area data: %w", err))
	}
	if len(nodes) == 0 {
		panic("idcardarea: empty embedded area data")
	}
	index := make(map[string]areaRecord)
	indexNodes(index, 0, Area{}, nodes)
	return index
}

func indexNodes(index map[string]areaRecord, parentCode int, parent Area, nodes []areaNode) {
	for _, node := range nodes {
		code := node.Code
		if code < 100_000 || code > 999_999 || node.Name == "" || node.Start <= 0 {
			panic(fmt.Sprintf("idcardarea: invalid embedded record for code %d", code))
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
			panic(fmt.Sprintf("idcardarea: invalid parent %d for code %d", parentCode, code))
		}
		key := strconv.Itoa(code)
		if node.Start > index[key].start {
			index[key] = areaRecord{
				area:  area,
				start: node.Start,
			}
		}
		indexNodes(index, code, area, node.Children)
	}
}
