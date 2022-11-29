package rows

import (
	"regexp"
	"strings"
	"unicode"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrNoRows Error = "no rows found"
)

var (
	re = regexp.MustCompile(`\S+\s*`)
)

func offsets(line string) [][]int {
	return re.FindAllStringIndex(line, -1)
}

func slice(s string, i []int) string {
	if i[0] >= len(s) {
		return ""
	}
	if i[1] >= len(s) {
		return s[i[0]:]
	}
	return s[i[0]:i[1]]
}

func isBlank(s string) bool {
	for _, r := range s {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

// Of returns the rows and cells of tabular text with leading and trailing blank lines and indentation omitted
func Of(table string) ([][]string, error) {
	lines := strings.Split(table, "\n")
	var header int
	for header = 0; header < len(lines); header++ {
		if !isBlank(lines[header]) {
			break
		}
	}
	if header == len(lines) {
		return nil, ErrNoRows
	}
	var end int
	for end = len(lines) - 1; end > 0; end-- {
		if !isBlank(lines[end]) {
			break
		}
	}
	lines = lines[header : end+1]
	o := offsets(lines[0])
	rows := make([][]string, 0)
	for _, line := range lines {
		row := make([]string, len(o))
		for j, offset := range o {
			row[j] = strings.TrimSpace(slice(line, offset))
		}
		rows = append(rows, row)
	}
	return rows, nil
}

// ToMap returns a table as an array of maps, where the keys are the first row of the table
func ToMap(table [][]string) []map[string]string {
	keys := table[0]
	rows := make([]map[string]string, 0)
	for _, row := range table[1:] {
		m := make(map[string]string)
		for j, key := range keys {
			m[key] = row[j]
		}
		rows = append(rows, m)
	}
	return rows
}
