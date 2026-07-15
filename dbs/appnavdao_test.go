package dbs

import "testing"

func TestNextAppAliasNo(t *testing.T) {
	tests := []struct {
		name     string
		maxAlias int64
		want     string
	}{
		{name: "empty table", maxAlias: 0, want: "100000"},
		{name: "below start", maxAlias: 99998, want: "100000"},
		{name: "at start", maxAlias: 100000, want: "100001"},
		{name: "increments maximum", maxAlias: 123456, want: "123457"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := nextAppAliasNo(test.maxAlias); got != test.want {
				t.Fatalf("got %s, want %s", got, test.want)
			}
		})
	}
}
