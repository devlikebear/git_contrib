package analyzer

type AuthorMetrics struct {
	Name             string
	TotalCommits     int
	Additions        int
	Deletions        int
	NetLOC           int
	ActiveDays       int
	RepoCount        int
	Repos            []string
	AvgCommitsPerDay float64
	AvgLOCPerDay     float64
}

type PeriodEntry struct {
	Commits   int
	Additions int
	Deletions int
}

type AuthorPeriodData struct {
	Author  string
	Entries map[string]*PeriodEntry
}

type RepoAuthorEntry struct {
	RepoName  string
	Author    string
	Commits   int
	Additions int
	Deletions int
}

type RepoReport struct {
	Name            string
	TotalCommits    int
	TotalAdditions  int
	TotalDeletions  int
	TotalActiveDays int
	Authors         []AuthorMetrics
	Daily           []AuthorPeriodData
	Weekly          []AuthorPeriodData
	Monthly         []AuthorPeriodData
	AllDays         []string
	AllWeeks        []string
	AllMonths       []string
}

type Report struct {
	Since           string
	Until           string
	TotalCommits    int
	TotalAdditions  int
	TotalDeletions  int
	TotalActiveDays int
	TotalRepos      int
	Authors         []AuthorMetrics
	Daily           []AuthorPeriodData
	Weekly          []AuthorPeriodData
	Monthly         []AuthorPeriodData
	AllDays         []string
	AllWeeks        []string
	AllMonths       []string
	RepoBreakdown   map[string][]RepoAuthorEntry
	RepoReports     []RepoReport
	Palette         []string
	GeneratedAt     string
}
