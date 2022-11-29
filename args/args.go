package args

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/ryantate13/space-table/output"
)

// Args struct for command line arguments
type Args struct {
	Input   *os.File      `help:"input file, defaults to stdin, pass - to set input to stdin explicitly"`
	Output  output.Format `help:"output format, defaults to csv, allowed values are: csv, json, md, markdown, yml, yaml"`
	Help    bool          `help:"show this help message and quit"`
	Version bool          `help:"show version and quit"`
}

// Help returns the help message for all command line arguments
func Help() string {
	msg := `Usage: space-table [options]...
Transforms space-aligned tabular data from stdin or a file to various formats.

Options:
`
	e := reflect.ValueOf(&Args{}).Elem()
	for i := 0; i < e.NumField(); i++ {
		f := e.Type().Field(i)
		n := strings.ToLower(f.Name)
		msg += fmt.Sprintf("  -%s, --%s\t%s\n", n[0:1], n, f.Tag.Get("help"))
	}

	msg += `
Examples:
  space-table < table.txt
  space-table --input table.txt
  space-table --output json < table.txt
  space-table --output csv --input table.txt
  kubectl get pods -o wide | space-table
  my-command-that-outputs-tabular-data | space-table -o json > table.json`

	return msg
}

// Parse parses command line arguments into an Args struct
func Parse(argv []string) (*Args, error) {
	args := &Args{
		Input:  os.Stdin,
		Output: output.CSV,
	}
	e := reflect.ValueOf(args).Elem()
	for i := 0; i < e.NumField(); i++ {
		f := e.Type().Field(i)
		n := strings.ToLower(f.Name)
		for j, a := range argv {
			if a == "--"+n || a == "-"+n[0:1] {
				switch f.Type.Kind() {
				case reflect.Bool:
					e.Field(i).SetBool(true)
				case reflect.TypeOf(output.Format(0)).Kind():
					if j+1 == len(argv) {
						return nil, errors.New("missing argument for --" + n)
					}
					switch argv[j+1] {
					case "csv":
						e.Field(i).SetUint(uint64(output.CSV))
					case "json":
						e.Field(i).SetUint(uint64(output.JSON))
					case "md", "markdown":
						e.Field(i).SetUint(uint64(output.Markdown))
					case "yml", "yaml":
						e.Field(i).SetUint(uint64(output.YAML))
					default:
						return nil, errors.New("invalid output format: " + argv[j+1])
					}
				case reflect.TypeOf((*os.File)(nil)).Kind():
					if j+1 == len(argv) {
						return nil, errors.New("missing argument for --" + n)
					}
					if argv[j+1] == "-" {
						e.Field(i).Set(reflect.ValueOf(os.Stdin))
					} else {
						f, err := os.Open(argv[j+1])
						if err != nil {
							return nil, err
						}
						e.Field(i).Set(reflect.ValueOf(f))
					}
				}
			}
		}
	}

	return args, nil
}
