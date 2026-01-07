package condition

type ConditionType string

const (
	ConditionTypeNone ConditionType = ""

	// compound condition
	ConditionTypeAnd  ConditionType = "and"
	ConditionTypeOr   ConditionType = "or"
	ConditionTypeNor  ConditionType = "nor"
	ConditionTypeNand ConditionType = "nand"
)
