package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	glob   string
	maxlen int
	passed = true
)

type longline struct {
	lineno int
	text   string
	fpath  string
}

func (l longline) Show() {
	fmt.Fprintf(os.Stderr, "%s:%d line length is %d\n",
		l.fpath, l.lineno, len(l.text))
}

func exitIfErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	flag.StringVar(&glob, "g", "*.go,**/*.go",
		"File globs to test, comma-separated")
	flag.IntVar(&maxlen, "l", 80, "Maximum line length to allow")
	flag.Parse()

	var files []string

	for _, g := range strings.Split(glob, ",") {
		f, err := filepath.Glob(g)
		exitIfErr(err)
		files = append(files, f...)
	}

	for _, fpath := range files {
		file, err := os.Open(fpath)
		exitIfErr(err)
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var lineno int
		for scanner.Scan() {
			lineno++
			ln := scanner.Text()
			if len(ln) > maxlen {
				ll := longline{lineno, ln, fpath}
				ll.Show()
				passed = false
			}
		}
		exitIfErr(scanner.Err())
	}

	if !passed {
		exitIfErr(errors.New("line length check failed"))
	}
}
