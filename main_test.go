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
	_ "github.com/lib/pq"
	"github.com/zzibert/3fs-rest-api/data"
	"github.com/zzibert/3fs-rest-api/handlers"
	. "gopkg.in/check.v1"
)

var db *gorm.DB
var l *log.Logger

// // Main Test func, initiliazing, running and clearing each test run
func TestMain(m *testing.M) {
	initialize()
	code := m.Run()
	clearDB()
	os.Exit(code)
}

// initializing the database connection and creating the logger
func initialize() {
	connection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", "localhost", "5433", "zanzibert", "nekineki", "test")
	l = log.New(os.Stdout, "3fs-rest-api", log.LstdFlags)
	db, _ = gorm.Open("postgres", connection)
}

// After each test run, clearing the database and reseting the index sequences
func clearDB() {
	db.Exec("delete from users")
	db.Exec("alter sequence users_id_seq restart with 1")
	db.Exec("delete from groups")
	db.Exec("alter sequence groups_id_seq restart with 1")
}

// The group test suite
type GroupTestSuite struct {
	mux          *mux.Router
	writer       *httptest.ResponseRecorder
	groupHandler *handlers.Groups
	group        *data.Group
}

// the user test suite
type UserTestSuite struct {
	mux         *mux.Router
	writer      *httptest.ResponseRecorder
	userHandler *handlers.Users
}

// Registering the test suites
func init() {
	Suite(&GroupTestSuite{})
	Suite(&UserTestSuite{})
}

// Integrating with the test package
func Test(t *testing.T) { TestingT(t) }

func (s *GroupTestSuite) SetUpSuite(c *C) {
	s.group = &data.Group{}
	s.mux = mux.NewRouter()
	s.groupHandler = handlers.NewGroups(l, db)
	s.writer = httptest.NewRecorder()

	getRouter := s.mux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/groups", s.groupHandler.ListAll)
	getRouter.HandleFunc("/groups/{id:[0-9]+}", s.groupHandler.ListSingle)
	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/groups/{id:[0-9]+}", groupHandler.Update)

}

func (s *GroupTestSuite) TearDownSuite(c *C) {
	clearDB()
}

// Testing that get all groups returns 0 groups
func (s *GroupTestSuite) TestHandleGetAllEmpty(c *C) {
	request, _ := http.NewRequest("GET", "/groups", nil)
	s.mux.ServeHTTP(s.writer, request)

	c.Check(s.writer.Code, Equals, 200)
	var groups []data.Group
	data.FromJSON(groups, s.writer.Body)
	c.Check(len(groups), Equals, 0)
}

// Creating a new group
// func (s *GroupTestSuite) TestHandlePost(c *C) {
// 	body := strings.NewReader(`{"name": "group 1"}`)
// 	request, _ := http.NewRequest("POST", "/groups", body)
// 	s.mux.ServeHTTP(s.writer, request)

// 	c.Check(s.writer.Code, Equals, 200)
// }

// Getting a single group with id
// func (s *GroupTestSuite) TestHandleListSingle(c *C) {
// 	request, _ := http.NewRequest("GET", "/groups/1", nil)
// 	s.mux.ServeHTTP(s.writer, request)

// 	var group data.Group
// 	data.FromJSON(group, s.writer.Body)

// 	c.Check(s.writer.Code, Equals, 200)
// 	c.Check(group.Name, Equals, "group 1")
// }

// Trying to create a group with the same name
// func (s *GroupTestSuite) TestHandlePostFail(c *C) {
// 	requestBody := strings.NewReader(`{"name": "group 1"}`)
// 	request, _ := http.NewRequest("POST", "/groups", requestBody)
// 	s.mux.ServeHTTP(s.writer, request)

// 	c.Check(s.writer.Code, Equals, 200)
// }
