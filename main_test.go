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

// Creates test suite
type GroupTestSuite struct {
	groupHandler *handlers.Groups
	group        *data.Group
	writer       *httptest.ResponseRecorder
	mux          *mux.Router
	l            *log.Logger
	db           *gorm.DB
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

	Suite(&GroupTestSuite{
		l:  l,
		db: db,
	})
}

// integrates with testing package
func Test(t *testing.T) { TestingT(t) }

func (s *GroupTestSuite) SetUpSuite(c *C) {
	setDB(s.db)
}

func (s *GroupTestSuite) SetUpTest(c *C) {
	s.writer = httptest.NewRecorder()
	s.group = &data.Group{}
	s.mux = mux.NewRouter()
	s.groupHandler = handlers.NewGroups(s.l, s.db)
}

func (s *GroupTestSuite) TearDownSuite(c *C) {
	clearDB(s.db)
}

func setDB(db *gorm.DB) {
	db.Exec("INSERT INTO groups(name) VALUES ('group 1')")
}

func clearDB(db *gorm.DB) {
	db.Exec("delete from groups")
	db.Exec("ALTER SEQUENCE groups_id_seq RESTART WITH 1")
	db.Exec("delete from users")
	db.Exec("ALTER SEQUENCE users_id_seq RESTART WITH 1")
}

// Tries to fetch a non-existent group with id 2
func (s *GroupTestSuite) TestGroupHandleGetSingleFail(c *C) {

	getRouter := s.mux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/groups/{id:[0-9]+}", s.groupHandler.ListSingle)

	request, _ := http.NewRequest("GET", "/groups/2", nil)
	s.mux.ServeHTTP(s.writer, request)

	c.Check(s.writer.Code, Equals, 404)
}

// Tries to create a new group
func (s *GroupTestSuite) TestGroupHandlePost(c *C) {
	postRouter := s.mux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/groups", s.groupHandler.Create)

	body := strings.NewReader(`{"name": "group 2"}`)
	request, _ := http.NewRequest("POST", "/groups", body)
	s.mux.ServeHTTP(s.writer, request)

	c.Check(s.writer.Code, Equals, 200)
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

// Tries to create a group with an existing name
func (s *GroupTestSuite) TestGroupHandlePostFail(c *C) {
	postRouter := s.mux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/groups", s.groupHandler.Create)

	body := strings.NewReader(`{"name": "group 1"}`)
	request, _ := http.NewRequest("POST", "/groups", body)
	s.mux.ServeHTTP(s.writer, request)

	c.Check(s.writer.Code, Equals, 500)
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
}
