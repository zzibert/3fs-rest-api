// Package classification 3fs API
//
// Documentation for 3fs API
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package handlers

import "github.com/zzibert/3fs-rest-api/data"

// Generic error message returned as a string
// swagger:response errorResponse
type errorResponseWrapper struct {
	// Description of the error
	// in: body
	Body GenericError
}

// GenericError is a generic error message
type GenericError struct {
	Message string `json:"message"`
}

// A list of groups
// swagger:response groupsResponse
type groupsResponseWrapper struct {
	// All current groups
	// in: body
	Body []data.Group
}

// A list of users
// swagger:response usersResponse
type usersResponseWrapper struct {
	// All current users
	// in: body
	Body []data.User
}

// A single group
// swagger:response groupResponse
type groupResponseWrapper struct {
	// a single group
	// in: body
	Body data.Group
}

// A single user
// swagger:response userResponse
type userResponseWrapper struct {
	// a single user
	// in: body
	Body data.User
}

// No content is returned by this API endpoint
// swagger:response noContentResponse
type noContentResponseWrapper struct {
}
