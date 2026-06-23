// Package publish writes the verified proxy set to every output format and tracks
// per-proxy uptime across runs.
package publish

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/zakky8/tg-proxy-list/internal/model"
)

// Dataset is the JSON document consumed by the website and any API client.
type Dataset struct {
	GeneratedAtUTC string         `json:"generated_at_utc"`
	Count          int            `json:"count"`
	Countries      map[string]int `json:"countries"`
	Proxies        []model.Proxy  `json:"proxies"`
}

// Write emits all output files. rootDir holds the canonical lists; docsDir gets a
// copy of proxies.json for the static site.
func Write(rootDir, docsDir string, proxies []model.Proxy, generatedAt string) error {
	sorted := model.SortByLatency(proxies)

	countries := map[string]int{}
	for _, p := range sorted {
		countries[p.Country]++
	}

	ds := Dataset{
		GeneratedAtUTC: generatedAt,
		Count:          len(sorted),
		Countries:      countries,
		Proxies:        sorted,
	}

	// proxies.json (root + docs copy)
	if err := writeJSON(filepath.Join(rootDir, "proxies.json"), ds); err != nil {
		return err
	}
	if err := os.MkdirAll(docsDir, 0o755); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(docsDir, "proxies.json"), ds); err != nil {
		return err
	}

	// all_proxies.txt — reference-compatible flat list, fastest first.
	var b strings.Builder
	for _, p := range sorted {
		b.WriteString(p.HTTPSLink())
		b.WriteByte('\n')
	}
	if err := os.WriteFile(filepath.Join(rootDir, "all_proxies.txt"), []byte(b.String()), 0o644); err != nil {
		return err
	}

	// sorted_by_latency.txt — human-readable table.
	b.Reset()
	b.WriteString("# latency_ms  country  status        link\n")
	for _, p := range sorted {
		fmt.Fprintf(&b, "%-11d %-7s %-13s %s\n", p.LatencyMS, p.Country, p.Status, p.HTTPSLink())
	}
	if err := os.WriteFile(filepath.Join(rootDir, "sorted_by_latency.txt"), []byte(b.String()), 0o644); err != nil {
		return err
	}

	// by_country/<CC>.txt
	if err := writeByCountry(filepath.Join(rootDir, "by_country"), sorted); err != nil {
		return err
	}
	return nil
}

func writeByCountry(dir string, sorted []model.Proxy) error {
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	groups := map[string][]model.Proxy{}
	for _, p := range sorted {
		cc := p.Country
		if cc == "" || cc == "??" {
			cc = "XX"
		}
		groups[cc] = append(groups[cc], p)
	}
	codes := make([]string, 0, len(groups))
	for cc := range groups {
		codes = append(codes, cc)
	}
	sort.Strings(codes)
	for _, cc := range codes {
		var b strings.Builder
		for _, p := range groups[cc] {
			b.WriteString(p.HTTPSLink())
			b.WriteByte('\n')
		}
		if err := os.WriteFile(filepath.Join(dir, cc+".txt"), []byte(b.String()), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}
