package lib

import (
	"regexp"
	"strconv"
)

var regWordSep = regexp.MustCompile(`\d+`)

// AutoPattern auto generate a regexp to match the index of sequential file names.
// Algorithm:
// 1. Create a table, the row is each path, the column is the words of each path.
// 2. Create histogram of the words in each column.
// 2. Analyze the histogram, use similar columns to create the regexp.
func AutoPattern(list []string) *regexp.Regexp {
	table := [][]string{}

	for _, l := range list {
		table = append(table, Split(l))
	}

	histograms := Histograms(table)
	indexCol := FindIndexCol(histograms)
	pattern := ""

	for col, histogram := range histograms {
		if col == indexCol {
			pattern += `(\d+)`
		} else {
			word, _ := FindWordForCol(histogram)
			pattern += regexp.QuoteMeta(word)
		}
	}

	return regexp.MustCompile(pattern)
}

// FindIndexCol find the column that looks like the index part of the paths.
// The column that has the most different numbers.
func FindIndexCol(list []map[string]int) int {
	max := 0
	col := 0
	for i, histogram := range list {
		// filter all the numbers
		nums := []int{}
		for word := range histogram {
			num, err := strconv.ParseInt(word, 10, 64)
			if err == nil {
				nums = append(nums, int(num))
			}
		}

		l := len(nums)
		if max < l {
			max = l
			col = i
		}

	}
	return col
}

// FindWordForCol the word that appear the most times
func FindWordForCol(histogram map[string]int) (string, int) {
	max := 0
	res := ""
	for word, count := range histogram {
		if max < count {
			max = count
			res = word
		}
	}
	return res, max
}

// Histograms histogram of the words of each path
func Histograms(table [][]string) []map[string]int {
	list := []map[string]int{}

	for col := 0; ; col++ {
		insufficient := true
		histogram := map[string]int{}

		for _, row := range table {
			if len(row) <= col {
				continue
			}
			insufficient = false
			histogram[row[col]]++
		}

		if len(histogram) > 0 {
			list = append(list, histogram)
		}

		if insufficient {
			return list
		}
	}
}

// Split use `\d+` as separator to split a path
func Split(path string) []string {
	locs := regWordSep.FindAllStringIndex(path, -1)

	if len(locs) == 0 {
		return nil
	}

	list := []string{}

	preLoc := []int{0, 0}
	var loc []int
	for _, loc = range locs {
		sepLeftSide := path[preLoc[1]:loc[0]]
		sep := path[loc[0]:loc[1]]
		list = append(list, sepLeftSide, sep)
		preLoc = loc
	}

	if loc[0] != loc[1] {
		list = append(list, path[loc[1]:])
	}

	return list
}
