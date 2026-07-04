# Code Companion — *Autonomous AI Infrastructure with Claude Code*

Each `ch-NN/` directory holds a frozen snapshot of the two repositories the book builds:

- **genesis** — the bootstrapper that scaffolds autonomous dev systems.
- **MaKlaude** — the dev system genesis bootstraps in the running example (autonomous Kubernetes operations).

The snapshots are vendored as **plain files (no `.git`)**, pinned to the exact version that chapter ends on, so you can read the precise code each chapter describes without checking out tags yourself.

## Version per chapter

Both repos share a tag scheme, and each chapter advances it by one.

| Chapter | genesis | MaKlaude | End state |
|---------|---------|----------|-----------|
| ch-03   | `v0.1`  | `v0.1`   | Bootstrapped: scaffold published, onboarding done |
| ch-04   | `v0.2`  | `v0.2`   | Milestone 1 complete: read-only foundation, shipped and verified |
| ch-05+  | later tags | later tags | added as each chapter lands |

## Refreshing a chapter's snapshot

Run from the repo root. Set the chapter and tag, then drop in each repo's tree at that tag with no `.git`.

```bash
CH=ch-04
TAG=v0.2

for repo in genesis MaKlaude; do
  rm -rf "$CH/$repo"
  git clone --depth 1 --branch "$TAG" "https://github.com/Sayfan-AI/$repo.git" "$CH/$repo"
  rm -rf "$CH/$repo/.git"
done
```

`MaKlaude` is a private repository, so you need access to it. `genesis` lives at `Sayfan-AI/genesis`.

If you already have local clones of both repos, `git archive` is faster and skips the network:

```bash
CH=ch-04
TAG=v0.2

for repo in genesis MaKlaude; do
  rm -rf "$CH/$repo" && mkdir -p "$CH/$repo"
  git -C ~/git/"$repo" archive "$TAG" | tar -x -C "$CH/$repo"
done
```

## Notes

- Snapshots are vendored as plain files, not git submodules, so the companion stays readable, offline, and pinned.
- The nested `.github/workflows/` inside each snapshot never run. GitHub Actions only executes workflows at a repository's own root, not in vendored copies.
- No secrets are included. The export carries only tracked files, and both repos gitignore their credentials (App private keys, tokens, central `.env`).
