package git

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

func CollectAll(repos []string, since, until string) ([]Commit, error) {
	g, _ := errgroup.WithContext(context.Background())
	results := make([][]Commit, len(repos))

	for i, repo := range repos {
		i, repo := i, repo
		g.Go(func() error {
			commits, err := collectRepo(repo, since, until)
			if err != nil {
				fmt.Fprintf(os.Stderr, "⚠️  건너뜀: %s (%v)\n", repo, err)
				return nil
			}
			results[i] = commits
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	var all []Commit
	for _, r := range results {
		all = append(all, r...)
	}
	return all, nil
}

func collectRepo(repoPath, since, until string) ([]Commit, error) {
	repoName := getRepoName(repoPath)
	fmt.Fprintf(os.Stderr, "📊 분석 중: %s (%s)\n", repoName, repoPath)

	commits, err := parseCommitLog(repoPath, since, until)
	if err != nil {
		return nil, err
	}

	stats, err := parseNumstat(repoPath, since, until)
	if err != nil {
		return nil, err
	}

	for i := range commits {
		commits[i].RepoName = repoName
		if s, ok := stats[commits[i].Hash]; ok {
			commits[i].Additions = s.additions
			commits[i].Deletions = s.deletions
		}
	}

	return commits, nil
}

func getRepoName(repoPath string) string {
	out, err := runGit(repoPath, "remote", "get-url", "origin")
	if err == nil {
		url := strings.TrimSpace(string(out))
		if url != "" {
			name := filepath.Base(url)
			name = strings.TrimSuffix(name, ".git")
			return name
		}
	}
	return filepath.Base(filepath.Clean(repoPath))
}

func runGit(repoPath string, args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = repoPath
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("git %s: %s", strings.Join(args, " "), stderr.String())
	}
	return stdout.Bytes(), nil
}

type numstatEntry struct {
	additions int
	deletions int
}

func parseCommitLog(repoPath, since, until string) ([]Commit, error) {
	out, err := runGit(repoPath,
		"log",
		"--format=%H|%aN|%aE|%ad|%s",
		"--date=iso-strict",
		fmt.Sprintf("--since=%s", since),
		fmt.Sprintf("--until=%s", until),
		"--all",
	)
	if err != nil {
		return nil, err
	}

	var commits []Commit
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 5)
		if len(parts) < 5 {
			continue
		}
		dt, err := time.Parse(time.RFC3339, parts[3])
		if err != nil {
			continue
		}
		commits = append(commits, Commit{
			Hash:    parts[0],
			Author:  strings.TrimSpace(parts[1]),
			Email:   strings.TrimSpace(parts[2]),
			Date:    dt,
			Subject: strings.TrimSpace(parts[4]),
		})
	}
	return commits, nil
}

func parseNumstat(repoPath, since, until string) (map[string]numstatEntry, error) {
	out, err := runGit(repoPath,
		"log",
		"--format=%H|%aN",
		"--numstat",
		fmt.Sprintf("--since=%s", since),
		fmt.Sprintf("--until=%s", until),
		"--all",
	)
	if err != nil {
		return nil, err
	}

	stats := make(map[string]numstatEntry)
	var currentHash string

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		if strings.Contains(trimmed, "|") && (len(trimmed) == 0 || !isDigitOrDash(trimmed[0])) {
			parts := strings.SplitN(trimmed, "|", 2)
			currentHash = strings.TrimSpace(parts[0])
			continue
		}

		if currentHash != "" && (isDigitOrDash(trimmed[0])) {
			fields := strings.Split(trimmed, "\t")
			if len(fields) >= 2 {
				add := parseIntOrZero(fields[0])
				del := parseIntOrZero(fields[1])
				entry := stats[currentHash]
				entry.additions += add
				entry.deletions += del
				stats[currentHash] = entry
			}
		}
	}
	return stats, nil
}

func isDigitOrDash(b byte) bool {
	return (b >= '0' && b <= '9') || b == '-'
}

func parseIntOrZero(s string) int {
	s = strings.TrimSpace(s)
	if s == "-" {
		return 0
	}
	v, _ := strconv.Atoi(s)
	return v
}
