# Code Companion - *Autonomous AI Infrastructure with Claude Code*

Each `ch-NN/` directory holds a frozen snapshot of the two repositories the book builds:

- **genesis** - the bootstrapper that scaffolds autonomous dev systems.
- **MaKlaude** - the dev system genesis bootstraps in the running example (autonomous Kubernetes operations).

The snapshots are vendored as **plain files (no `.git`)**, pinned to the exact version that chapter ends on, so you can read the precise code each chapter describes without checking out tags yourself.

## Version per chapter

Each repo carries its own tag, and a chapter pins whichever tag each repo ends on. The two started in lockstep but now diverge: a chapter that only touches MaKlaude advances MaKlaude's tag while genesis stays put. A repo is re-vendored only in the chapters where it actually changed, so where a chapter left one untouched its directory is omitted and you read the snapshot from the most recent chapter that did change it. A chapter can leave both untouched, which is what ch-07 does.

| Chapter | genesis | MaKlaude | End state |
|---------|---------|----------|-----------|
| ch-03   | `v0.1`  | `v0.1`   | Bootstrapped: scaffold published, onboarding done |
| ch-04   | `v0.2`  | `v0.2`   | Milestone 1 complete: read-only foundation, shipped and verified |
| ch-05   | unchanged (`v0.2`) | `v0.3` | Human-in-the-loop steering: comms channels and escalation |
| ch-06   | `v0.3`  | `v0.4`   | The evolver: self-improvement, seed backport, the two-tier gap |
| ch-07   | unchanged (`v0.3`) | unchanged (`v0.4`) | Permissions and trust: the chapter audits the posture already shipped rather than adding to it |
| ch-08   | `v0.4`  | `v0.5`   | Observability: hooks to Loki, the activity dashboard, session cost and the continuation ladder |

## Refreshing a chapter's snapshot

Use the script. Pass the chapter and a tag for each repo the chapter changed, and omit the flag for a repo it left alone:

```bash
scripts/refresh-chapter.sh ch-08 --genesis v0.4 --maklaude v0.5
scripts/refresh-chapter.sh ch-05 --maklaude v0.3          # genesis unchanged
```

It exports with `git archive`, so the snapshot carries tracked files only: no `.git`, no ignored credentials, no local scratch. It reads local clones at `~/git/genesis` and `~/git/MaKlaude`, overridable with `GENESIS_SRC` and `MAKLAUDE_SRC`, and it refuses to vendor a tag that hasn't been pushed, since a local-only tag would leave the snapshot unreproducible for anyone reading the book.

Both source repos are public: `Sayfan-AI/genesis` and `Sayfan-AI/MaKlaude`. Without local clones, clone the tag first and point the script at it:

```bash
git clone --branch v0.5 https://github.com/Sayfan-AI/MaKlaude.git /tmp/MaKlaude
MAKLAUDE_SRC=/tmp/MaKlaude scripts/refresh-chapter.sh ch-08 --maklaude v0.5
```

## Notes

- Snapshots are vendored as plain files, not git submodules, so the companion stays readable, offline, and pinned.
- The nested `.github/workflows/` inside each snapshot never run. GitHub Actions only executes workflows at a repository's own root, not in vendored copies.
- No secrets are included. The export carries only tracked files, and both repos gitignore their credentials (App private keys, tokens, central `.env`).
