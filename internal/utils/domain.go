package utils

import (
	"regexp"
)

var domainRegexp = regexp.MustCompile(`^(?i)[a-z0-9-]+(\.[a-z0-9-]+)+\.?$`)

// IsValidDomain returns true if the domain is valid.
//
// It uses a simple regular expression to check the domain validity.
func IsValidDomain(domain string) bool {
	return domainRegexp.MatchString(domain)
}
