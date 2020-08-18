package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/zzibert/3fs-rest-api/handlers"
	. "gopkg.in/check.v1"
)

// Creates test suite
type GroupTestSuite struct {
	groupHandler *handlers.Groups
}

// Registering test suite
func init() {
	connection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", "localhost", "5432", "zanzibert", "nekineki", "postgres")

	// Opening a connection to the postgres database
	db, err := gorm.Open("postgres", connection)
	if err != nil {
		panic(err)
	}
	l := log.New(os.Stdout, "3fs-rest-api", log.LstdFlags)

	Suite(&GroupTestSuite{
		groupHandler: handlers.NewGroups(l, db),
	})
}

// integrates with testing package
func Test(t *testing.T) { TestingT(t) }

func (s *GroupTestSuite) TestHandleGet(c *C) {
	sm := mux.NewRouter()
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/users/{id:[0-9]+}", s.groupHandler.ListSingle)
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/groups/1", nil)
	sm.ServeHTTP(writer, request)

	c.Check(writer.Code, Equals, 200)
}
