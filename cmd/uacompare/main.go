// Command uacompare measures this parser against an external reference corpus.
//
// Two sources are supported:
//
//	matomo    github.com/matomo-org/device-detector (LGPL-3.0-or-later)
//	uap-core  github.com/ua-parser/uap-core (Apache-2.0)
//
// No project shares a taxonomy with this one: the references name every product
// they know, while this package exposes a small set of constants. Reference
// labels without a counterpart here are reported as out of scope rather than as
// errors, and deliberate differences are declared in a deviations file, so the
// mismatches that remain are the ones worth a decision.
//
// The fixtures are downloaded on demand and are not part of this repository.
// That matters for the Matomo corpus, whose data carries the LGPL: fetching it
// locally to measure against is use, not redistribution.
//
// Usage:
//
//	go run ./cmd/uacompare -fetch              # download the Matomo fixtures
//	go run ./cmd/uacompare                     # compare
//	go run ./cmd/uacompare -source uap-core -fetch
//	go run ./cmd/uacompare -source uap-core
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	sourceName := flag.String("source", "matomo", "reference corpus: matomo or uap-core")
	dir := flag.String("dir", "", "directory holding the fixtures (default testdata/<source>)")
	deviationPath := flag.String("deviations", filepath.Join("testdata", "deviations.tsv"), "file declaring intentional differences")
	fetch := flag.Bool("fetch", false, "download the fixtures into -dir and exit")
	show := flag.Int("show", 10, "number of example mismatches to print per dimension")
	strict := flag.Bool("strict", false, "exit with a non-zero status when unexplained mismatches remain")
	flag.Parse()

	source, err := lookupSource(*sourceName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if *dir == "" {
		*dir = filepath.Join("testdata", source.name)
	}

	if *fetch {
		if err := source.fetch(*dir); err != nil {
			fmt.Fprintln(os.Stderr, "fetch:", err)
			os.Exit(1)
		}
		fmt.Printf("Fixtures written to %s\n", *dir)
		return
	}

	cases, err := source.load(*dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load: %v\nRun with -fetch to download the fixtures.\n", err)
		os.Exit(1)
	}
	declared, err := loadDeviations(*deviationPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "deviations:", err)
		os.Exit(1)
	}

	rep := compare(cases, declared)
	fmt.Print(rep.String(source, len(cases), *show))
	if *strict && rep.unexplained() > 0 {
		os.Exit(1)
	}
}

type report struct {
	dimensions []*dimension
}

type dimension struct {
	name     string
	compared int
	hits     int
	known    int
	skipped  int
	examples []string
	causes   map[string]int
}

func (d *dimension) hit() {
	d.compared++
	d.hits++
}

func (d *dimension) miss(ua, want, got string) {
	d.compared++
	d.examples = append(d.examples, fmt.Sprintf("want %-24s got %-24s %s", want, got, truncate(ua, 92)))
	if d.causes == nil {
		d.causes = map[string]int{}
	}
	d.causes[fmt.Sprintf("%s → %s", want, got)]++
}

// topCauses returns the most frequent want/got pairs, which is usually enough to
// tell a systematic difference from a handful of odd user agents.
func (d *dimension) topCauses(n int) []string {
	pairs := make([]string, 0, len(d.causes))
	for pair := range d.causes {
		pairs = append(pairs, pair)
	}
	sort.Slice(pairs, func(i, j int) bool {
		if d.causes[pairs[i]] != d.causes[pairs[j]] {
			return d.causes[pairs[i]] > d.causes[pairs[j]]
		}
		return pairs[i] < pairs[j]
	})
	if n > 0 && len(pairs) > n {
		pairs = pairs[:n]
	}
	out := make([]string, len(pairs))
	for i, pair := range pairs {
		out[i] = fmt.Sprintf("%6d  %s", d.causes[pair], pair)
	}
	return out
}

func (d *dimension) misses() int { return d.compared - d.hits }

func (r *report) unexplained() int {
	total := 0
	for _, d := range r.dimensions {
		total += d.misses()
	}
	return total
}

func (r *report) String(source source, cases, show int) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s comparison — %d user agents\nreference: %s (%s)\n\n", source.name, cases, source.repo, source.license)
	fmt.Fprintf(&b, "%-16s %14s %8s   %-10s %s\n", "dimension", "hits/compared", "rate", "declared", "out of scope")
	for _, d := range r.dimensions {
		fmt.Fprintf(&b, "%-16s %7d/%-6d %8s   %-10d %d\n",
			d.name, d.hits, d.compared, percent(d.hits, d.compared), d.known, d.skipped)
	}

	if show == 0 {
		return b.String()
	}
	for _, d := range r.dimensions {
		if len(d.examples) == 0 {
			continue
		}
		fmt.Fprintf(&b, "\n%s: %d mismatches\n", d.name, d.misses())
		fmt.Fprintf(&b, "  most frequent differences:\n")
		for _, cause := range d.topCauses(show) {
			fmt.Fprintf(&b, "  %s\n", cause)
		}

		shown := d.examples
		if show > 0 && len(shown) > show {
			shown = shown[:show]
		}
		fmt.Fprintf(&b, "  examples:\n")
		for _, e := range shown {
			fmt.Fprintf(&b, "    %s\n", e)
		}
	}
	return b.String()
}

func percent(num, den int) string {
	if den == 0 {
		return "n/a"
	}
	return fmt.Sprintf("%.1f%%", float64(num)/float64(den)*100)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "…"
}
