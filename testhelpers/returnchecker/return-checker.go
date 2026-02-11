package returnchecker

import "testing"

type baseReturnChecker struct {
	t *testing.T
}

func (rc *baseReturnChecker) checkErr(err error) {
	rc.t.Helper()
	if err != nil {
		rc.t.Fatalf("unexpected error: %v", err)
	}
}

type ReturnChecker1[R any] struct {
	baseReturnChecker
}

func New1ReturnChecker[R any](t *testing.T) *ReturnChecker1[R] {
	return &ReturnChecker1[R]{baseReturnChecker{t: t}}
}

func (rc *ReturnChecker1[R]) Check(res R, err error) R {
	rc.t.Helper()
	rc.checkErr(err)
	return res
}

type ReturnChecker2[R1 any, R2 any] struct {
	baseReturnChecker
}

func New2ReturnChecker[R1 any, R2 any](t *testing.T) *ReturnChecker2[R1, R2] {
	return &ReturnChecker2[R1, R2]{baseReturnChecker{t: t}}
}

func (rc *ReturnChecker2[R1, R2]) Check(res1 R1, res2 R2, err error) (R1, R2) {
	rc.t.Helper()
	rc.checkErr(err)
	return res1, res2
}

type ReturnChecker3[R1 any, R2 any, R3 any] struct {
	baseReturnChecker
}

func New3ReturnChecker[R1 any, R2 any, R3 any](t *testing.T) *ReturnChecker3[R1, R2, R3] {
	return &ReturnChecker3[R1, R2, R3]{baseReturnChecker{t: t}}
}

func (rc *ReturnChecker3[R1, R2, R3]) Check(res1 R1, res2 R2, res3 R3, err error) (R1, R2, R3) {
	rc.t.Helper()
	rc.checkErr(err)
	return res1, res2, res3
}
