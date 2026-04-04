# git-contrib

여러 Git 리포지토리를 분석하여 개발자별 기여 통계를 HTML/JSON 리포트로 생성하는 CLI 도구입니다.

## Features

- **멀티리포 분석** — 여러 리포지토리를 동시에 분석 (goroutine 병렬 수집)
- **개발자별 기여 통계** — 커밋 수, 추가/삭제 LOC, 순증감, 활동일, 일평균 커밋/LOC
- **기간별 상세** — 일간/주간/월간 탭으로 Top 20 개발자의 기간별 활동 확인
- **리포별 탭** — 종합 + 리포별 독립 분석 (통계, 차트, 순위, 기간별 상세)
- **타임시리즈 차트** — 일간 활동 라인 차트 (커밋/LOC 전환, 개발자별 필터)
- **Author alias** — 동일 인물의 여러 git username을 하나로 병합
- **다크 테마 HTML 리포트** — 단일 HTML 파일, 외부 의존성 없음
- **JSON 출력** — 분석 결과를 구조화된 JSON으로 내보내기

## Installation

### Homebrew (macOS)

```bash
brew tap devlikebear/tap
brew install git-contrib
```

### Go install

```bash
go install github.com/devlikebear/git-contrib/cmd@latest
```

### Binary

[Releases](https://github.com/devlikebear/git-contrib/releases)에서 플랫폼에 맞는 바이너리를 다운로드하세요.

## Usage

### 1. 설정 파일 생성

```bash
cp config.yaml.sample config.yaml
```

`config.yaml`을 편집하여 분석할 리포지토리 경로와 기간을 설정합니다:

```yaml
since: "2024-01-01"
until: "2024-12-31"
repos:
  - /path/to/repo1
  - /path/to/repo2
output: "report.html"

# Author alias mapping (optional, case-insensitive)
authors:
  Teknium:
    - teknium1
    - teknum
```

### 2. 실행

```bash
# config.yaml 사용
git-contrib

# CLI 플래그로 오버라이드
git-contrib --repos /path/repo1,/path/repo2 --since 2025-01-01 --until 2025-12-31

# 설정 파일 경로 지정
git-contrib -c my-config.yaml -o my-report.html

# JSON 출력 (파일 확장자로 자동 감지)
git-contrib -o report.json
```

### 3. 리포트 확인

```bash
open report.html
```

## CLI Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--config` | `-c` | 설정 파일 경로 | `config.yaml` |
| `--since` | | 시작 날짜 (YYYY-MM-DD) | config 값 |
| `--until` | | 종료 날짜 (YYYY-MM-DD) | config 값 |
| `--output` | `-o` | 출력 파일명 | config 값 |
| `--repos` | `-r` | 리포지토리 경로 (쉼표 구분) | config 값 |
| `--version` | | 버전 출력 | |

## Report Sections

### 종합 탭
- 커밋 수 바 차트 (Top 20)
- 일간 활동 타임시리즈 (커밋/LOC 전환, 개발자 필터)
- 개발자 종합 순위 테이블
- 기간별 상세 (일간/주간/월간)

### 리포별 탭
각 리포지토리별로 동일한 분석 결과를 독립적으로 확인할 수 있습니다.

## Development

```bash
# 빌드
make build

# 실행
make run

# 테스트
make test
```

## License

MIT
