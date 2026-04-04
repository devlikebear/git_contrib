---
name: release
description: 버전 범프, 문서 갱신, 릴리즈 트리거 워크플로우
---

# Release Workflow

## 순서

1. `VERSION.txt` 버전 범프 (semver)
2. `CHANGELOG.md` 새 버전 섹션 추가 (Keep a Changelog 형식)
3. `README.md` 갱신 (새 기능 반영)
4. 커밋: `chore: bump version to X.Y.Z`
5. 푸시 → 태그 생성 → 태그 푸시

```bash
git add VERSION.txt CHANGELOG.md README.md
git commit -m "chore: bump version to X.Y.Z"
git push
git tag vX.Y.Z
git push origin vX.Y.Z
```

6. `gh run watch` 로 워크플로우 성공 확인

## 자동화 (GitHub Actions)
- 태그 `v*` 푸시 → GoReleaser 빌드 → GitHub Release 생성
- release 성공 → homebrew-tap Formula 자동 갱신
