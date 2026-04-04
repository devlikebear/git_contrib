package analyzer

type AuthorMetrics struct {
	Name             string   `json:"name"`
	TotalCommits     int      `json:"total_commits"`
	Additions        int      `json:"additions"`
	Deletions        int      `json:"deletions"`
	NetLOC           int      `json:"net_loc"`
	ActiveDays       int      `json:"active_days"`
	RepoCount        int      `json:"repo_count"`
	Repos            []string `json:"repos"`
	AvgCommitsPerDay float64  `json:"avg_commits_per_day"`
	AvgLOCPerDay     float64  `json:"avg_loc_per_day"`
}

type PeriodEntry struct {
	Commits   int `json:"commits"`
	Additions int `json:"additions"`
	Deletions int `json:"deletions"`
}

type AuthorPeriodData struct {
	Author  string                  `json:"author"`
	Entries map[string]*PeriodEntry `json:"entries"`
}

type RepoAuthorEntry struct {
	RepoName  string `json:"repo_name"`
	Author    string `json:"author"`
	Commits   int    `json:"commits"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
}

type RepoReport struct {
	Name            string             `json:"name"`
	TotalCommits    int                `json:"total_commits"`
	TotalAdditions  int                `json:"total_additions"`
	TotalDeletions  int                `json:"total_deletions"`
	TotalActiveDays int                `json:"total_active_days"`
	Authors         []AuthorMetrics    `json:"authors"`
	Daily           []AuthorPeriodData `json:"daily"`
	Weekly          []AuthorPeriodData `json:"weekly"`
	Monthly         []AuthorPeriodData `json:"monthly"`
	AllDays         []string           `json:"all_days"`
	AllWeeks        []string           `json:"all_weeks"`
	AllMonths       []string           `json:"all_months"`
}

type Report struct {
	Since           string                          `json:"since"`
	Until           string                          `json:"until"`
	TotalCommits    int                             `json:"total_commits"`
	TotalAdditions  int                             `json:"total_additions"`
	TotalDeletions  int                             `json:"total_deletions"`
	TotalActiveDays int                             `json:"total_active_days"`
	TotalRepos      int                             `json:"total_repos"`
	Authors         []AuthorMetrics                 `json:"authors"`
	Daily           []AuthorPeriodData              `json:"daily"`
	Weekly          []AuthorPeriodData              `json:"weekly"`
	Monthly         []AuthorPeriodData              `json:"monthly"`
	AllDays         []string                        `json:"all_days"`
	AllWeeks        []string                        `json:"all_weeks"`
	AllMonths       []string                        `json:"all_months"`
	RepoBreakdown   map[string][]RepoAuthorEntry    `json:"repo_breakdown"`
	RepoReports     []RepoReport                    `json:"repo_reports"`
	Palette         []string                        `json:"-"`
	GeneratedAt     string                          `json:"generated_at"`
}
