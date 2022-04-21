package utils

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tables := []struct {
		a []string
		f Flags
	}{
		{[]string{"-t", "domain", "-p", "1", "google.pt"}, Flags{Type: "domain", Format: "text", Domain: "google.pt", Port: []int{1}}},
		{[]string{"--t", "domain", "-p", "1", "google.pt"}, Flags{Type: "domain", Format: "text", Domain: "google.pt", Port: []int{1}}},
		{[]string{"--t=domain", "-p", "1", "google.pt"}, Flags{Type: "domain", Format: "text", Domain: "google.pt", Port: []int{1}}},
		{[]string{"--t=domain", "-p", "1", "google.pt"}, Flags{Type: "domain", Format: "text", Domain: "google.pt", Port: []int{1}}},
		{[]string{"--t=domain", "-p", "1-5", "google.pt"}, Flags{Type: "domain", Format: "text", Domain: "google.pt", Port: []int{1, 2, 3, 4, 5}}},
		{[]string{"--t=domain", "-p", "1-5", "-w", "a.txt", "google.pt"}, Flags{Type: "domain", WordList: "a.txt", Format: "text", Domain: "google.pt", Port: []int{1, 2, 3, 4, 5}}},
		{[]string{"--t=domain", "-p", "1-5", "-w", "a.txt", "google.pt"}, Flags{Type: "domain", WordList: "a.txt", Format: "text", Domain: "google.pt", Port: []int{1, 2, 3, 4, 5}}},
		{[]string{"--t=domain", "-p", "1,3,4", "-w", "a.txt", "google.pt"}, Flags{Type: "domain", WordList: "a.txt", Format: "text", Domain: "google.pt", Port: []int{1, 3, 4}}},
		{[]string{"--t=domain", "-p", "1,3,4-7,20", "-w", "a.txt", "google.pt"}, Flags{Type: "domain", WordList: "a.txt", Format: "text", Domain: "google.pt", Port: []int{1, 3, 4, 5, 6, 7, 20}}},
	}

	for _, table := range tables {
		if flags := Parse(table.a); !reflect.DeepEqual(flags, table.f) {
			t.Errorf("Parse was incorrect, got: %v, want: %v.", flags, table.f)
		}
	}
}
