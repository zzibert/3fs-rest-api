package data

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// ErrUserNotFound is an error raised when a user can not be found in the database
var ErrUserNotFound = fmt.Errorf("User not found")

// User defines the structure for an API User
// swagger:model
type User struct {
	// the id of the user
	//
	// required: false
	// min:1
	ID int `json:"id" gorm:"primary_key"`

	// the name of the user
	//
	// required: true
	// max length: 255
	Name string `json:"name" gorm:"type:varchar(20);unique"`

	// the email of the user
	//
	// required: true
	// max length: 255
	Email string `json:"email" gorm:"type:varchar(100);unique"`

	// the password of the user
	//
	// required: true
	// max length: 255
	Password string `json:"password" gorm:"type:varchar(20)"`

	// the id of the group that the user belongs to
	//
	// required: true
	// min: 1
	GroupID int `json:"groupId" `

	// The group that the user belongs to
	//
	// required: false
	Group Group `json:"-" gorm:"foreignkey:GroupID"`
}

// GetUsers returns all users from the database
func GetUsers(db *gorm.DB) (users []*User) {
	db.Find(&users)
	return
}

// GetUserById returns a single user with the specified id
// If the user is not found this func retuns UserNotFound error
func GetUserById(id int, db *gorm.DB) (user User, err error) {
	if err = db.First(&user, id).Error; err != nil {
		err = ErrUserNotFound
	}
	return
}

// UpdateUser replaces a user with the given item
// If the user is not found this func returns UserNotFound error
func UpdateUser(originalUser User, user User, db *gorm.DB) (err error) {
	return db.Model(&originalUser).Updates(user).Error
}

// AddUser adds a new user to the database
func AddUser(user *User, db *gorm.DB) (err error) {
	err = db.Create(user).Error
	return
}

// DeleteUser deletes an user from the database
func DeleteUser(id int, db *gorm.DB) (err error) {
	var user *User
	if err = db.First(user, id).Error; err != nil {
		err = ErrUserNotFound
	} else {
		db.Delete(user)
	}
	return
}
