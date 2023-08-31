package handler

import "regexp"

var validSegmentNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func isValidSegmentName(name string) bool {
	return validSegmentNamePattern.MatchString(name) && name != ""
}
