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

// Groups defines a slice of Group
type Groups []*Group

// GetGroups returns all groups from the database
func GetGroups() Groups {

}

// GetGroupById returns a single group with the specified id
// If a group is not found this func returns a GroupNotFound error
func GetGroupById(id int) (*Group, error) {

}

// UpdateGroup replaces a group with the given item
// If a group is not found this func returns a GroupNotFound error
func UpdateGroup(g Group) error {

}

// AddGroup adds a group to the database
func AddGroup(g Group) {

}

// DeleteGroup deletes a group from the database
func DeleteGroup(id int) error {

}
