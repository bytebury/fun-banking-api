package controllers

import "strings"

/**
 * Determines if the given error is associated to a SQL UNIQUE constraint.
 */
func isDuplicateError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "duplicate key")
}

func isRecordNotFoundError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "record not found")
}
