package regex

const (
	NumberRegex               = `(\d+(\.\d+)?)`
	ZeroOrWhiteSpace          = `\s*`
	EndOrWhiteSpace           = `(?:\s|$)`
	EndWhiteSpaceOrBrackets   = `(?:\s|[(\[})\]},]|$)`
	StartOrWhiteSpace         = `(?:^|\s)`
	StartWhiteSpaceOrBrackets = `(?:^|\s|[(\[})\]},])`
)
