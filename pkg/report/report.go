package report

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/devlikebear/git-contrib/pkg/analyzer"
)

//go:embed templates/report.html
var templateFS embed.FS

func Generate(outputPath string, data *analyzer.Report) error {
	funcMap := template.FuncMap{
		"pct": func(val, max int) float64 {
			if max == 0 {
				return 0
			}
			return float64(val) / float64(max) * 100
		},
		"formatNum": func(n int) string {
			if n < 0 {
				return "-" + formatPositive(-n)
			}
			return formatPositive(n)
		},
		"signedNum": func(n int) string {
			if n >= 0 {
				return "+" + formatPositive(n)
			}
			return "-" + formatPositive(-n)
		},
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"rankBadge": func(i int) string {
			switch i {
			case 0:
				return "🥇"
			case 1:
				return "🥈"
			case 2:
				return "🥉"
			default:
				return fmt.Sprintf("#%d", i+1)
			}
		},
		"paletteColor": func(i int, palette []string) string {
			if len(palette) == 0 {
				return "#888"
			}
			return palette[i%len(palette)]
		},
		"fmtFloat1": func(f float64) string {
			return fmt.Sprintf("%.1f", f)
		},
		"fmtFloat0": func(f float64) string {
			return fmt.Sprintf("%.0f", f)
		},
		"maxCommits": func(authors []analyzer.AuthorMetrics) int {
			max := 0
			for _, a := range authors {
				if a.TotalCommits > max {
					max = a.TotalCommits
				}
			}
			return max
		},
		"getPeriodEntry": func(apd analyzer.AuthorPeriodData, key string) *analyzer.PeriodEntry {
			if e, ok := apd.Entries[key]; ok {
				return e
			}
			return nil
		},
		"sortedRepoNames": func(m map[string][]analyzer.RepoAuthorEntry) []string {
			keys := make([]string, 0, len(m))
			for k := range m {
				keys = append(keys, k)
			}
			for i := 0; i < len(keys); i++ {
				for j := i + 1; j < len(keys); j++ {
					if keys[i] > keys[j] {
						keys[i], keys[j] = keys[j], keys[i]
					}
				}
			}
			return keys
		},
		"topNAuthors": func(authors []analyzer.AuthorMetrics, n int) []analyzer.AuthorMetrics {
			if len(authors) <= n {
				return authors
			}
			return authors[:n]
		},
		"topNPeriod": func(data []analyzer.AuthorPeriodData, n int) []analyzer.AuthorPeriodData {
			if len(data) <= n {
				return data
			}
			return data[:n]
		},
		"seq": func(n int) []int {
			s := make([]int, n)
			for i := range s {
				s[i] = i
			}
			return s
		},
		"slugify": func(s string) string {
			r := strings.NewReplacer(" ", "_", "/", "_", ".", "_", "-", "_")
			return r.Replace(strings.ToLower(s))
		},
		"timeseriesJSON": func(days []string, periodData []analyzer.AuthorPeriodData, palette []string, n int) template.HTMLAttr {
			type authorSeries struct {
				Name    string `json:"name"`
				Color   string `json:"color"`
				Commits []int  `json:"commits"`
				LOC     []int  `json:"loc"`
			}
			type tsData struct {
				Days    []string       `json:"days"`
				Authors []authorSeries `json:"authors"`
			}
			top := periodData
			if len(top) > n {
				top = top[:n]
			}
			d := tsData{Days: days}
			if d.Days == nil {
				d.Days = []string{}
			}
			for i, apd := range top {
				s := authorSeries{
					Name:    apd.Author,
					Color:   palette[i%len(palette)],
					Commits: make([]int, len(days)),
					LOC:     make([]int, len(days)),
				}
				for j, day := range days {
					if e, ok := apd.Entries[day]; ok {
						s.Commits[j] = e.Commits
						s.LOC[j] = e.Additions + e.Deletions
					}
				}
				d.Authors = append(d.Authors, s)
			}
			if d.Authors == nil {
				d.Authors = []authorSeries{}
			}
			b, _ := json.Marshal(d)
			return template.HTMLAttr(b)
		},
	}

	tmpl, err := template.New("report.html").Funcs(funcMap).ParseFS(templateFS, "templates/report.html")
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}

func formatPositive(n int) string {
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}
	var parts []string
	for i := len(s); i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		parts = append([]string{s[start:i]}, parts...)
	}
	return strings.Join(parts, ",")
}
