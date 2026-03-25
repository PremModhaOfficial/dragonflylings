# Open Questions

## dragonflylings - 2025-03-25
- [ ] Rename Go module from `dragonflyLearnings` to `dragonflylings`? -- Affects all import paths, should decide before scaffolding begins
- [ ] Test isolation strategy: unique key prefixes per test vs FLUSHDB in setup? -- Unique prefixes are safer for parallel test runs but more verbose
- [ ] Should the existing `cmd/main.go` be archived to `cmd/sandbox/main.go` as a free-form playground? -- Preserves user's existing work while making room for the runner CLI
- [ ] Dragonfly version to pin in docker-compose.yml? -- Newer versions may have different feature sets (e.g., JSON support)
- [ ] Should capstone exercises have time limits to simulate interview/production pressure? -- Adds realism but may frustrate learners
- [ ] Include a `solutions/` directory with reference implementations, or keep solutions only in test assertions? -- Solutions help self-checking but risk learners peeking too early
