package regex

import (
	"regexp"
	"strconv"
)

// TextToFloat32 is a helper function to find and convert the first number in the given text to a float32
func TextToFloat32(text string) *float32 {
	if text == "" {
		return nil
	}

	re := regexp.MustCompile(NumberRegex)
	numberStr := re.FindString(text)

	num, err := strconv.ParseFloat(numberStr, 64)
	if err != nil {
		return nil
	}

	num32 := float32(num)

	return &num32
}

func TextToUint(text string) *uint {
	num := TextToFloat32(text)
	if num == nil {
		return nil
	}

	numUint := uint(*num)
	return &numUint
}
