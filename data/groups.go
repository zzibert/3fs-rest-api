package data

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// ErrGroupNotFound is an error raised when a group can not be found in the database
var ErrGroupNotFound = fmt.Errorf("Group not found")

// ErrGroupConstraintViolation is an error raised when a group can not be created because of constraint violations
var ErrGroupConstraintViolation = fmt.Errorf("group has constraints violation")

// Group defines the structure for an API group
// swagger:model
type Group struct {
	// the id of the group
	//
	// required: false
	// min: 1
	ID int `json:"id"`

	// the name for the group
	//
	// required: true
	// max length: 255
	Name string `json:"name"`

	// the list of users belonging to this group
	//
	// required: false
	Users []User `json:"users"`
}

// GetGroups returns all groups from the database
func GetGroups(db *gorm.DB) (groups []*Group) {
	db.Preload("Users").Find(&groups)
	return
}

// GetGroupById returns a single group with the specified id
// If a group is not found this func returns a GroupNotFound error
func GetGroupById(id int, db *gorm.DB) (group Group, err error) {
	group.ID = id
	if err = db.Preload("Users").First(&group).Error; err != nil {
		err = ErrGroupNotFound
	}
	return
}

// UpdateGroup replaces the set of values within the given group
// If a group is not found this func returns a GroupNotFound error
// if the update would make a constraint violation the func returns a ErrGroupConstraintViolation error
func UpdateGroup(id int, groupMap map[string]interface{}, db *gorm.DB) (err error) {
	var group Group
	if err = db.First(&group, id).Error; err != nil {
		err = ErrGroupNotFound
		return
	}

	if err = db.Model(&group).Updates(groupMap).Error; err != nil {
		err = ErrGroupConstraintViolation
	}
	return
}

// AddGroup adds a group to the database
// if the group would make a constraint violation the func returns a ErrGroupConstraintViolation error
func AddGroup(group *Group, db *gorm.DB) (err error) {
	if err = db.Create(group).Error; err != nil {
		err = ErrGroupConstraintViolation
	}
	return
}

// DeleteGroup deletes a group from the database
// if the deletion of the group would make a constraint violation the func returns a ErrGroupConstraintViolation error
func DeleteGroup(id int, db *gorm.DB) (err error) {
	var group Group
	if err = db.First(&group, id).Error; err != nil {
		err = ErrGroupNotFound
		return
	}
	if err = db.Delete(&group).Error; err != nil {
		err = ErrGroupConstraintViolation
	}
	return
}
