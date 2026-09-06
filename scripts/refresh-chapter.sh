#!/usr/bin/env bash
# Vendor a chapter's frozen snapshots from the two source repos.
#
#   scripts/refresh-chapter.sh ch-08 --genesis v0.4 --maklaude v0.5
#   scripts/refresh-chapter.sh ch-05 --maklaude v0.3
#
# Omit a repo's flag when the chapter left that repo untouched. Its directory is
# then left exactly as it is, which for a chapter that never had one means the
# reader falls back to the most recent chapter that did change it.
#
# Source trees default to ~/git/genesis and ~/git/MaKlaude. Override with
# GENESIS_SRC and MAKLAUDE_SRC. Uses `git archive`, so the export carries tracked
# files only: no `.git`, no ignored credentials, no local scratch.
set -euo pipefail

usage() {
    echo "usage: $0 <ch-NN> [--genesis TAG] [--maklaude TAG] [--with-history]" >&2
    exit 2
}

GENESIS_SRC="${GENESIS_SRC:-$HOME/git/genesis}"
MAKLAUDE_SRC="${MAKLAUDE_SRC:-$HOME/git/MaKlaude}"

CH="${1:-}"
[ -n "$CH" ] || usage
shift

GENESIS_TAG=""
MAKLAUDE_TAG=""
while [ $# -gt 0 ]; do
    case "$1" in
        --genesis)  GENESIS_TAG="${2:-}"; shift 2 ;;
        --maklaude) MAKLAUDE_TAG="${2:-}"; shift 2 ;;
        --with-history) WITH_HISTORY=1; shift ;;
        *) echo "unknown argument: $1" >&2; usage ;;
    esac
done

[ -n "$GENESIS_TAG" ] || [ -n "$MAKLAUDE_TAG" ] || {
    echo "nothing to do: pass --genesis and/or --maklaude" >&2
    exit 2
}

# Run from the companion repo root, so a mistyped chapter cannot scatter trees
# into whatever directory the caller happened to be sitting in.
cd "$(dirname "$0")/.."
[ -d "$CH" ] || { echo "no such chapter directory: $CH" >&2; exit 1; }

vendor() {
    local name="$1" src="$2" tag="$3" dest="$CH/$1"

    [ -d "$src" ] || { echo "$name: no source tree at $src" >&2; exit 1; }
    git -C "$src" rev-parse -q --verify "refs/tags/$tag" >/dev/null || {
        echo "$name: tag $tag does not exist in $src" >&2
        exit 1
    }

    # A tag that exists only locally makes the snapshot unreproducible for
    # anyone reading the book, which is the one thing this repo is for.
    local remote
    remote="$(git -C "$src" remote get-url origin)"
    git -C "$src" ls-remote --tags "$remote" "refs/tags/$tag" 2>/dev/null | grep -q . || {
        echo "$name: tag $tag is not pushed to $remote; push it first" >&2
        exit 1
    }

    rm -rf "$dest"
    mkdir -p "$dest"
    git -C "$src" archive "$tag" | tar -x -C "$dest"
    echo "$name $tag -> $dest ($(find "$dest" -type f | wc -l | tr -d ' ') files)"

    # Optional: ship the revision history as a bundle beside the plain files.
    # A chapter that MEASURES how the code changed cannot be served by a flat
    # export, because `git log` against it returns nothing. Chapter 10 is that
    # chapter.
    #
    # Two details are deliberate. The bundle is built from the TAG alone, never
    # from `main`, so it stops exactly where the chapter stops and a reader
    # cannot accidentally read future work into a chapter's conclusions. And it
    # stays a plain file rather than a nested `.git`: git converts an embedded
    # `.git` into a gitlink at mode 160000, so clones of THIS repo would carry
    # a submodule pointer and no history at all. Measured, not assumed.
    if [ -n "${WITH_HISTORY:-}" ]; then
        # Absolute, because `git -C "$src"` resolves a relative output path
        # against the SOURCE repo rather than this one.
        local bundle="$PWD/$CH/$name.bundle"
        rm -f "$bundle"
        git -C "$src" bundle create "$bundle" "$tag" >/dev/null 2>&1 || {
            echo "$name: could not build a history bundle for $tag" >&2
            exit 1
        }
        # Prove it before claiming it. A bundle that cannot be cloned is worse
        # than no bundle, because the README promises a reader it works.
        local probe
        probe="$(mktemp -d)"
        git clone -q --branch "$tag" "$bundle" "$probe/x" 2>/dev/null || {
            echo "$name: the bundle for $tag does not clone; refusing to ship it" >&2
            rm -rf "$probe"; exit 1
        }
        echo "$name $tag -> $bundle ($(du -h "$bundle" | cut -f1), $(git -C "$probe/x" rev-list --count HEAD) commits)"
        rm -rf "$probe"
    fi
}

[ -n "$GENESIS_TAG" ]  && vendor genesis  "$GENESIS_SRC"  "$GENESIS_TAG"
[ -n "$MAKLAUDE_TAG" ] && vendor MaKlaude "$MAKLAUDE_SRC" "$MAKLAUDE_TAG"

# Belt and suspenders on the two things that would be embarrassing in a public
# companion repo: a nested .git, and a credential the source repo forgot to
# gitignore. Neither should be possible via git archive, which is why they are
# cheap to assert rather than trust.
if find "$CH" -name .git -print -quit | grep -q .; then
    echo "refusing to leave a nested .git under $CH" >&2
    exit 1
fi
if find "$CH" -name '*.pem' -o -name '.env' -print -quit | grep -q .; then
    echo "warning: a credential-shaped file landed under $CH; check before committing" >&2
fi

echo "done. review with: git status --short $CH"
