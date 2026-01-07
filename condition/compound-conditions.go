package condition

import "fmt"

func checkConditionsNonNil[C any](subConditions ...Condition[C]) error {
	if len(subConditions) == 0 {
		return fmt.Errorf("compound conditions must have at least one sub-condition")
	}

	for _, subCond := range subConditions {
		if subCond == nil {
			return fmt.Errorf("sub-condition cannot be nil")
		}
	}
	return nil
}

type AndCondition[C any] struct {
	subConditions []Condition[C]
}

// NewAndCondition creates a new AndCondition instance.
func NewAndCondition[C any](subConditions ...Condition[C]) (Condition[C], error) {
	err := checkConditionsNonNil(subConditions...)
	if err != nil {
		return nil, err
	}

	cond := &AndCondition[C]{
		subConditions: subConditions,
	}
	return cond, nil
}

// Evaluate evaluates the AndCondition by evaluating all its sub-conditions.
// It returns true only if all sub-conditions evaluate to true.
// If any sub-condition evaluates to false, it returns false immediately.
// If any sub-condition returns an error, it returns that error immediately.
func (c *AndCondition[C]) Evaluate(conditionCtx C) (bool, error) {
	for _, subCondition := range c.subConditions {
		result, err := subCondition.Evaluate(conditionCtx)
		if err != nil {
			return false, err
		}
		if !result {
			return false, nil
		}
	}
	return true, nil
}

type OrCondition[C any] struct {
	subConditions []Condition[C]
}

// NewOrCondition creates a new OrCondition instance.
func NewOrCondition[C any](subConditions ...Condition[C]) (Condition[C], error) {
	err := checkConditionsNonNil(subConditions...)
	if err != nil {
		return nil, err
	}

	cond := &OrCondition[C]{
		subConditions: subConditions,
	}
	return cond, nil
}

// Evaluate evaluates the OrCondition by evaluating all its sub-conditions.
// It returns true if any sub-condition evaluates to true.
// If all sub-conditions evaluate to false, it returns false.
// If any sub-condition returns an error, it returns that error immediately.
func (c *OrCondition[C]) Evaluate(conditionCtx C) (bool, error) {
	for _, subCondition := range c.subConditions {
		result, err := subCondition.Evaluate(conditionCtx)
		if err != nil {
			return false, err
		}
		if result {
			return true, nil
		}
	}
	return false, nil
}

type NoTCondition[C any] struct {
	subCondition Condition[C]
}

func NewNorCondition[C any](subConditions ...Condition[C]) (Condition[C], error) {
	orCond, err := NewOrCondition(subConditions...)
	if err != nil {
		return nil, err
	}

	cond := &NoTCondition[C]{
		subCondition: orCond,
	}
	return cond, nil
}

func NewNandCondition[C any](subConditions ...Condition[C]) (Condition[C], error) {
	andCond, err := NewAndCondition(subConditions...)
	if err != nil {
		return nil, err
	}

	cond := &NoTCondition[C]{
		subCondition: andCond,
	}
	return cond, nil
}

func (c *NoTCondition[C]) Evaluate(conditionCtx C) (bool, error) {
	pass, err := c.subCondition.Evaluate(conditionCtx)
	if err != nil {
		return false, err
	}
	return !pass, nil
}
