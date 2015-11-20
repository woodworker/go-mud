package game

import (
	"testing"
)

var usernameTests = []struct {
	in  string
	out bool
}{
	{"a", true},
	{"aa", true},
	{"aWaa", true},
	{"aaa", true},
	{"aaa", true},
	{"asd", true},
	{"afsdfs", true},
	{"af-sdfs", true},
	{"afs_dfs", true},
	{"afs_-_2dfs", true},
	{"afs_2-dfs", true},
	{"afsdf324s", true},
	{"$#%3", false},
	{"../dsf", false},
	{"", false},
	{"päter", false},
}

func TestLevelIsValidUsername(t *testing.T) {
	s := Server{}
	for _, tt := range usernameTests {
		s := s.IsValidUsername(tt.in)
		if s != tt.out {
			t.Errorf("tests for username %q failed, should be %v", tt.in, tt.out )
		}
	}
}