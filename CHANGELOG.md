# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [0.1.0] - 2026-04-04

### Added
- 멀티리포 Git 기여도 분석 (goroutine 병렬 수집)
- 개발자별 통계: 커밋 수, 추가/삭제 LOC, 순증감, 활동일, 일평균 커밋/LOC
- 다크 테마 HTML 리포트 생성 (단일 파일, 외부 의존성 없음)
- 커밋 수 바 차트 (Top 20)
- 기간별 상세 테이블 (일간/주간/월간 탭)
- 리포별 탭 전환 (종합 + 각 리포 독립 분석)
- 일간 활동 타임시리즈 차트 (Canvas 기반)
  - 커밋/LOC 지표 전환
  - 개발자별 필터 (Top 20)
  - 마우스 호버 툴팁
- YAML 설정 파일 (`config.yaml`)
- CLI 플래그 오버라이드 (`--since`, `--until`, `--repos`, `--output`, `--config`)
- `--version` 플래그
- GoReleaser 기반 릴리즈 파이프라인
- GitHub Actions: 태그 푸시 시 자동 릴리즈 + Homebrew Tap 갱신
- Homebrew 설치 지원 (`brew install devlikebear/tap/git-contrib`)

[0.1.0]: https://github.com/devlikebear/git-contrib/releases/tag/v0.1.0
