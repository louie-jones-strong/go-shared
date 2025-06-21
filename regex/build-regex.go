package regex

import (
	"regexp"
	"strings"
)

func BuildRegex(pattenParts ...string) *regexp.Regexp {
	patten := strings.Join(pattenParts, "")
	re := regexp.MustCompile(patten)

	return re
}
