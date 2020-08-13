package handlers

import "log"

// Groups Handler for getting and updating groups
type Groups struct {
	l *log.Logger
}

// NewGroups returns a new groups handler
func NewGroups(l *log.Logger) *Groups {
	return &Groups{l}
}
