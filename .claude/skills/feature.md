---
name: feature
description: 새 기능 추가 시 따라야 할 개발 워크플로우
---

# Feature Development

## 구현 순서

1. 관련 코드 읽기 (변경할 파일 파악)
2. 구현 (최소 변경 원칙)
3. `make build` 로 빌드 확인
4. 실제 데이터로 동작 테스트: `./bin/git-contrib`
5. 커밋 & 푸시: `feat: <설명>`

## 주요 확장 포인트

- **설정 추가**: `pkg/config/config.go` → Config 구조체 필드 추가 → `config.yaml.sample` 갱신
- **분석 로직**: `pkg/analyzer/analyzer.go` (집계), `pkg/analyzer/types.go` (구조체, json 태그 포함)
- **리포트 출력**: `pkg/report/report.go` (HTML/JSON 분기), `pkg/report/templates/report.html`
- **git 수집**: `pkg/git/git.go` (파싱), `pkg/git/types.go` (Commit 구조체)
- **CLI 플래그**: `cmd/main.go` (cobra 플래그, alias 적용 등 전처리)

## HTML 템플릿 주의사항

- Go `html/template`의 `<script>` 태그 내 변수 이스케이핑 문제 → `data-*` 속성에 JSON 삽입 후 JS에서 파싱
- 탭 전환 시 Canvas 차트는 `display:none` 상태에서 크기 0 → `setTimeout` 후 re-render 필요
- `slugify` 함수는 underscore 사용 (JS 변수명 호환)
