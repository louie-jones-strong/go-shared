package condition

type Condition[C any] interface {
	Evaluate(conditionCtx C) (bool, error)
}

func Filter[C any](items []C, cond Condition[C]) ([]C, error) {
	res := make([]C, 0, len(items))
	for _, item := range items {
		pass, err := cond.Evaluate(item)
		if err != nil {
			return nil, err
		}
		if pass {
			res = append(res, item)
		}
	}

	return res, nil
}
