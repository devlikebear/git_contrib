package git

import "time"

type Commit struct {
	Hash      string
	Author    string
	Email     string
	Subject   string
	RepoName  string
	Date      time.Time
	Additions int
	Deletions int
}
