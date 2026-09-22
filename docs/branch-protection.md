# Branch Protection Setup

CI (`.github/workflows/ci.yml`) runs on every push and PR, but it only blocks
merges once branch protection is turned on — that's a repo setting, not
something a workflow file can do by itself.

## Via GitHub UI
1. Repo → Settings → Branches → Add branch protection rule
2. Branch name pattern: `main`
3. Enable "Require status checks to pass before merging"
4. Search for and select: **CI OK** (the `ci-ok` job — it only goes green once
   both `web` and `api` pass, so you only need to require this one check)
5. Enable "Require branches to be up to date before merging"
6. Save

## Via GitHub CLI
```bash
gh api repos/:owner/:repo/branches/main/protection \
  --method PUT \
  -f required_status_checks.strict=true \
  -f 'required_status_checks.contexts[]=CI OK' \
  -f enforce_admins=true \
  -f required_pull_request_reviews.required_approving_review_count=1
```

## Why point at `ci-ok` instead of listing `web` and `api`
If you add a third CI job later (e.g. a Phase E addition, or a future
observability check), you'd have to remember to also update branch
protection settings. Pointing at the single `ci-ok` aggregator job means new
jobs are covered automatically — just add them to its `needs:` list in
`ci.yml`.
