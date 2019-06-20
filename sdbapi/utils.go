package sdbapi

import "strings"

func IsDuplicateEntryError(err error) bool {
	if strings.Contains(err.Error(), "Duplicate entry") {
		return true
	}
	return false
}
