package utils

import (
	"regexp"
	"strings"
)

var alphanumericRegex = regexp.MustCompile(`[a-z0-9]`)

// metadata.name: Invalid value: "feature/test": a lowercase RFC 1123 label must consist of lower case alphanumeric characters or '-',
// and must start and end with an alphanumeric character (e.g. 'my-name',  or '123-abc', regex used for validation is '[a-z0-9]([-a-z0-9]*[a-z0-9])?')
// https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-label-names
func EncodeRFC1123(label string) string {
	validLabel := "kka-"

	for _, r := range strings.ToLower(label) {
		if alphanumericRegex.MatchString(string(r)) {
			validLabel += string(r)
		} else {
			validLabel += "-"
		}
	}

	// Since the namespace MUST end with an alphanumeric character and still contain at most 63 characters
	// we cut the string and append "-z"
	if len(validLabel) > 60 {
		validLabel = validLabel[0:60]
		return validLabel + "-z"
	}

	return validLabel
}
