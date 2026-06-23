package parse

import "testing"

func TestLineFormats(t *testing.T) {
	cases := []struct {
		name string
		line string
		want int // expected parsed proxies
	}{
		{"tg link", "tg://proxy?server=1.2.3.4&port=443&secret=ee" + hex32() + "6578616d706c652e636f6d", 1},
		{"tme link", "https://t.me/proxy?server=host.example&port=8888&secret=" + hex32(), 1},
		{"triplet colon", "1.2.3.4:443:" + hex32(), 1},
		{"triplet space", "host.example 443 dd" + hex32(), 1},
		{"comment", "# just a comment", 0},
		{"junk", "this is not a proxy", 0},
		{"short secret", "1.2.3.4:443:abcd", 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Line(c.line)
			if len(got) != c.want {
				t.Fatalf("Line(%q) = %d proxies, want %d", c.line, len(got), c.want)
			}
		})
	}
}

func TestLineParsesFields(t *testing.T) {
	got := Line("https://t.me/proxy?server=host.example&port=8888&secret=" + hex32())
	if len(got) != 1 {
		t.Fatalf("got %d proxies", len(got))
	}
	p := got[0]
	if p.Server != "host.example" || p.Port != 8888 {
		t.Fatalf("bad fields: %+v", p)
	}
	if p.Type != "plain" {
		t.Fatalf("type = %q, want plain", p.Type)
	}
}

func TestTextDedupes(t *testing.T) {
	body := "1.2.3.4:443:" + hex32() + "\n" + "1.2.3.4:443:" + hex32() + "\n"
	if got := Text(body); len(got) != 1 {
		t.Fatalf("Text dedupe = %d, want 1", len(got))
	}
}

// hex32 returns a 32-char hex secret.
func hex32() string { return "0123456789abcdef0123456789abcdef" }
