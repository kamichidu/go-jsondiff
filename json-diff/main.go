package main

import (
	"github.com/kamichidu/go-jsondiff/json-diff/cmd"
)

var (
	// these variables are filled by goreleaser
	version = "N/A"
	commit  = "N/A"
	date    = "N/A"
	builtBy = "N/A"
)

func main() {
	cmd.Execute(cmd.VersionDescriptor{
		Version: version,
		Commit:  commit,
		Date:    date,
		BuiltBy: builtBy,
	})
}
