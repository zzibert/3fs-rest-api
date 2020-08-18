package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/zzibert/3fs-rest-api/data"
	"github.com/zzibert/3fs-rest-api/handlers"
	. "gopkg.in/check.v1"
)

// Creates group test suite
type GroupTestSuite struct {
	groupHandler *handlers.Groups
	group        *data.Group
	writer       *httptest.ResponseRecorder
	mux          *mux.Router
	l            *log.Logger
	db           *gorm.DB
}

// Creates user test suite
type UserTestSuite struct {
	userHandler *handlers.Users
	user        *data.User
	writer      *httptest.ResponseRecorder
	mux         *mux.Router
	l           *log.Logger
	db          *gorm.DB
}

// Registering test suite
func init() {
	connection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", "localhost", "5433", "zanzibert", "nekineki", "test")

	// Opening a connection to the postgres database
	db, err := gorm.Open("postgres", connection)
	if err != nil {
		panic(err)
	}
	l := log.New(os.Stdout, "3fs-rest-api", log.LstdFlags)

	Suite(&GroupTestSuite{l: l, db: db})
	Suite(&UserTestSuite{l: l, db: db})
}

// integrates with testing package
func Test(t *testing.T) { TestingT(t) }

func (s *GroupTestSuite) SetUpTest(c *C) {
	s.writer = httptest.NewRecorder()
	s.group = &data.Group{}
	s.mux = mux.NewRouter()
	s.groupHandler = handlers.NewGroups(s.l, s.db)
	setDB(s.db)
}

func (s *GroupTestSuite) TearDownTest(c *C) {
	clearDB(s.db)
}

func (s *UserTestSuite) SetUpTest(c *C) {
	s.writer = httptest.NewRecorder()
	s.user = &data.User{}
	s.mux = mux.NewRouter()
	s.userHandler = handlers.NewUsers(s.l, s.db)
	setDB(s.db)
}

func (s *UserTestSuite) TearDownTest(c *C) {
	clearDB(s.db)
}

func setDB(db *gorm.DB) {
	db.Exec("INSERT INTO groups(name) VALUES ('group 1')")
	db.Exec("INSERT INTO groups(name) VALUES ('group 2')")
	db.Exec("INSERT INTO users(name, password, email, group_id) VALUES ('user 1', 'pass', 'user@email.com', 1)")
	db.Exec("INSERT INTO users(name, password, email, group_id) VALUES ('user 2', 'pass', 'user2@email.com', 1)")
}

func clearDB(db *gorm.DB) {
	db.Exec("delete from users")
	db.Exec("ALTER SEQUENCE users_id_seq RESTART WITH 1")
	db.Exec("delete from groups")
	db.Exec("ALTER SEQUENCE groups_id_seq RESTART WITH 1")
}

//GROUP TESTS

// Tries to fetch a non-existent group with id 2
func (s *GroupTestSuite) TestGroupHandleGetSingleFail(c *C) {

	getRouter := s.mux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/groups/{id:[0-9]+}", s.groupHandler.ListSingle)

	request, _ := http.NewRequest("GET", "/groups/3", nil)
	s.mux.ServeHTTP(s.writer, request)

	c.Check(s.writer.Code, Equals, 404)
}

// Tries to fetch a group with id 1
func (s *GroupTestSuite) TestGroupHandleGetSingle(c *C) {

	getRouter := s.mux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/groups/{id:[0-9]+}", s.groupHandler.ListSingle)

	request, _ := http.NewRequest("GET", "/groups/1", nil)
	s.mux.ServeHTTP(s.writer, request)

	c.Check(s.writer.Code, Equals, 200)
	json.Unmarshal(s.writer.Body.Bytes(), s.group)
	c.Check(s.group.Name, Equals, "group 1")
}

// Tries to create a new group
func (s *GroupTestSuite) TestGroupHandlePost(c *C) {
	postRouter := s.mux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/groups", s.groupHandler.Create)

	body := strings.NewReader(`{"name": "group 3"}`)
	request, _ := http.NewRequest("POST", "/groups", body)
	s.mux.ServeHTTP(s.writer, request)

	c.Check(s.writer.Code, Equals, 200)
}

// Tries to create a group with an existing name
func (s *GroupTestSuite) TestGroupHandlePostFail(c *C) {
	postRouter := s.mux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/groups", s.groupHandler.Create)

	body := strings.NewReader(`{"name": "group 1"}`)
	request, _ := http.NewRequest("POST", "/groups", body)
	s.mux.ServeHTTP(s.writer, request)

	c.Check(s.writer.Code, Equals, 400)
}

// Tries to fetch all groups
func (s *GroupTestSuite) TestGroupHandleGetAll(c *C) {

	getRouter := s.mux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/groups", s.groupHandler.ListAll)

	request, _ := http.NewRequest("GET", "/groups", nil)
	s.mux.ServeHTTP(s.writer, request)

	c.Check(s.writer.Code, Equals, 200)

	var groups *[]data.Group
	json.Unmarshal(s.writer.Body.Bytes(), &groups)
	c.Check((*groups)[0].Name, Equals, "group 1")
	c.Check((*groups)[1].Name, Equals, "group 2")
}

// Trying to update a group with id 1
func (s *GroupTestSuite) TestGroupHandlePut(c *C) {
	putRouter := s.mux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/groups/{id:[0-9]+}", s.groupHandler.Update)

	body := strings.NewReader(`{"name": "new group name"}`)
	request, _ := http.NewRequest("PUT", "/groups/1", body)
	s.mux.ServeHTTP(s.writer, request)

	c.Check(s.writer.Code, Equals, 204)

	getRouter := s.mux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/groups/{id:[0-9]+}", s.groupHandler.ListSingle)

	request, _ = http.NewRequest("GET", "/groups/1", nil)
	s.writer = httptest.NewRecorder()
	s.mux.ServeHTTP(s.writer, request)

	c.Check(s.writer.Code, Equals, 200)
	json.Unmarshal(s.writer.Body.Bytes(), s.group)
	c.Check(s.group.Name, Equals, "new group name")
}

//trying to delete a group with users referenced to it
func (s *GroupTestSuite) TestGroupHandleDelete(c *C) {
	deleteRouter := s.mux.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/groups/{id:[0-9]+}", s.groupHandler.Delete)

	request, _ := http.NewRequest("DELETE", "/groups/1", nil)
	s.mux.ServeHTTP(s.writer, request)

	c.Check(s.writer.Code, Equals, 400)
}

//trying to delete a group
func (s *GroupTestSuite) TestGroupHandleDeleteFailOne(c *C) {
	deleteRouter := s.mux.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/groups/{id:[0-9]+}", s.groupHandler.Delete)

	request, _ := http.NewRequest("DELETE", "/groups/2", nil)
	s.mux.ServeHTTP(s.writer, request)

	c.Check(s.writer.Code, Equals, 204)
}

// Trying to delete an non-existend group
func (s *GroupTestSuite) TestGroupHandleDeleteFailTwo(c *C) {
	deleteRouter := s.mux.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/groups/{id:[0-9]+}", s.groupHandler.Delete)

	request, _ := http.NewRequest("DELETE", "/groups/5", nil)
	s.mux.ServeHTTP(s.writer, request)

	c.Check(s.writer.Code, Equals, 404)
}

// USERS TESTS

// Tries to fetch a non-existent user with id 3
func (s *UserTestSuite) TestUserHandleGetSingleFail(c *C) {

	getRouter := s.mux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/users/{id:[0-9]+}", s.userHandler.ListSingle)

	request, _ := http.NewRequest("GET", "/users/3", nil)
	s.mux.ServeHTTP(s.writer, request)

	c.Check(s.writer.Code, Equals, 404)
}

// Tries to fetch a user with id 1
func (s *UserTestSuite) TestUserHandleGetSingle(c *C) {

	getRouter := s.mux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/users/{id:[0-9]+}", s.userHandler.ListSingle)

	request, _ := http.NewRequest("GET", "/users/1", nil)
	s.mux.ServeHTTP(s.writer, request)

	c.Check(s.writer.Code, Equals, 200)
	json.Unmarshal(s.writer.Body.Bytes(), s.user)
	c.Check(s.user.Name, Equals, "user 1")
}

// Tries to create a new user
func (s *UserTestSuite) TestUserHandlePost(c *C) {
	postRouter := s.mux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/users", s.userHandler.Create)

	body := strings.NewReader(`{"name": "user 3", "password": "pass", "email": "user3@email.com", "groupID": 1}`)
	request, _ := http.NewRequest("POST", "/users", body)
	s.mux.ServeHTTP(s.writer, request)

	c.Check(s.writer.Code, Equals, 200)
}

// Tries to create a user with an existing name
func (s *UserTestSuite) TestUserHandlePostFail(c *C) {
	postRouter := s.mux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/users", s.userHandler.Create)

	body := strings.NewReader(`{"name": "user 1", "password": "pass", "email": "user43@email.com", "groupID": 1}`)
	request, _ := http.NewRequest("POST", "/users", body)
	s.mux.ServeHTTP(s.writer, request)

	c.Check(s.writer.Code, Equals, 400)
}

// Tries to fetch all users
func (s *UserTestSuite) TestUserHandleGetAll(c *C) {

	getRouter := s.mux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/users", s.userHandler.ListAll)

	request, _ := http.NewRequest("GET", "/users", nil)
	s.mux.ServeHTTP(s.writer, request)

	c.Check(s.writer.Code, Equals, 200)

	var users []data.Group
	json.Unmarshal(s.writer.Body.Bytes(), &users)
	c.Check(users[0].Name, Equals, "user 1")
	c.Check(users[1].Name, Equals, "user 2")
}

// Trying to update a user with id 1
func (s *UserTestSuite) TestUserHandlePut(c *C) {
	putRouter := s.mux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/users/{id:[0-9]+}", s.userHandler.Update)

	body := strings.NewReader(`{"name": "new user name"}`)
	request, _ := http.NewRequest("PUT", "/users/1", body)
	s.mux.ServeHTTP(s.writer, request)

	c.Check(s.writer.Code, Equals, 204)

	getRouter := s.mux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/users/{id:[0-9]+}", s.userHandler.ListSingle)

	request, _ = http.NewRequest("GET", "/users/1", nil)
	s.writer = httptest.NewRecorder()
	s.mux.ServeHTTP(s.writer, request)

	c.Check(s.writer.Code, Equals, 200)
	json.Unmarshal(s.writer.Body.Bytes(), s.user)
	c.Check(s.user.Name, Equals, "new user name")
}

// Trying to update an users groupID to non-existent group
func (s *UserTestSuite) TestUserHandlePutFail(c *C) {
	putRouter := s.mux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/users/{id:[0-9]+}", s.userHandler.Update)

	body := strings.NewReader(`{"groupID": 66}`)
	request, _ := http.NewRequest("PUT", "/users/1", body)
	s.mux.ServeHTTP(s.writer, request)

	c.Check(s.writer.Code, Equals, 404)
}

// Trying to delete a user
func (s *UserTestSuite) TestUserHandleDelete(c *C) {
	deleteRouter := s.mux.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/users/{id:[0-9]+}", s.userHandler.Delete)

	request, _ := http.NewRequest("DELETE", "/users/1", nil)
	s.mux.ServeHTTP(s.writer, request)

	c.Check(s.writer.Code, Equals, 204)
}

// Trying to delete a non-existent user
func (s *UserTestSuite) TestUserHandleDeleteFail(c *C) {
	deleteRouter := s.mux.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/users/{id:[0-9]+}", s.userHandler.Delete)

	request, _ := http.NewRequest("DELETE", "/users/14", nil)
	s.mux.ServeHTTP(s.writer, request)

	c.Check(s.writer.Code, Equals, 404)
}
