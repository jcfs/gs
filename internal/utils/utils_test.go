package utils

import "testing"

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
