# Chapter 11 - Building Your Own - From Goal to Running System

Frozen snapshots of the three repositories as they stood at the end of this chapter. They
are plain files with no `.git`, so you can read exactly the code the chapter describes
without checking out tags yourself.

| Directory | Repository | Pinned at |
|-----------|------------|-----------|
| `genesis/` | [`Sayfan-AI/genesis`](https://github.com/Sayfan-AI/genesis) - the bootstrapper that scaffolds autonomous dev systems | `v0.6` |
| `MaKlaude/` | [`Sayfan-AI/MaKlaude`](https://github.com/Sayfan-AI/MaKlaude) - the dev system genesis bootstraps, for autonomous Kubernetes operations | `v0.6` |
| `genesis-dashboard/` | [`Sayfan-AI/genesis-dashboard`](https://github.com/Sayfan-AI/genesis-dashboard) - the dev system this chapter builds from a goal, a dashboard of a dev system's state | `v0.1` |

**End state:** Genesis Dashboard bootstrapped from a goal and run in local mode through
milestone 1, and genesis running an evolver of its own. MaKlaude is unchanged from
Chapter 10. `genesis-dashboard-goal.md` is the goal the chapter feeds to genesis.

## What you need

Only if you want to *run* the code. Reading it needs nothing installed.

| Tool | Version | Needed for |
|------|---------|-----------|
| [Claude Code](https://docs.claude.com/en/docs/claude-code) | 2.1 | Running the agents. All three repos are driven through it |
| Python | 3.12 or later, with [`uv`](https://docs.astral.sh/uv/) | `genesis/` - the scaffolder and the local control plane |
| [`gh` CLI](https://cli.github.com) | any current | `genesis/` - repo creation, issues, workflow control. Needs `gh auth login` |
| A GitHub App | - | `genesis/` - the identity a dev system acts as. Permissions are listed in `genesis/README.md` |
| Node.js | 22 | `genesis-dashboard/` - building and testing the dashboard |
| Go | 1.24 | `MaKlaude/` - building and testing the dev system |
| Docker | any current | `MaKlaude/` - the container runtime under KinD |
| [KinD](https://kind.sigs.k8s.io) | `kindest/node:v1.31.2` | `MaKlaude/` - the local Kubernetes cluster its end-to-end tests run against |

Those are the versions this chapter was written and tested against. Later versions will
usually work, and the pins are what you fall back to when something behaves oddly.

## Running it

Each snapshot carries its own README with the commands for that tree:

```bash
cd genesis           && cat README.md
cd genesis-dashboard && cat README.md
cd MaKlaude          && cat README.md
```

A quick check that your toolchain is in order. Run it on a copy, so the snapshot stays
clean. A few genesis tests ask git what the repo tracks, so its copy needs to be a git
repository first:

```bash
cp -R genesis /tmp/genesis && cd /tmp/genesis
git init -q && git add -A
uv sync --locked && uv run --no-sync pytest -q

cp -R genesis-dashboard /tmp/genesis-dashboard && cd /tmp/genesis-dashboard
npm ci && npm test && npm run typecheck
```

## Notes

- The nested `.github/workflows/` directories never run. GitHub Actions only executes
  workflows at a repository's own root, not inside a vendored copy.
- No credentials are vendored. The snapshots carry tracked files only. genesis and
  MaKlaude gitignore their App private keys, tokens and `.env` files, and
  genesis-dashboard never writes one: its GitHub token lives in browser memory only.
- To refresh a snapshot, see `scripts/refresh-chapter.sh` and the root `README.md`.
