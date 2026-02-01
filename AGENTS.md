# Coding Guide

This document contains concise coding standards and guidelines for contributors and for the agent editing files in this repository.

## Goals
- Keep code readable, consistent, and idiomatic for Go and project conventions.
- Make small, focused changes; avoid wide refactors.
- Run and verify tests/build where feasible after modifications.

## General Principles
- Prefer clarity over cleverness.
- Keep functions small and single-responsibility.
- Avoid changing unrelated files or formatting unless required.
- Add a small unit test for behavioral changes when practical.

## Go-specific
- Use `gofmt`/`go fmt` formatting. Follow `go vet` and `revive` linting rules in `revive.toml`.
- Keep exported identifiers capitalized; internal unexported lowercase.
- Return errors rather than panicking in library code. Panic only for programmer errors and truly unrecoverable states.
- All errors should be handled or explicitly ignored with comments.
- Prefer concrete types in function signatures where appropriate; use interfaces for behaviour only.
- When editing files, minimize API surface changes (don't rename public types/functions without PR-level coordination).

## New Features / Changes
- Add a focused TODO to track the change if it spans multiple steps.
- Run `go test ./...` and `go vet` locally or via CI before finalising changes.

## Agent Behaviour Guidance
- Always add a TODO entry for multi-step tasks and update it as steps complete.
- Before making changes, read the current file contents to avoid clobbering recent edits.
- If adding imports, ensure they are used; avoid shadowing package names with local variables.
- If code compiles errors after edits, attempt up to 3 small fixes; if unresolved, ask the user for guidance.

## Quick Checklist Before Pushing
- [ ] Code builds locally (`go build ./...`).
- [ ] Tests pass (`go test ./...`).
- [ ] Linter passes or lint issues are justified.
- [ ] TODO list updated to reflect progress.

## Unit Test Style & Tips

- **Use table-driven tests:** Prefer table-driven style for variations and edge-cases. Keep each test case small and self-contained.
- **Test helpers:** Put common test helpers in `_test.go` helper files in the same package to keep tests readable and DRY.
- **Test data files:** Store fixtures under `test-data/` and load them with explicit paths. Use `t.TempDir()` for temporary files.
- **Naming:** Name tests `TestXxx` for unit tests and `BenchmarkXxx` for benchmarks. Use descriptive subtest names in `t.Run`.
- **Assertions:** Use `require`/`assert` from `testify` when helpful (project vendor includes `stretchr/testify`). Prefer `require` when failing fast is appropriate.
- **Isolation & determinism:** Avoid shared global state; run tests deterministically (seed rand when needed). Use `t.Parallel()` only when the test is safe to run concurrently.
- **Mocks & fakes:** Provide small, focused mocks in `mock.go` or local test files rather than large frameworks. Make behaviour explicit in the mock implementation.
- **Error handling:** Check and assert errors (do not ignore them). Prefer returning errors from helpers so tests can surface failures clearly.
- **Small, focused tests:** Each test should assert a single behaviour. If setup is heavy, consider helper builders or fixtures.
- **Benchmarks:** Keep benchmarks isolated and use `b.ResetTimer()`/`b.ReportAllocs()` appropriately. Run with `go test -bench=. -benchmem`.
- **Local runs:** Use `go test ./...` for full run; use `-run`, `-bench`, and `-v` for focused debugging: e.g. `go test ./... -run TestName -v`.
- **Formatting & vet:** Run `gofmt`/`go fmt`, `go vet`, and project linters before pushing changes.

These practices reflect patterns and existing tests in this repository: table-driven tests, `test-data/` usage, small test helpers, and vendorized `testify` usage. Follow them to keep tests consistent and maintainable.

## Example: Table-driven Unit Test

Below is an example table-driven unit test pattern used in the project. It demonstrates subtests, inputs, expected results and error assertions.

```go
func TestUnit_FunctionName(t *testing.T) {

	tests := []struct {
		name        string
		input1      int
		input2      int
		expectedRes int
		expectedErr error
	}{
		{
			name:        "Test Case 1",
			input1:      1,
			input2:      2,
			expectedRes: 0,
			expectedErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			res, err := FunctionName(tc.input1, tc.input2)

			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRes, res)
			}
		})
	}
}
```

Notes:
- Replace `FunctionName` with the function under test.
- This example uses `assert` helpers (from `testify`) — if your package avoids external deps, replace with standard `testing` checks.