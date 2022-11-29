package args

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ryantate13/space-table/output"
)

func TestHelp(t *testing.T) {
	t.Run("it returns help message for all args", func(t *testing.T) {
		msg := Help()
		require.True(t, strings.Contains(msg, "Options:"))
		require.True(t, strings.Contains(msg, "Examples:"))
		e := reflect.TypeOf(&Args{})
		for i := 0; i < e.Elem().NumField(); i++ {
			f := e.Elem().Field(i)
			n := strings.ToLower(f.Name)
			require.True(t, strings.Contains(msg, "-"+n[0:1]+", --"+n))
			require.True(t, strings.Contains(msg, f.Tag.Get("help")))
		}
	})
}

func TestParse(t *testing.T) {
	tests := []struct {
		it     string
		argv   []string
		assert func(t *testing.T, args *Args, err error)
	}{
		{
			it:   "it returns default args when no args are passed",
			argv: []string{},
			assert: func(t *testing.T, args *Args, err error) {
				require.NoError(t, err)
				require.Equal(t, &Args{
					Input:  os.Stdin,
					Output: output.CSV,
				}, args)
			},
		},
		{
			it:   "correctly parses short args",
			argv: []string{"-h", "-v", "-i", "args_test.go", "-o", "json"},
			assert: func(t *testing.T, args *Args, err error) {
				require.NoError(t, err)
				f, err := os.Open("args_test.go")
				require.NoError(t, err)
				require.Equal(t, f.Name(), args.Input.Name())
				require.Equal(t, output.JSON, args.Output)
				require.True(t, args.Version)
				require.True(t, args.Help)
			},
		},
		{
			it:   "correctly parses long args",
			argv: []string{"--help", "--version", "--input", "args_test.go", "--output", "json"},
			assert: func(t *testing.T, args *Args, err error) {
				require.NoError(t, err)
				f, err := os.Open("args_test.go")
				require.NoError(t, err)
				require.Equal(t, f.Name(), args.Input.Name())
				require.Equal(t, output.JSON, args.Output)
				require.True(t, args.Version)
				require.True(t, args.Help)
			},
		},
		{
			it:   "correctly parses yaml format",
			argv: []string{"-o", "yaml"},
			assert: func(t *testing.T, args *Args, err error) {
				require.NoError(t, err)
				require.Equal(t, output.YAML, args.Output)
			},
		},
		{
			it:   "correctly parses yml format",
			argv: []string{"-o", "yml"},
			assert: func(t *testing.T, args *Args, err error) {
				require.NoError(t, err)
				require.Equal(t, output.YAML, args.Output)
			},
		},
		{
			it:   "correctly parses csv format",
			argv: []string{"-o", "csv"},
			assert: func(t *testing.T, args *Args, err error) {
				require.NoError(t, err)
				require.Equal(t, output.CSV, args.Output)
			},
		},
		{
			it:   "correctly parses md format",
			argv: []string{"-o", "md"},
			assert: func(t *testing.T, args *Args, err error) {
				require.NoError(t, err)
				require.Equal(t, output.Markdown, args.Output)
			},
		},
		{
			it:   "correctly parses markdown format",
			argv: []string{"-o", "markdown"},
			assert: func(t *testing.T, args *Args, err error) {
				require.NoError(t, err)
				require.Equal(t, output.Markdown, args.Output)
			},
		},
		{
			it:   "returns an error if output argument is missing",
			argv: []string{"--output"},
			assert: func(t *testing.T, args *Args, err error) {
				require.Error(t, err)
				require.Equal(t, "missing argument for --output", err.Error())
			},
		},
		{
			it:   "returns an error if output argument is invalid",
			argv: []string{"--output", "INVALID"},
			assert: func(t *testing.T, args *Args, err error) {
				require.Error(t, err)
				require.Equal(t, "invalid output format: INVALID", err.Error())
			},
		},
		{
			it:   "returns an error if input argument is missing",
			argv: []string{"--input"},
			assert: func(t *testing.T, args *Args, err error) {
				require.Error(t, err)
				require.Equal(t, "missing argument for --input", err.Error())
			},
		},
		{
			it:   "returns an error if input file cannot be opened",
			argv: []string{"--input", "INVALID"},
			assert: func(t *testing.T, args *Args, err error) {
				require.Error(t, err)
				require.True(t, os.IsNotExist(err))
			},
		},
		{
			it:   "sets input file to stdin when passed -",
			argv: []string{"--input", "-"},
			assert: func(t *testing.T, args *Args, err error) {
				require.NoError(t, err)
				require.Equal(t, os.Stdin, args.Input)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.it, func(t *testing.T) {
			args, err := Parse(tt.argv)
			tt.assert(t, args, err)
		})
	}
}
