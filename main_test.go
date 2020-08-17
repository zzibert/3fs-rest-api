package main

import (
	. "gopkg.in/check.v1"
)

// Creates test suite
type GroupTestSuite struct {
}

// Registering test suite
func init() {
	Suite(&GroupTestSuite{})
}
