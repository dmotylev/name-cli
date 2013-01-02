package main

import (
	"flag"
	"fmt"
	"github.com/dmotylev/goproperties"
	"github.com/dmotylev/name-cli/api"
	"os"
	"strings"
)

var (
	rcFile = flag.String("c",
		os.Getenv("HOME")+string(os.PathSeparator)+".name-cli.rc",
		"full path to rc file")
)

const (
	tsLayout   = "2006-01-02 15:04:05"
	lineFormat = "%-20s  %-19s  %-19s\n"
)

func main() {
	var (
		err error
		d   map[string]api.Domain
	)
	flag.Parse()
	rc, _ := properties.Load(*rcFile)
	ep := api.NewEndPoint(rc.GetString("api.endpoint", "https://api.name.com"))
	err = ep.Login(rc.GetString("api.username", ""), rc.GetString("api.token", ""))
	if err != nil {
		goto Error
	}
	d, err = ep.ListDomains()
	if err != nil {
		goto Error
	}

	fmt.Printf(lineFormat, "Domain", "Created", "Expired")
	fmt.Printf(lineFormat, strings.Repeat("-", 20), strings.Repeat("-", 19), strings.Repeat("-", 19))
	for k, v := range d {
		fmt.Printf(lineFormat, k, v.Created.Format(tsLayout), v.Expired.Format(tsLayout))
	}

	ep.Logout()
	return

Error:
	fmt.Fprintf(os.Stderr, "error: %s\n", err)
}
