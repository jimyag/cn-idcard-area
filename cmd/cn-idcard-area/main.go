package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/jimyag/cn-idcard-area"
)

var prefixPattern = regexp.MustCompile(`^[0-9]{6}$`)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) != 1 || !prefixPattern.MatchString(args[0]) {
		fmt.Fprintln(stderr, "usage: cn-idcard-area <6 ASCII digits>")
		return 2
	}

	area, ok := idcardarea.Lookup(args[0])
	var err error
	if ok {
		err = json.NewEncoder(stdout).Encode(area)
	} else {
		_, err = fmt.Fprintln(stdout, "未知")
	}
	if err != nil {
		fmt.Fprintf(stderr, "write result: %v\n", err)
		return 1
	}
	return 0
}
