package handlers

import (
	"log"
	"net/http"

	"github.com/zzibert/3fs-rest-api/data"
)

// Groups Handler for getting and updating groups
type Groups struct {
	l *log.Logger
}

// NewGroups returns a new groups handler
func NewGroups(l *log.Logger) *Groups {
	return &Groups{l}
}

// ErrInvalidGroupPath is an error message when the group path is not valid
var ErrInvalidGroupPath = fmt.Errof("Invalid Path, path should be /groups/[id]")

// swagger:route GET /groups groups ListGroups
// Return a list of groups from the database
// responses:
//  200: groupsResponse

// ListAll handles GET requests and returns all groups
func (g *Groups) ListAll(rw http.ResponseWriter, r *http.Request) {
	g.l.Println("get all groups")

	groups := data.GetGroups()

	err := data.ToJSON(groups, rw)
	if err != nil {
		g.l.Println("error encoding groups")
	}
}

// swagger:route GET /groups/{id} groups ListGroup
// returns a single group from the database
// responses:
//  200: groupResponse
//  404: errorResponse

// ListSingle handles GET requests with id parameter
func (g *Groups) ListSingle(rw http.ResponseWriter, r *http.Request) {
	id := getId(r)

	group, err := data.GetGroupById(id)

	switch err {
	case nil:

	case data.ErrGroupNotFound:
		g.l.Println("Error fetching group", err)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

	default:
		g.l.Println("Error fetching group", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = data.ToJSON(group, rw)
	if err != nil {
		g.l.Println("Error encoding group", err)
	}

}
