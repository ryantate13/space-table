package output

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		it    string
		table [][]string
		enc   Encoder
		want  string
	}{
		{
			it: "encodes a table to csv",
			table: [][]string{
				{"hello", "world"},
				{"foo", "bar"},
			},
			enc:  CSVEncoder,
			want: "hello,world\nfoo,bar",
		},
		{
			it: "encodes a table to json with header row",
			table: [][]string{
				{"hello", "world"},
				{"foo", "bar"},
			},
			enc: JSONEncoder,
			want: func() string {
				j, _ := json.MarshalIndent([]map[string]string{{
					"hello": "foo",
					"world": "bar",
				}}, "", "  ")
				return string(j)
			}(),
		},
		{
			it:    "encodes a table to json with single row",
			table: [][]string{{"hello", "world"}},
			enc:   JSONEncoder,
			want:  `["hello","world"]`,
		},
		{
			it: "encodes a table to markdown",
			table: [][]string{
				{"hello", "world"},
				{"foo", "bar"},
			},
			enc:  MarkdownEncoder,
			want: "| hello | world |\n|-------|-------|\n| foo   | bar   |",
		},
		{
			it: "encodes a table to yaml with header row",
			table: [][]string{
				{"hello", "world"},
				{"foo", "bar"},
			},
			enc: YAMLEncoder,
			want: func() string {
				j, _ := yaml.Marshal([]map[string]string{{
					"hello": "foo",
					"world": "bar",
				}})
				return string(j[:len(j)-1])
			}(),
		},
		{
			it:    "encodes a table to yaml with single row",
			table: [][]string{{"hello", "world"}},
			enc:   YAMLEncoder,
			want:  strings.Join([]string{"- hello", "- world"}, "\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.it, func(t *testing.T) {
			require.Equal(t, tt.want, tt.enc.Encode(tt.table))
		})
	}
}
