package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/zzibert/3fs-rest-api/handlers"
	. "gopkg.in/check.v1"
)

var db *gorm.DB
var l *log.Logger

func TestMain(m *testing.M) {
	initialize()
	code := m.Run()
	clearDB()
	os.Exit(code)
}

func initialize() {
	connection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", "localhost", "5433", "zanzibert", "nekineki", "test")
	l = log.New(os.Stdout, "3fs-rest-api", log.LstdFlags)
	db, _ = gorm.Open("postgres", connection)
}

func clearDB() {
	db.Exec("delete from users")
	db.Exec("alter sequence users_id_seq restart with 1")
	db.Exec("delete from groups")
	db.Exec("alter sequence groups_id_seq restart with 1")
}

type GroupTestSuite struct {
	mux          *http.ServeMux
	writer       *httptest.ResponseRecorder
	groupHandler *handlers.Groups
}

type UserTestSuite struct {
	mux         *http.ServeMux
	writer      *httptest.ResponseRecorder
	userHandler *handlers.Users
}

func init() {
	Suite(&GroupTestSuite{})
	Suite(&UserTestSuite{})
}

func Test(t *testing.T) { TestingT(t) }

func (s *PostTestSuite) SetUpSuite(c *C) {
	s.mux = http.NewServeMux()
	s.groupHandler = handlers.NewGroups(l, db)
	s.mux.HandleFunc("/groups", s.groupHandler.ListAll)
	s.writer = httptest.NewRecorder()
}

func (s *PostTestSuite) TearDownSuite(c *C) {

}

func (s *PostTestSuite) TestHandleGet(c *C) {
	request, _ := http.NewRequest("GET", "/groups", nil)
	s.mux.ServeHTTP(s.writer, request)

	c.Check(s.writer.Code, Equals, 200)
}
