package main

import (
	"fmt"
	"io"
	"os"

	"github.com/ryantate13/space-table/args"
	"github.com/ryantate13/space-table/internal/version"
	"github.com/ryantate13/space-table/output"
	"github.com/ryantate13/space-table/rows"
)

func main() {
	opts, err := args.Parse(os.Args[1:])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n%s\n", err.Error(), args.Help())
	}
	if opts.Help {
		fmt.Println(args.Help())
		return
	}
	if opts.Version {
		fmt.Println(version.Get())
		return
	}
	in, err := io.ReadAll(opts.Input)
	if err != nil {
		panic(err)
	}
	table, err := rows.Of(string(in))
	if err != nil {
		panic(err)
	}
	var e output.Encoder
	switch opts.Output {
	case output.CSV:
		e = output.CSVEncoder
	case output.JSON:
		e = output.JSONEncoder
	case output.Markdown:
		e = output.MarkdownEncoder
	case output.YAML:
		e = output.YAMLEncoder
	}
	fmt.Println(e.Encode(table))
}
