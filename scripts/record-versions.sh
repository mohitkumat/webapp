#!/usr/bin/env bash
set -euo pipefail
v() { echo "$1: $($1 2>/dev/null | head -1 || echo 'not installed')"; }
echo "os: $(uname -s) $(uname -r) $(uname -m)"
v "git --version"
v "gh --version"
v "node --version"
v "pnpm --version"
v "go version"
v "python3 --version"
v "uv --version"
v "docker --version"
v "docker compose version"
v "claude --version"
v "codex --version"
v "curl --version"
v "jq --version"
v "openssl version"
