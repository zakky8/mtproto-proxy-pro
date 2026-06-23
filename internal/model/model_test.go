package model

import "testing"

func TestFakeTLSDomain(t *testing.T) {
	// ee + 16 bytes (32 hex) + hex("example.com")
	secret := "ee" + "0123456789abcdef0123456789abcdef" + "6578616d706c652e636f6d"
	p := Proxy{Server: "x", Port: 443, Secret: secret}
	if got := p.FakeTLSDomain(); got != "example.com" {
		t.Fatalf("FakeTLSDomain = %q, want example.com", got)
	}
}

func TestFakeTLSDomainNonEE(t *testing.T) {
	p := Proxy{Secret: "0123456789abcdef0123456789abcdef"}
	if got := p.FakeTLSDomain(); got != "" {
		t.Fatalf("FakeTLSDomain = %q, want empty", got)
	}
}

func TestClassifySecret(t *testing.T) {
	cases := map[string]string{
		"0123456789abcdef0123456789abcdef": TypePlain,
		"dd0123456789abcdef":               TypeDD,
		"ee0123456789abcdef":               TypeEE,
	}
	for s, want := range cases {
		if got := ClassifySecret(s); got != want {
			t.Errorf("ClassifySecret(%q) = %q, want %q", s, got, want)
		}
	}
}

func TestSortByLatency(t *testing.T) {
	in := []Proxy{
		{Server: "c", LatencyMS: 300, Status: StatusReachable},
		{Server: "a", LatencyMS: 100, Status: StatusReachable},
		{Server: "b", LatencyMS: 100, Status: StatusHandshakeOK},
	}
	out := SortByLatency(in)
	if out[0].Server != "b" || out[1].Server != "a" || out[2].Server != "c" {
		t.Fatalf("unexpected order: %s %s %s", out[0].Server, out[1].Server, out[2].Server)
	}
}

func TestValidate(t *testing.T) {
	good := Proxy{Server: "h", Port: 443, Secret: "0123456789abcdef0123456789abcdef"}
	if err := good.Validate(); err != nil {
		t.Fatalf("good proxy rejected: %v", err)
	}
	bad := Proxy{Server: "", Port: 443, Secret: "0123456789abcdef0123456789abcdef"}
	if bad.Validate() == nil {
		t.Fatal("empty server accepted")
	}
}
