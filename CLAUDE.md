# CLAUDE.md

## 프로젝트
git-contrib — Go CLI, 멀티리포 Git 기여도 분석 → HTML/JSON 리포트 생성

## 구조
```
cmd/main.go          # cobra CLI, author alias 전처리
pkg/config/          # YAML 설정 로딩/검증, AuthorAlias()
pkg/git/             # git log 실행/파싱, errgroup 병렬 수집
pkg/analyzer/        # 집계 (작성자별, 기간별, 리포별), json 태그 포함
pkg/report/          # HTML(embed.FS 템플릿) / JSON 출력
```

## 빌드 & 테스트
```bash
make build           # VERSION.txt에서 버전 읽어 ldflags 주입
make test
./bin/git-contrib    # config.yaml 기반 실행
./bin/git-contrib --version
```

## 릴리즈
VERSION.txt 범프 → CHANGELOG.md/README.md 갱신 → 커밋 → `git tag vX.Y.Z && git push origin vX.Y.Z`
GitHub Actions가 GoReleaser 빌드 + homebrew-tap 자동 갱신

## 주의사항
- `config.yaml`은 .gitignore 대상 (민감 경로 포함), `config.yaml.sample`만 추적
- HTML 템플릿 내 JS 데이터는 `data-*` 속성으로 전달 (`<script>` 내 Go 변수 이스케이핑 문제 회피)
- analyzer types에 json 태그 필수 유지
- 리포 이름은 `git_contrib` (underscore), go 모듈은 `git-contrib` (hyphen)
