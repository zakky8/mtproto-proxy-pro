package geo

import (
	"net"
	"strings"
	"testing"
)

func TestLookup(t *testing.T) {
	csv := "1.0.0.0,1.0.0.255,AU\n1.0.1.0,1.0.3.255,CN\n8.8.8.0,8.8.8.255,US\n"
	db, err := parse(strings.NewReader(csv))
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		ip   string
		want string
	}{
		{"1.0.0.50", "AU"},
		{"1.0.2.7", "CN"},
		{"8.8.8.8", "US"},
		{"9.9.9.9", "??"}, // gap
		{"::1", "??"},     // ipv6
	}
	for _, c := range cases {
		if got := db.LookupIP(net.ParseIP(c.ip)); got != c.want {
			t.Errorf("LookupIP(%s) = %q, want %q", c.ip, got, c.want)
		}
	}
}

func TestEmptyDBSafe(t *testing.T) {
	db := &DB{}
	if got := db.LookupIP(net.ParseIP("1.2.3.4")); got != "??" {
		t.Errorf("empty DB lookup = %q, want ??", got)
	}
}
