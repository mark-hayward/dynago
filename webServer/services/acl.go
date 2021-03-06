package services

import (
	"errors"
)

// ACLService regulates access control.
type ACLService interface {
	CheckPermission(*User, Permission) error
}

// NewACLService creates an access control service
func NewACLService() ACLService {
	return &aclService{}
}

type aclService struct{}

// Permission is a type of permission
type Permission string

//AdministratorRole is an administrator
var AdministratorRole = "administrator"

// Role is a user role
type Role struct {
	Name        string
	Description string
	Permissions []Permission
}

// CheckPermission returns true if the user is a member of a role that has the permission.
func (a *aclService) CheckPermission(user *User, permission Permission) error {
	if user == nil {
		return errors.New("CheckPermission: No user supplied")
	}

	if permission == "" {
		return errors.New("CheckPermission: You must supply a valid permission to check against.")
	}

	if user.HasRole(AdministratorRole) {
		// Admins can do anything
		return nil
	}

	return errors.New("CheckPermission: User not authorized")
}
