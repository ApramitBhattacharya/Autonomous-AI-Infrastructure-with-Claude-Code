# Chapter 8 - Observability - Watching the Watchers

Frozen snapshots of the two repositories as they stood at the end of this chapter. They
are plain files with no `.git`, so you can read exactly the code the chapter describes
without checking out tags yourself.

| Directory | Repository | Pinned at |
|-----------|------------|-----------|
| `genesis/` | [`Sayfan-AI/genesis`](https://github.com/Sayfan-AI/genesis) - the bootstrapper that scaffolds autonomous dev systems | `v0.4` |
| `MaKlaude/` | [`Sayfan-AI/MaKlaude`](https://github.com/Sayfan-AI/MaKlaude) - the dev system genesis bootstraps, for autonomous Kubernetes operations | `v0.5` |

**End state:** Observability: hooks to Loki, the activity dashboard, session cost and the continuation ladder.

## What you need

Only if you want to *run* the code. Reading it needs nothing installed.

| Tool | Version | Needed for |
|------|---------|-----------|
| [Claude Code](https://docs.claude.com/en/docs/claude-code) | 2.1 | Running the agents. Both repos are driven through it |
| Python | 3.12 or later, with [`uv`](https://docs.astral.sh/uv/) | `genesis/` - the scaffolder and the local control plane |
| [`gh` CLI](https://cli.github.com) | any current | `genesis/` - repo creation, issues, workflow control. Needs `gh auth login` |
| A GitHub App | - | `genesis/` - the identity a dev system acts as. Permissions are listed in `genesis/README.md` |
| Go | 1.24 | `MaKlaude/` - building and testing the dev system |
| Docker | any current | `MaKlaude/` - the container runtime under KinD |
| [KinD](https://kind.sigs.k8s.io) | `kindest/node:v1.31.2` | `MaKlaude/` - the local Kubernetes cluster its end-to-end tests run against |

Those are the versions this chapter was written and tested against. Later versions will
usually work, and the pins are what you fall back to when something behaves oddly.

### Also for this chapter

A free **Grafana Cloud** account, for Loki and the dashboards built on it. The chapter uses the free tier throughout, and the `GENESIS_LOKI_URL`, `GENESIS_LOKI_USER` and `GENESIS_LOKI_TOKEN` environment variables carry the credentials. Nothing here is committed.

## Running it

Each snapshot carries its own README with the commands for that tree:

```bash
cd genesis   && cat README.md
cd MaKlaude  && cat README.md
```

A quick check that your toolchain is in order:

```bash
cd genesis  && uv run --no-sync pytest -q
cd MaKlaude && go test ./...
```

## Notes

- The nested `.github/workflows/` directories never run. GitHub Actions only executes
  workflows at a repository's own root, not inside a vendored copy.
- No credentials are vendored. The snapshots carry tracked files only, and both source
  repos gitignore their App private keys, tokens and `.env` files.
- To refresh a snapshot, see `scripts/refresh-chapter.sh` and the root `README.md`.
