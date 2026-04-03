package analyzer

import (
	"fmt"
	"sort"
	"time"

	"github.com/devlikebear/git-contrib/pkg/git"
)

var defaultPalette = []string{
	"#2563eb", "#dc2626", "#16a34a", "#ca8a04", "#9333ea",
	"#0891b2", "#e11d48", "#65a30d", "#c026d3", "#ea580c",
	"#4f46e5", "#059669", "#d97706", "#7c3aed", "#0d9488",
}

func Analyze(commits []git.Commit, since, until string) *Report {
	type authorAcc struct {
		commits    int
		additions  int
		deletions  int
		activeDays map[string]bool
		repos      map[string]bool
	}

	authors := make(map[string]*authorAcc)
	dailyMap := make(map[string]map[string]*PeriodEntry)   // author -> day -> entry
	weeklyMap := make(map[string]map[string]*PeriodEntry)  // author -> week -> entry
	monthlyMap := make(map[string]map[string]*PeriodEntry) // author -> month -> entry
	repoMap := make(map[string]map[string]*RepoAuthorEntry)

	allDays := make(map[string]bool)
	allWeeks := make(map[string]bool)
	allMonths := make(map[string]bool)

	for _, c := range commits {
		a := c.Author
		if _, ok := authors[a]; !ok {
			authors[a] = &authorAcc{
				activeDays: make(map[string]bool),
				repos:      make(map[string]bool),
			}
		}
		acc := authors[a]
		acc.commits++
		acc.additions += c.Additions
		acc.deletions += c.Deletions
		acc.activeDays[c.Date.Format("2006-01-02")] = true
		acc.repos[c.RepoName] = true

		dayKey := c.Date.Format("2006-01-02")
		year, week := c.Date.ISOWeek()
		weekKey := fmt.Sprintf("%d-W%02d", year, week)
		monthKey := c.Date.Format("2006-01")

		allDays[dayKey] = true
		allWeeks[weekKey] = true
		allMonths[monthKey] = true

		addToPeriod(dailyMap, a, dayKey, c)
		addToPeriod(weeklyMap, a, weekKey, c)
		addToPeriod(monthlyMap, a, monthKey, c)

		// repo breakdown
		if _, ok := repoMap[c.RepoName]; !ok {
			repoMap[c.RepoName] = make(map[string]*RepoAuthorEntry)
		}
		if _, ok := repoMap[c.RepoName][a]; !ok {
			repoMap[c.RepoName][a] = &RepoAuthorEntry{RepoName: c.RepoName, Author: a}
		}
		ra := repoMap[c.RepoName][a]
		ra.Commits++
		ra.Additions += c.Additions
		ra.Deletions += c.Deletions
	}

	// Build sorted author metrics
	var authorList []AuthorMetrics
	for name, acc := range authors {
		repos := sortedKeys(acc.repos)
		days := len(acc.activeDays)
		avgCommits := 0.0
		avgLOC := 0.0
		if days > 0 {
			avgCommits = float64(acc.commits) / float64(days)
			avgLOC = float64(acc.additions+acc.deletions) / float64(days)
		}
		authorList = append(authorList, AuthorMetrics{
			Name:             name,
			TotalCommits:     acc.commits,
			Additions:        acc.additions,
			Deletions:        acc.deletions,
			NetLOC:           acc.additions - acc.deletions,
			ActiveDays:       days,
			RepoCount:        len(acc.repos),
			Repos:            repos,
			AvgCommitsPerDay: avgCommits,
			AvgLOCPerDay:     avgLOC,
		})
	}
	sort.Slice(authorList, func(i, j int) bool {
		return authorList[i].TotalCommits > authorList[j].TotalCommits
	})

	// Build period data ordered by author rank
	daily := buildAuthorPeriodList(authorList, dailyMap)
	weekly := buildAuthorPeriodList(authorList, weeklyMap)
	monthly := buildAuthorPeriodList(authorList, monthlyMap)

	// Build repo breakdown
	repoBreakdown := make(map[string][]RepoAuthorEntry)
	for repo, authorMap := range repoMap {
		var entries []RepoAuthorEntry
		for _, e := range authorMap {
			entries = append(entries, *e)
		}
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Commits > entries[j].Commits
		})
		repoBreakdown[repo] = entries
	}

	totalAdd, totalDel := 0, 0
	for _, c := range commits {
		totalAdd += c.Additions
		totalDel += c.Deletions
	}

	return &Report{
		Since:           since,
		Until:           until,
		TotalCommits:    len(commits),
		TotalAdditions:  totalAdd,
		TotalDeletions:  totalDel,
		TotalActiveDays: len(allDays),
		TotalRepos:      len(repoBreakdown),
		Authors:         authorList,
		Daily:           daily,
		Weekly:          weekly,
		Monthly:         monthly,
		AllDays:         sortedStringSet(allDays),
		AllWeeks:        sortedStringSet(allWeeks),
		AllMonths:       sortedStringSet(allMonths),
		RepoBreakdown:   repoBreakdown,
		Palette:         defaultPalette,
		GeneratedAt:     time.Now().Format("2006-01-02 15:04"),
	}
}

func addToPeriod(m map[string]map[string]*PeriodEntry, author, key string, c git.Commit) {
	if _, ok := m[author]; !ok {
		m[author] = make(map[string]*PeriodEntry)
	}
	if _, ok := m[author][key]; !ok {
		m[author][key] = &PeriodEntry{}
	}
	e := m[author][key]
	e.Commits++
	e.Additions += c.Additions
	e.Deletions += c.Deletions
}

func buildAuthorPeriodList(authors []AuthorMetrics, m map[string]map[string]*PeriodEntry) []AuthorPeriodData {
	var result []AuthorPeriodData
	for _, a := range authors {
		entries := make(map[string]*PeriodEntry)
		if am, ok := m[a.Name]; ok {
			for k, v := range am {
				cp := *v
				entries[k] = &cp
			}
		}
		result = append(result, AuthorPeriodData{
			Author:  a.Name,
			Entries: entries,
		})
	}
	return result
}

func sortedKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortedStringSet(m map[string]bool) []string {
	return sortedKeys(m)
}
