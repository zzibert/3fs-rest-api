package main

import (
	"testing"

	"github.com/gorilla/mux"
	. "gopkg.in/check.v1"
)

// Creates test suite
type GroupTestSuite struct {
}

// Registering test suite
func init() {
	Suite(&GroupTestSuite{})
}

// integrates with testing package
func Test(t *testing.T) { TestingT(t) }

func (s *GroupTestSuite) TestHandleGet(c *C) {
  mux := mux.NewRouter()
  mux.HandleFunc("/groups", )
}
