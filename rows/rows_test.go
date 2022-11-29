package rows

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Error(t *testing.T) {
	tests := []struct {
		it   string
		err  error
		want string
	}{
		{
			it:   "returns error string",
			err:  Error("hello"),
			want: "hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.it, func(t *testing.T) {
			require.Equal(t, tt.want, tt.err.Error())
		})
	}
}

func Test_offsets(t *testing.T) {
	tests := []struct {
		it     string
		line   string
		assert func(t *testing.T, got [][]int)
	}{
		{
			it:   "returns single column indices without spaces",
			line: "hello",
			assert: func(t *testing.T, got [][]int) {
				require.Equal(t, [][]int{{0, 5}}, got)
			},
		},
		{
			it:   "returns single column indices with spaces",
			line: "hello     ",
			assert: func(t *testing.T, got [][]int) {
				require.Equal(t, [][]int{{0, 10}}, got)
			},
		},
		{
			it:   "returns multi-column indices with single space",
			line: "hello world",
			assert: func(t *testing.T, got [][]int) {
				require.Equal(t, [][]int{{0, 6}, {6, 11}}, got)
			},
		},
		{
			it:   "returns multi-column indices with multiple spaces",
			line: "hello   world   ",
			assert: func(t *testing.T, got [][]int) {
				require.Equal(t, [][]int{{0, 8}, {8, 16}}, got)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.it, func(t *testing.T) {
			tt.assert(t, offsets(tt.line))
		})
	}
}

func Test_slice(t *testing.T) {
	tests := []struct {
		it   string
		s    string
		idx  []int
		want string
	}{
		{
			it:   "returns an empty string if start index is greater than string length",
			s:    "123",
			idx:  []int{4, 100},
			want: "",
		},
		{
			it:   "returns from start index to end of string if end index is greater than string length",
			s:    "123",
			idx:  []int{0, 100},
			want: "123",
		},
		{
			it:   "returns from start index to end index if both are within string length",
			s:    "123",
			idx:  []int{0, 2},
			want: "12",
		},
	}
	for _, tt := range tests {
		t.Run(tt.it, func(t *testing.T) {
			require.Equal(t, tt.want, slice(tt.s, tt.idx))
		})
	}
}

func Test_isBlank(t *testing.T) {
	tests := []struct {
		it   string
		s    string
		want bool
	}{
		{
			it:   "returns true for empty string",
			s:    "",
			want: true,
		},
		{
			it:   "returns true for string with only spaces",
			s:    "   ",
			want: true,
		},
		{
			it:   "returns true for string with only white space characters",
			s:    "\t\n\r   \t\n\r",
			want: true,
		},
		{
			it:   "returns false for string with non-space characters",
			s:    "hello",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.it, func(t *testing.T) {
			if tt.want {
				require.True(t, isBlank(tt.s))
			} else {
				require.False(t, isBlank(tt.s))
			}
		})
	}
}

func TestOf(t *testing.T) {
	tests := []struct {
		it     string
		table  string
		assert func(t *testing.T, got [][]string, err error)
	}{
		{
			it:    "returns an error when no rows are found",
			table: "\n",
			assert: func(t *testing.T, got [][]string, err error) {
				require.Nil(t, got)
				require.Equal(t, ErrNoRows, err)
			},
		},
		{
			it:    "returns the rows of a text table",
			table: "HELLO WORLD\nfoo   bar",
			assert: func(t *testing.T, got [][]string, err error) {
				require.NoError(t, err)
				require.Equal(t, [][]string{
					{"HELLO", "WORLD"},
					{"foo", "bar"},
				}, got)
			},
		},
		{
			it: "strips indentation",
			table: `
				HELLO WORLD
				foo   bar
				1     2
			`,
			assert: func(t *testing.T, got [][]string, err error) {
				require.NoError(t, err)
				require.Equal(t, [][]string{
					{"HELLO", "WORLD"},
					{"foo", "bar"},
					{"1", "2"},
				}, got)
			},
		},
		{
			it: "keeps sparse rows in the center of a table",
			table: `
				HELLO WORLD

				1     2
			`,
			assert: func(t *testing.T, got [][]string, err error) {
				require.NoError(t, err)
				require.Equal(t, [][]string{
					{"HELLO", "WORLD"},
					{"", ""},
					{"1", "2"},
				}, got)
			},
		},
		{
			it: "does not include empty rows at the end of a table",
			table: `
				HELLO WORLD




			`,
			assert: func(t *testing.T, got [][]string, err error) {
				require.NoError(t, err)
				require.Equal(t, [][]string{{"HELLO", "WORLD"}}, got)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.it, func(t *testing.T) {
			got, err := Of(tt.table)
			tt.assert(t, got, err)
		})
	}
}

func TestToMap(t *testing.T) {
	tests := []struct {
		it    string
		table [][]string
		want  []map[string]string
	}{
		{
			it: "returns an array of maps with an entry for each row",
			table: [][]string{
				{"key", "value"},
				{"HELLO", "WORLD"},
				{"foo", "bar"},
			},
			want: []map[string]string{
				{"key": "HELLO", "value": "WORLD"},
				{"key": "foo", "value": "bar"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.it, func(t *testing.T) {
			require.Equal(t, tt.want, ToMap(tt.table))
		})
	}
}
