# Code Companion - *Autonomous AI Infrastructure with Claude Code*

Each `ch-NN/` directory holds a frozen snapshot of the two repositories the book builds:

- **genesis** - the bootstrapper that scaffolds autonomous dev systems.
- **MaKlaude** - the dev system genesis bootstraps in the running example (autonomous Kubernetes operations).

The snapshots are vendored as **plain files (no `.git`)**, pinned to the exact version that chapter ends on, so you can read the precise code each chapter describes without checking out tags yourself.

## Version per chapter

Each repo carries its own tag, and a chapter pins whichever tag each repo ends on. The two started in lockstep but now diverge: a chapter that only touches MaKlaude advances MaKlaude's tag while genesis stays put.

Every chapter directory is self-contained and holds both repositories, so you never have to look in a different chapter to find the code. Where the table says *unchanged*, that chapter's snapshot is byte-identical to the previous one, because the chapter didn't move that repo. ch-07 is the clearest case: it audits the security posture the earlier chapters shipped rather than adding to it, so both of its trees match ch-06.

| Chapter | genesis | MaKlaude | End state |
|---------|---------|----------|-----------|
| ch-03   | `v0.1`  | `v0.1`   | Bootstrapped: scaffold published, onboarding done |
| ch-04   | `v0.2`  | `v0.2`   | Milestone 1 complete: read-only foundation, shipped and verified |
| ch-05   | unchanged (`v0.2`) | `v0.3` | Human-in-the-loop steering: comms channels and escalation |
| ch-06   | `v0.3`  | `v0.4`   | The evolver: self-improvement, seed backport, the two-tier gap |
| ch-07   | unchanged (`v0.3`) | unchanged (`v0.4`) | Permissions and trust: the chapter audits the posture already shipped rather than adding to it |
| ch-08   | `v0.4`  | `v0.5`   | Observability: hooks to Loki, the activity dashboard, session cost and the continuation ladder |
| ch-09   | `v0.5`  | `v0.6`   | Evaluation and chaos: the replay harness, action-correctness scoring, and Kaos Monkey at both levels |
| ch-10   | unchanged (`v0.5`) | unchanged (`v0.6`) | Memory and learning: the chapter measures how the durable store behaved rather than adding to it, so neither repo advances. Ships history bundles |

## History bundles, where a chapter needs them

Snapshots are plain files, which is enough to read code but not enough to measure how it
changed. Chapter 10 measures exactly that, so `ch-10/` also carries `genesis.bundle` and
`MaKlaude.bundle`, each cloneable offline:

```bash
git clone --branch v0.6 ch-10/MaKlaude.bundle MaKlaude-history
```

Bundles are built from the chapter's **tag**, never from `main`, so they stop where the
chapter stops. They stay plain files rather than nested `.git` directories because git
turns an embedded `.git` into a gitlink, and clones of this repository would then carry a
submodule pointer with no history in it. Add one with `--with-history` (below).

The bundles are byte-reproducible, so a refresh that changes nothing leaves the file
untouched rather than showing up as a spurious diff. That takes `pack.threads=1`, which
the script sets. By default git splits delta compression across threads and the delta
base each object gets depends on how that partition fell, so two runs over the same
objects emit different bytes (287,056 against 271,547 for the same 698 objects, measured
on genesis `v0.4`). Single-threaded packing is deterministic, and it comes out slightly
smaller besides, because one delta window sees every object.

## Refreshing a chapter's snapshot

Use the script. Pass the chapter and a tag for each repo. Since every chapter carries both, an unchanged repo takes the same tag as the previous chapter:

```bash
scripts/refresh-chapter.sh ch-08 --genesis v0.4 --maklaude v0.5
scripts/refresh-chapter.sh ch-07 --genesis v0.3 --maklaude v0.4   # unchanged from ch-06
scripts/refresh-chapter.sh ch-10 --genesis v0.5 --maklaude v0.6 --with-history
```

Passing only one flag refreshes only that repo and leaves the other directory alone, which is handy when re-vendoring a single tree.

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
