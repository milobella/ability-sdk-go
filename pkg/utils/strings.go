package utils

import "strings"

func StringSliceContains(slice []string, element string) bool {
	for _, sliceElt := range slice {
		if sliceElt == element {
			return true
		}
	}
	return false
}

func StringFuzzyMatch(original string, match string) bool {
	originalLower := strings.ToLower(original)
	matchLower := strings.ToLower(match)
	return strings.Contains(originalLower, matchLower) || strings.Contains(matchLower, originalLower)
}
