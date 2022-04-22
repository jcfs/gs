package utils

import (
	"reflect"
	"testing"
)

func TestTrimPrefix(t *testing.T) {
	tables := []struct {
		g string
		w string
		e string
	}{
		{"-", "-", ""},
		{"-verbose", "-v", "erbose"},
		{"--verbose", "-", "verbose"},
		{"--verbose", "--", "verbose"},
		{"xxxaxxa", "x", "axxa"},
		{"xxxaxxa", "xx", "xaxxa"},
		{"xxxaxxa", "xxx", "axxa"},
	}

	for _, table := range tables {
		if trimmed := TrimPrefix(table.g, table.w); trimmed != table.e {
			t.Errorf("TrimPrefix was incorrect, got: %s, want: %s.", trimmed, table.e)
		}
	}

}

func TestChunks(t *testing.T) {
	tables := []struct {
		given    []int
		size     int
		expected [][]int
	}{
		{[]int{1, 2, 3}, 1, [][]int{{1}, {2}, {3}}},
		{[]int{1, 2, 3}, 2, [][]int{{1, 2}, {3}}},
		{[]int{1, 2, 3}, 3, [][]int{{1, 2, 3}}},
		{[]int{1, 2, 3}, 4, [][]int{{1, 2, 3}}},
		{[]int{1, 2, 3, 4}, 2, [][]int{{1, 2}, {3, 4}}},
	}

	for _, i := range tables {
		chunks := Chunks(i.given, i.size)
		if !reflect.DeepEqual(chunks, i.expected) {
			t.Errorf("Chunk result was incorrect, got: %v, want: %v.", chunks, i.expected)
		}
	}
}
