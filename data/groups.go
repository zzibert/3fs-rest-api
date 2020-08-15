package data

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// ErrGroupNotFound is an error raised when a group can not be found in the database
var ErrGroupNotFound = fmt.Errorf("Group not found")

// Group defines the structure for an API group
// swagger:model
type Group struct {
	// the id of the group
	//
	// required: false
	// min: 1
	ID int `json:"id" gorm:"primary_key"`

	// the name for the group
	//
	// required: true
	// max length: 255
	Name string `json:"name" gorm:"type:varchar(20);unique"`

	// the list of users belonging to this group
	//
	// required: false
	Users []User `json:"users" gorm:"foreignkey:groupID"`
}

// GetGroups returns all groups from the database
func GetGroups(db *gorm.DB) (groups []*Group) {
	db.Find(&groups)
	return
}

// GetGroupById returns a single group with the specified id
// If a group is not found this func returns a GroupNotFound error
func GetGroupById(id int, db *gorm.DB) (group Group, err error) {
	group.ID = id
	if err = db.First(&group).Error; err != nil {
		err = ErrGroupNotFound
	}
	return
}

// UpdateGroup replaces a group with the given item
// If a group is not found this func returns a GroupNotFound error
func UpdateGroup(group *Group, db *gorm.DB) (err error) {
	if err = db.Save(group).Error; err != nil {
		err = ErrGroupNotFound
	}
	return
}

// AddGroup adds a group to the database
func AddGroup(group *Group, db *gorm.DB) {
	db.Create(group)
}

// DeleteGroup deletes a group from the database
func DeleteGroup(id int, db *gorm.DB) (err error) {
	var group *Group
	if err = db.First(group, id).Error; err != nil {
		err = ErrGroupNotFound
	} else {
		db.Delete(group)
	}
	return
}
