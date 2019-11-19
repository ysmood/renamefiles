package lib_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ysmood/renamefiles/lib"
)

func TestSplit(t *testing.T) {
	assert.Equal(t, []string{"", "03", "-a-b-", "01", ".txt"}, lib.Split("03-a-b-01.txt"))
	assert.Equal(t, []string{"-", "03", "-a-b-", "01", ".txt"}, lib.Split("-03-a-b-01.txt"))
	assert.Equal(t, []string(nil), lib.Split("a.txt"))
}

func TestHistograms(t *testing.T) {
	assert.Equal(t,
		[]map[string]int{
			map[string]int{"a": 6},                            // col1
			map[string]int{"01": 1, "b": 3, "s": 1},           // col2
			map[string]int{"01": 2, "02": 1, "03": 1, "b": 1}, // col3
		},
		lib.Histograms([][]string{
			//       col1 col2 col3
			[]string{"a"},
			[]string{"a", "b", "01"},
			[]string{"a", "b", "02"},
			[]string{"a", "b", "03"},
			[]string{"a", "s", "01"},
			[]string{"a", "01", "b"},
		}),
	)
}

func TestFindIncrementalCol(t *testing.T) {
	assert.Equal(t, 1,
		lib.FindIndexCol([]map[string]int{
			map[string]int{"01": 6},                           // col1
			map[string]int{"02": 2, "01": 1, "03": 1, "b": 1}, // col2
			map[string]int{"01": 1, "b": 3, "s": 1},           // col3
		}),
	)
}

func TestAutoPatternFunc(t *testing.T) {
	assert.Equal(t, `01-a-b-(\d+)\.txt`, lib.AutoPattern([]string{
		"01-a-b-01.txt",
		"01-a-b-02.txt",
		"01-a-b-03.txt",
		"a",
		"b",
	}).String())
}
