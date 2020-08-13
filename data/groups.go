package data

import "fmt"

// ErrGroupNotFound is an error raised when a group can not be found in the database
var ErrGroupNotFound = fmt.Errorf("Group not found")

// Group defines the structure for an API group
type Group struct {
	// the id of the group
	//
	// required: false
	// min: 1
	Id int `json:"id"`

	// the name for the group
	//
	// required: true
	// max length: 255
	Name string `json:"name"`
}
