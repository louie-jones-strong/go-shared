package condition

import (
	"testing"

	"github.com/louie-jones-strong/go-shared/returnchecker"
	"github.com/stretchr/testify/assert"
)

func TestUnit_Evaluate(t *testing.T) {

	rc := returnchecker.New1ReturnChecker[Condition[int]](t)

	alwaysTrueCond := NewMockCondition[int](true, nil)
	alwaysFalseCond := NewMockCondition[int](false, nil)
	errorCond := NewMockCondition[int](false, assert.AnError)

	tests := []struct {
		name        string
		cond        Condition[int]
		expectedRes bool
		expectedErr error
	}{
		{
			name:        "always true condition",
			cond:        alwaysTrueCond,
			expectedRes: true,
			expectedErr: nil,
		},
		{
			name:        "always false condition",
			cond:        alwaysFalseCond,
			expectedRes: false,
			expectedErr: nil,
		},
		{
			name:        "always error condition",
			cond:        errorCond,
			expectedRes: false,
			expectedErr: assert.AnError,
		},
		{
			name: "and condition with single true sub-conditions",
			cond: rc.Check(NewAndCondition(
				alwaysTrueCond,
			)),
			expectedRes: true,
			expectedErr: nil,
		},
		{
			name: "and condition with multiple true sub-conditions",
			cond: rc.Check(NewAndCondition(
				alwaysTrueCond,
				alwaysTrueCond,
				alwaysTrueCond,
			)),
			expectedRes: true,
			expectedErr: nil,
		},
		{
			name: "and condition with a false sub-condition",
			cond: rc.Check(NewAndCondition(
				alwaysTrueCond,
				alwaysFalseCond,
				alwaysTrueCond,
			)),
			expectedRes: false,
			expectedErr: nil,
		},
		{
			name: "and condition where a sub-condition returns error",
			cond: rc.Check(NewAndCondition(
				alwaysTrueCond,
				errorCond,
			)),
			expectedRes: false,
			expectedErr: assert.AnError,
		},
		{
			name: "or condition with single true sub-conditions",
			cond: rc.Check(NewOrCondition(
				alwaysTrueCond,
			)),
			expectedRes: true,
			expectedErr: nil,
		},
		{
			name: "or condition with multiple false sub-conditions",
			cond: rc.Check(NewOrCondition(
				alwaysFalseCond,
				alwaysFalseCond,
			)),
			expectedRes: false,
			expectedErr: nil,
		},
		{
			name: "or condition with a later true sub-condition",
			cond: rc.Check(NewOrCondition(
				alwaysFalseCond,
				alwaysTrueCond,
			)),
			expectedRes: true,
			expectedErr: nil,
		},
		{
			name: "or condition where a sub-condition returns error",
			cond: rc.Check(NewOrCondition(
				alwaysFalseCond,
				errorCond,
			)),
			expectedRes: false,
			expectedErr: assert.AnError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			res, err := tc.cond.Evaluate(0)

			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
				assert.False(t, res)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRes, res)
			}
		})
	}
}
