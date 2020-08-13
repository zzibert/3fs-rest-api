package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/zzibert/3fs-rest-api/data"
)

// Users handler for getting and updating users
type Users struct {
	l  *log.Logger
	Db *sql.DB
}

// NewUsers returns a new users handler with the given logger
func NewUsers(l *log.Logger) *Users {
	return &Users{l}
}

// ErrInvalidUserPath is an error when user path is not valid
var ErrInvalidUserPath = fmt.Errorf("Invalid path, path should be /users/[id]")

// GenericError is a generic error message
type GenericError struct {
	Message string `json:"message"`
}

// swagger:route GET /users users ListUsers
// Returns a list of users from the database
// responses:
//  200: UsersResponse

// ListAll handles GET requests and returns all current users
func (u *Users) ListAll(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Get all users")

	users := data.GetUsers()

	err := data.ToJSON(users, rw)
	if err != nil {
		u.l.Println("error encoding users")
	}
}

// swagger:route GET /users/{id} users ListUser
// Returns a single user from the database
// responses:
//  200: userResponse
//  404: errorResponse

// ListSingle handles GET requests with id
func (u *Users) ListSingle(rw http.ResponseWriter, r *http.Request) {
	id := getId(r)

	u.l.Println("get user id", id)

	user, err := data.GetUserById(id)

	switch err {

	}
}

// swagger:route PUT /users users updateUser
// update an user
//
// responses:
//  201: noContentResponse
//  404: errorResponse

// Update handles PUT to update users
func (u *Users) Update(rw http.ResponseWriter, r *http.Request) {

	// create the user from the request body
	var user data.User
	err := data.FromJSON(user, r.Body)
	if err != nil {
		u.l.Println("Error couldnt parse user from request body", err)

		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	err = data.UpdateUser(user)
	if err != nil {
		u.l.Println("Error user not found", user)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: "User not found in database"}, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

// swagger:route POST /users users createUser
// Create a new User
//
// responses:
//  200: userResponse
//  501: errorResponse

// Create handles POST requests to add new users
func (u *Users) Create(rw http.ResponseWriter, r *http.Request) {
	var user data.User
	err := data.FromJSON(user, r.Body)
	if err != nil {
		u.l.Println("Error couldnt parse user from request body", err)

		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	data.AddUser(user)
}

// swagger:route DELETE /users/{id} users deleteUser
// Deletes an user from the database
//
// responses:
//  201: noContentResponse
//  404: errorResponse
//  501: errorResponse

// DeleteUser handles DELETE requests for deleting an user from the database
func (u *Users) Delete(rw http.ResponseWriter, r *http.Request) {
	id := getId(r)

	u.l.Println("Deleting user with id", id)

	err := data.DeleteUser(id)
	switch err {
	case nil:

	case data.ErrUserNotFound:
		u.l.Println("Error user id does not exist")

		rw.WriteHeader(http.StatusNotFound)
		return
	default:
		u.l.Println("Error deleting user", err)

		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
