package main

import (
	_ "embed"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
)

//go:embed data/result.csv
var areaCSV string

const (
	codeColumn = iota
	provinceColumn
	cityColumn
	nameColumn
	levelColumn
	_ // status
	startYearColumn
)

func main() {
	if len(os.Args) != 2 || !validPrefix(os.Args[1]) {
		fmt.Fprintln(os.Stderr, "usage: cn-idcard-area <6 ASCII digits>")
		os.Exit(2)
	}

	path, err := lookup(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "read area data: %v\n", err)
		os.Exit(1)
	}
	if _, err := fmt.Println(path); err != nil {
		fmt.Fprintf(os.Stderr, "write result: %v\n", err)
		os.Exit(1)
	}
}

func validPrefix(prefix string) bool {
	return len(prefix) == 6 && !strings.ContainsFunc(prefix, func(r rune) bool {
		return r < '0' || r > '9'
	})
}

func lookup(prefix string) (string, error) {
	rows, err := csv.NewReader(strings.NewReader(strings.TrimPrefix(areaCSV, "\ufeff"))).ReadAll()
	if err != nil {
		return "", err
	}
	header := []string{"代码", "一级行政区", "二级行政区", "名称", "级别", "状态", "启用时间", "变更/弃用时间", "新代码"}
	if len(rows) < 2 || !slices.Equal(rows[0], header) {
		return "", errors.New("unexpected area data schema")
	}

	// The prefix contains no assignment year. Pick the latest record by its four-digit start year.
	var latest []string
	for _, row := range rows[1:] {
		if row[codeColumn] == prefix && (latest == nil || row[startYearColumn] > latest[startYearColumn]) {
			latest = row
		}
	}
	if latest == nil {
		return "未知", nil
	}

	parts := []string{latest[provinceColumn]}
	if city := latest[cityColumn]; city != "" && city != "直辖" {
		parts = append(parts, city)
	}
	if latest[levelColumn] == "县级" {
		parts = append(parts, latest[nameColumn])
	}
	return strings.Join(parts, " / "), nil
}
