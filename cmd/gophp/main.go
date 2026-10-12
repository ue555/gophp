package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"example.com/gophp/internal/php"
)

func run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("gophp", flag.ContinueOnError)
	flags.SetOutput(stderr)
	target := flags.String("php-version", string(php.PHP84), "target PHP version: 8.3, 8.4 or 8.5")
	flags.Usage = func() {
		fmt.Fprintln(stderr, "usage: gophp [--php-version=8.3|8.4|8.5] <file.php>")
	}
	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}
	if flags.NArg() != 1 {
		flags.Usage()
		return 2
	}
	version, err := php.ParseVersion(*target)
	if err != nil {
		fmt.Fprintf(stderr, "gophp: %v\n", err)
		return 2
	}
	filename, err := filepath.Abs(flags.Arg(0))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	source, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return php.RunWithVersion(string(source), filename, stdout, stderr, version)
}

func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }
