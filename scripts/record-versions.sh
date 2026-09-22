#!/usr/bin/env bash
set -euo pipefail
echo "git --version: $(git --version)"
echo "node --version: $(node --version 2>/dev/null || echo 'not installed')"
echo "pnpm --version: $(pnpm --version 2>/dev/null || echo 'not installed')"
echo "go version: $(go version 2>/dev/null || echo 'not installed')"
echo "docker --version: $(docker --version 2>/dev/null || echo 'not installed')"
echo "docker compose version: $(docker compose version 2>/dev/null || echo 'not installed')"
