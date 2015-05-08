// Command godefaultinstance generates at the package level a default instance of a type
// and functions that call methods on it.
//
// This command is designed to be called with go generate; run godefaultinstance -h for help.
//
// Check out the example at github.com/tcard/godefaultinstance/example.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tcard/godefaultinstance/godefaultinstance"
)

var genType = flag.String("type", "", "Name of the type to instantiate. Prefix it with a * to make a pointer.")
var instanceName = flag.String("name", "", "(Optional.) Name of the instance. If not set, it will be Default<Name of the type>.")
var exclude = flag.String("exclude", "", "(Optional.) Comma-separated list of methods not to wrap.")

func main() {
	flag.Parse()
	if *genType == "" {
		flag.Usage()
		os.Exit(1)
	}
	c, err := godefaultinstance.NewConfig(
		os.Getenv("GOPACKAGE"),
		*genType,
		*instanceName,
		strings.Split(*exclude, ","),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	c.RawArgs = strings.Join(os.Args[1:], " ")

	err = c.Generate(filepath.Dir(os.Getenv("GOFILE")))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
