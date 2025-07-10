package apptype

import (
	"fmt"
	"strconv"
	"time"
)

type Type string

const (
	String   Type = "string"
	Int      Type = "int"
	Float    Type = "float"
	Bool     Type = "bool"
	DateTime Type = "datetime"
)

func FindType(arr []string) (Type, error) {
	var hasDates, hasFloats, hasInts, hasBools, hasStrings bool
	for _, str := range arr {
		if str == "" || str == "NaN" {
			continue
		}
		if _, err := StringToTime(str); err == nil {
			hasDates = true
		}
		if _, err := StringToInt(str); err == nil {
			hasInts = true
			continue
		}
		if _, err := StringToFloat(str); err == nil {
			hasFloats = true
			continue
		}
		if _, err := StringToBool(str); err == nil {
			hasBools = true
			continue
		}

		hasStrings = true
	}

	switch {
	case hasStrings:
		return String, nil
	case hasBools:
		return Bool, nil
	case hasFloats:
		return Float, nil
	case hasInts:
		return Int, nil
	case hasDates:
		return DateTime, nil
	default:
		return String, fmt.Errorf("couldn't detect type")
	}
}

func StringToTime(str string) (time.Time, error) {

	timeFormats := []string{
		time.Layout,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
		time.DateTime,
		time.DateOnly,
		time.TimeOnly,
	}

	for _, timeFormat := range timeFormats {
		t, err := time.Parse(timeFormat, str)
		if err == nil {
			return t, nil
		}
	}

	var out time.Time
	return out, fmt.Errorf("cannot parse string: \"%v\" as time", str)
}

func StringToInt(str string) (int, error) {
	return strconv.Atoi(str)
}

func StringToFloat(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}

func StringToBool(str string) (bool, error) {

	if str == "true" {
		return true, nil
	}

	if str == "false" {
		return false, nil
	}

	return false, fmt.Errorf("cannot convert string: %v to bool", str)
}

func ConvertArr[I any, O any](arr []I, convertor func(I) (O, error)) ([]O, error) {

	res := make([]O, len(arr))
	for i := 0; i < len(arr); i++ {
		out, err := convertor(arr[i])
		if err != nil {
			return nil, err
		}

		res[i] = out
	}

	return res, nil
}
