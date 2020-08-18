package data

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// ErrUserNotFound is an error raised when a user can not be found in the database
var ErrUserNotFound = fmt.Errorf("User not found")

// ErrUserConstraintViolation is an error raised when an user can not be created because of constraint violations
var ErrUserConstraintViolation = fmt.Errorf("user has constraints violation")

// User defines the structure for an API User
// swagger:model
type User struct {
	// the id of the user
	//
	// required: false
	// min:1
	ID int `json:"id"`

	// the name of the user
	//
	// required: true
	// max length: 255
	Name string `json:"name"`

	// the email of the user
	//
	// required: true
	// max length: 255
	Email string `json:"email"`

	// the password of the user
	//
	// required: true
	// max length: 255
	Password string `json:"password"`

	// the id of the group that the user belongs to
	//
	// required: true
	// min: 1
	GroupID int `json:"groupID"`

	// The group that the user belongs to
	//
	// required: false
	Group Group `json:"-"`
}

// GetUsers returns all users from the database
func GetUsers(db *gorm.DB) (users []*User) {
	db.Find(&users)
	return
}

// GetUserById returns a single user with the specified id
// If the user is not found this func retuns UserNotFound error
func GetUserById(id int, db *gorm.DB) (user User, err error) {
	user.ID = id
	if err = db.First(&user).Error; err != nil {
		err = ErrUserNotFound
	}
	return
}

// UpdateUser replaces the set of values within the given user
// If a user is not found this func returns a UserNotFound error
// if the update would make a constraint violation the func returns a ErrUserConstraintViolation error
func UpdateUser(id int, userMap map[string]interface{}, db *gorm.DB) (err error) {
	var user User
	if err = db.First(&user, id).Error; err != nil {
		err = ErrUserNotFound
		return
	}

	if err = db.Model(&user).Updates(userMap).Error; err != nil {
		err = ErrUserConstraintViolation
	}
	return
}

// AddUser adds a user to the database
// if the user would make a constraint violation the func returns a ErrUserConstraintViolation error
func AddUser(user *User, db *gorm.DB) (err error) {
	if err = db.Create(user).Error; err != nil {
		err = ErrUserConstraintViolation
	}
	return
}

// DeleteUser deletes an user from the database
func DeleteUser(id int, db *gorm.DB) (err error) {
	var user User
	if err = db.First(&user, id).Error; err != nil {
		err = ErrUserNotFound
	} else {
		db.Delete(&user)
	}
	return
}
