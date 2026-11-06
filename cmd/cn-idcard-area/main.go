package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/jimyag/cn-idcard-area"
	"github.com/jimyag/cn-idcard-area/internal/addresscode"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	var code string
	if len(args) == 1 {
		code = addresscode.Prefix(args[0])
	}
	if code == "" {
		fmt.Fprintln(stderr, "usage: cn-idcard-area <6-digit code or 15-/18-character ID number>")
		return 2
	}
	if err := idcardarea.DataError(); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	area, ok := idcardarea.Lookup(code)
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
