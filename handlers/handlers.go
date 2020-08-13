package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// getId returnes the Id from the URL
// panics if it cannot convert the id into an integer
func getId(r *http.Request) int {
	// parse the user id from the url
	vars := mux.Vars(r)

	// convert the id into an integer
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}

	return id
}
