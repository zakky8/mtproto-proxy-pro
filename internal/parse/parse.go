// Package parse turns the many text formats public proxy lists use into model.Proxy records.
//
// Accepted shapes (one per line, surrounding junk tolerated):
//
//	tg://proxy?server=HOST&port=PORT&secret=SECRET
//	https://t.me/proxy?server=HOST&port=PORT&secret=SECRET
//	HOST:PORT:SECRET
//	HOST PORT SECRET            (whitespace separated)
package parse

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/zakky8/mtproto-proxy-pro/internal/model"
)

// linkRe extracts proxy links embedded anywhere in a line.
var linkRe = regexp.MustCompile(`(?i)(?:tg://proxy\?|https?://t\.me/proxy\?)[^\s"'<>]+`)

// tripletRe matches HOST:PORT:SECRET or HOST PORT SECRET anywhere in a line.
var tripletRe = regexp.MustCompile(`(?i)([a-z0-9.\-]+)[\s:]+(\d{2,5})[\s:]+([0-9a-f]{32,})`)

// Line parses a single line into zero or more proxies.
func Line(line string) []model.Proxy {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return nil
	}

	var out []model.Proxy
	seen := map[string]bool{}
	add := func(p model.Proxy, ok bool) {
		if !ok {
			return
		}
		p.Secret = strings.ToLower(strings.TrimSpace(p.Secret))
		// Normalize the host: Telegram is picky about trailing dots and casing in
		// deep links, and scraped lists often carry "Host.Example." verbatim.
		p.Server = strings.ToLower(strings.TrimRight(strings.TrimSpace(p.Server), "."))
		p.Type = model.ClassifySecret(p.Secret)
		if p.Validate() != nil {
			return
		}
		if k := p.Key(); !seen[k] {
			seen[k] = true
			out = append(out, p)
		}
	}

	for _, m := range linkRe.FindAllString(line, -1) {
		add(fromLink(m))
	}
	if len(out) == 0 {
		for _, m := range tripletRe.FindAllStringSubmatch(line, -1) {
			add(fromTriplet(m[1], m[2], m[3]))
		}
	}
	return out
}

// Text parses a whole document, deduping across all lines.
func Text(body string) []model.Proxy {
	var out []model.Proxy
	seen := map[string]bool{}
	for _, ln := range strings.Split(body, "\n") {
		for _, p := range Line(ln) {
			if k := p.Key(); !seen[k] {
				seen[k] = true
				out = append(out, p)
			}
		}
	}
	return out
}

func fromLink(link string) (model.Proxy, bool) {
	q := link
	if i := strings.Index(link, "?"); i >= 0 {
		q = link[i+1:]
	}
	v, err := url.ParseQuery(q)
	if err != nil {
		return model.Proxy{}, false
	}
	port, err := strconv.Atoi(v.Get("port"))
	if err != nil {
		return model.Proxy{}, false
	}
	return model.Proxy{
		Server: strings.TrimSpace(v.Get("server")),
		Port:   port,
		Secret: v.Get("secret"),
	}, true
}

func fromTriplet(host, port, secret string) (model.Proxy, bool) {
	p, err := strconv.Atoi(port)
	if err != nil {
		return model.Proxy{}, false
	}
	return model.Proxy{Server: host, Port: p, Secret: secret}, true
}
