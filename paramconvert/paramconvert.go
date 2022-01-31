package paramconvert

import "regexp"

// ToColon is a regular expression that converts {param} to :param.
var ToColon = regexp.MustCompile(`(?s)\{(.*?)\}`)

// BraceToColon converts a URL with parameters that are surrounded in braces
// to parameters that start with a colon.
func BraceToColon(URL string) string {
	return ToColon.ReplaceAllString(URL, ":$1")
}
