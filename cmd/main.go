package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/devlikebear/git-contrib/pkg/analyzer"
	"github.com/devlikebear/git-contrib/pkg/config"
	"github.com/devlikebear/git-contrib/pkg/git"
	"github.com/devlikebear/git-contrib/pkg/report"
	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	var (
		cfgPath string
		since   string
		until   string
		output  string
		repos   []string
	)

	rootCmd := &cobra.Command{
		Use:     "git-contrib",
		Short:   "Git 멀티리포 개발자 기여도 분석 도구",
		Long:    "여러 Git 리포지토리를 분석하여 개발자별 기여 통계를 HTML 리포트로 생성합니다.",
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("config load: %w", err)
			}

			// CLI flag overrides
			if since != "" {
				cfg.Since = since
			}
			if until != "" {
				cfg.Until = until
			}
			if output != "" {
				cfg.Output = output
			}
			if len(repos) > 0 {
				cfg.Repos = repos
			}

			if err := cfg.Validate(); err != nil {
				return fmt.Errorf("config validation: %w", err)
			}

			fmt.Fprintf(os.Stderr, "\n⚡ Git Contribution Analyzer\n")
			fmt.Fprintf(os.Stderr, "   기간: %s → %s\n", cfg.Since, cfg.Until)
			fmt.Fprintf(os.Stderr, "   레포: %d개\n\n", len(cfg.Repos))

			commits, err := git.CollectAll(cfg.Repos, cfg.Since, cfg.Until)
			if err != nil {
				return fmt.Errorf("collect: %w", err)
			}

			// Apply author aliases
			if aliasMap := cfg.AuthorAlias(); len(aliasMap) > 0 {
				for i := range commits {
					if canonical, ok := aliasMap[strings.ToLower(commits[i].Author)]; ok {
						commits[i].Author = canonical
					}
				}
			}

			if len(commits) == 0 {
				return fmt.Errorf("해당 기간에 커밋 데이터가 없습니다")
			}

			fmt.Fprintf(os.Stderr, "✅ 총 %d개 커밋 분석 완료\n", len(commits))

			rpt := analyzer.Analyze(commits, cfg.Since, cfg.Until)

			if err := report.Generate(cfg.Output, rpt); err != nil {
				return fmt.Errorf("report: %w", err)
			}

			fmt.Fprintf(os.Stderr, "📄 리포트 생성: %s\n", cfg.Output)
			fmt.Fprintf(os.Stderr, "\n   open %s  # 브라우저에서 열기\n\n", cfg.Output)
			return nil
		},
	}

	rootCmd.Flags().StringVarP(&cfgPath, "config", "c", "config.yaml", "설정 파일 경로")
	rootCmd.Flags().StringVar(&since, "since", "", "시작 날짜 (YYYY-MM-DD)")
	rootCmd.Flags().StringVar(&until, "until", "", "종료 날짜 (YYYY-MM-DD)")
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "출력 파일명")
	rootCmd.Flags().StringSliceVarP(&repos, "repos", "r", nil, "분석할 리포지토리 경로 (쉼표 구분)")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
