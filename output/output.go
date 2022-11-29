package output

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"strings"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"

	"github.com/ryantate13/space-table/rows"
)

// Format enum for output types
type Format uint8

const (
	Markdown Format = iota
	JSON
	CSV
	YAML
)

var (
	// CSVEncoder encodes a table to CSV
	CSVEncoder Encoder = csvOut{}
	// JSONEncoder encodes a table to JSON
	JSONEncoder Encoder = jsonOut{}
	// MarkdownEncoder encodes a table to Markdown
	MarkdownEncoder Encoder = markdownOut{}
	// YAMLEncoder encodes a table to YAML
	YAMLEncoder Encoder = yamlOut{}
)

// Encoder is an interface for encoding a table to a string
type Encoder interface {
	Encode(table [][]string) string
}

type csvOut struct{}

func (c csvOut) Encode(table [][]string) string {
	b := bytes.NewBuffer(nil)
	w := csv.NewWriter(b)
	w.WriteAll(table)
	return strings.TrimSpace(b.String())
}

type jsonOut struct{}

func (j jsonOut) Encode(table [][]string) string {
	var b []byte
	if len(table) == 1 {
		b, _ = json.Marshal(table[0])
	} else {
		b, _ = json.MarshalIndent(rows.ToMap(table), "", "  ")
	}
	return string(b)
}

type markdownOut struct{}

func (m markdownOut) Encode(table [][]string) string {
	b := bytes.NewBuffer(nil)
	t := tablewriter.NewWriter(b)
	t.SetAutoFormatHeaders(false)
	t.SetHeader(table[0])
	t.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	t.SetCenterSeparator("|")
	if len(table) > 1 {
		t.AppendBulk(table[1:])
	}
	t.Render()
	return strings.TrimSpace(b.String())
}

type yamlOut struct{}

func (y yamlOut) Encode(table [][]string) string {
	var b []byte
	if len(table) == 1 {
		b, _ = yaml.Marshal(table[0])
	} else {
		m := rows.ToMap(table)
		b, _ = yaml.Marshal(m)
	}
	return strings.TrimSpace(string(b))
}
