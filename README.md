# Code Companion — *Autonomous AI Infrastructure with Claude Code*

Each `ch-NN/` directory holds a frozen snapshot of the two repositories the book builds:

- **genesis** — the bootstrapper that scaffolds autonomous dev systems.
- **MaKlaude** — the dev system genesis bootstraps in the running example (autonomous Kubernetes operations).

The snapshots are vendored as **plain files (no `.git`)**, pinned to the exact version that chapter ends on, so you can read the precise code each chapter describes without checking out tags yourself.

## Version per chapter

Each repo carries its own tag, and a chapter pins whichever tag each repo ends on. The two started in lockstep but now diverge: a chapter that only touches MaKlaude advances MaKlaude's tag while genesis stays put. So genesis is re-vendored only in the chapters where it actually changed — where a chapter left genesis untouched, its `genesis/` directory is omitted and you read the snapshot from the most recent chapter that did change it.

| Chapter | genesis | MaKlaude | End state |
|---------|---------|----------|-----------|
| ch-03   | `v0.1`  | `v0.1`   | Bootstrapped: scaffold published, onboarding done |
| ch-04   | `v0.2`  | `v0.2`   | Milestone 1 complete: read-only foundation, shipped and verified |
| ch-05   | unchanged (`v0.2`) | `v0.3` | Human-in-the-loop steering: comms channels and escalation |
| ch-06   | `v0.3`  | `v0.4`   | The evolver: self-improvement, seed backport, the two-tier gap |

## Refreshing a chapter's snapshot

Run from the repo root. Set the chapter and each repo's tag, then drop in each tree at that tag with no `.git`. Skip the repo that didn't change in the chapter.

```bash
CH=ch-06
GENESIS_TAG=v0.3
MAKLAUDE_TAG=v0.4

rm -rf "$CH/genesis"
git clone --depth 1 --branch "$GENESIS_TAG" https://github.com/Sayfan-AI/genesis.git "$CH/genesis"
rm -rf "$CH/genesis/.git"

rm -rf "$CH/MaKlaude"
git clone --depth 1 --branch "$MAKLAUDE_TAG" https://github.com/Sayfan-AI/MaKlaude.git "$CH/MaKlaude"
rm -rf "$CH/MaKlaude/.git"
```

`MaKlaude` is a private repository, so you need access to it. `genesis` lives at `Sayfan-AI/genesis`.

If you already have local clones of both repos, `git archive` is faster and skips the network:

```bash
CH=ch-06
GENESIS_TAG=v0.3
MAKLAUDE_TAG=v0.4

rm -rf "$CH/genesis" && mkdir -p "$CH/genesis"
git -C ~/git/genesis archive "$GENESIS_TAG" | tar -x -C "$CH/genesis"

rm -rf "$CH/MaKlaude" && mkdir -p "$CH/MaKlaude"
git -C ~/git/MaKlaude archive "$MAKLAUDE_TAG" | tar -x -C "$CH/MaKlaude"
```

## Notes

- Snapshots are vendored as plain files, not git submodules, so the companion stays readable, offline, and pinned.
- The nested `.github/workflows/` inside each snapshot never run. GitHub Actions only executes workflows at a repository's own root, not in vendored copies.
- No secrets are included. The export carries only tracked files, and both repos gitignore their credentials (App private keys, tokens, central `.env`).
