# Chapter 10 - Memory and Learning

Frozen snapshots of the two repositories as they stood at the end of this chapter, plus
their **revision history**, which this chapter is unusual in needing.

| Directory | Repository | Pinned at |
|-----------|------------|-----------|
| `genesis/` | [`Sayfan-AI/genesis`](https://github.com/Sayfan-AI/genesis) - the bootstrapper that scaffolds autonomous dev systems | `v0.5` |
| `MaKlaude/` | [`Sayfan-AI/MaKlaude`](https://github.com/Sayfan-AI/MaKlaude) - the dev system genesis bootstraps, for autonomous Kubernetes operations | `v0.6` |

**End state:** Memory and learning. The chapter measures how the durable store behaved
across the project rather than adding to it, so neither repository advances here.

## The history bundles

Every other chapter ships its snapshot as plain files, which is enough to read the code.
This chapter measures how `CLAUDE.md` *changed* over many revisions, and `git log`
against a flat export returns nothing at all. So the snapshots come with a companion
bundle each:

```bash
git clone --branch v0.5 genesis.bundle  genesis-history
git clone --branch v0.6 MaKlaude.bundle MaKlaude-history
```

No network needed. Every count in the chapter can be re-derived from those clones, and
the headline one is:

```bash
cd MaKlaude-history
git log --reverse --format='%ad' --date=short --numstat -- CLAUDE.md
```

which totals 24 revisions, 148 lines added and 8 removed.

Two things about the bundles are deliberate. They are built from the **tag**, not from
`main`, so they stop exactly where the chapter stops and no later work can leak into a
chapter's conclusions. And they are plain files rather than nested `.git` directories,
because git converts an embedded `.git` into a gitlink and clones of this repository
would carry a submodule pointer with no history in it.

## What you need

Only if you want to *run* the code. Reading it, and querying the bundles, needs only git.

| Tool | Version | Needed for |
|------|---------|-----------|
| git | any current | Cloning the history bundles above |
| [Claude Code](https://docs.claude.com/en/docs/claude-code) | 2.1 | Running the agents. Both repos are driven through it |
| Python | 3.12 or later, with [`uv`](https://docs.astral.sh/uv/) | `genesis/` - the scaffolder and the local control plane |
| [`gh` CLI](https://cli.github.com) | any current | `genesis/` - repo creation, issues, workflow control. Needs `gh auth login` |
| A GitHub App | - | `genesis/` - the identity a dev system acts as. Permissions are listed in `genesis/README.md` |
| Go | 1.24 | `MaKlaude/` - building and testing the dev system |
| Docker | any current | `MaKlaude/` - the container runtime under KinD |
| [KinD](https://kind.sigs.k8s.io) | `kindest/node:v1.31.2` | `MaKlaude/` - the local Kubernetes cluster its end-to-end tests run against |

Those are the versions this chapter was written and tested against. Later versions will
usually work, and the pins are what you fall back to when something behaves oddly.

## Notes

- The nested `.github/workflows/` directories never run. GitHub Actions only executes
  workflows at a repository's own root, not inside a vendored copy.
- No credentials are vendored. The snapshots carry tracked files only, and both source
  repos gitignore their App private keys, tokens and `.env` files.
- To refresh a snapshot, see `scripts/refresh-chapter.sh` and the root `README.md`. This
  chapter needs the `--with-history` flag.
