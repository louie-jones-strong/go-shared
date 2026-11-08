package enum

type FlagEnum int

func (f FlagEnum) Next() FlagEnum {
	return f << 1
}

func (f FlagEnum) Has(flag FlagEnum) bool {
	return f&flag != 0
}

func (f FlagEnum) Add(flag FlagEnum) FlagEnum {
	return f | flag
}

func (f FlagEnum) Remove(flag FlagEnum) FlagEnum {
	return f &^ flag
}

func (f FlagEnum) Toggle(flag FlagEnum) FlagEnum {
	return f ^ flag
}
