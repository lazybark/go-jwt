package storage

import (
	"fmt"
	"time"
)

var (
	ErrEntityExists   = fmt.Errorf("entity_exists")
	ErrEntityNotExist = fmt.Errorf("entity_not_exist")
	ErrInternal       = fmt.Errorf("internal_storage_error")
)

type User struct {
	ID int `json:"user_id"`
	//ServiceId represents service that uses ths user data
	ServiceId    Service   `json:"service_id"`
	Login        string    `json:"login"`
	PasswordHash string    `json:"pwd_hash"`
	Name         string    `json:"name"`
	SecName      string    `json:"sec_name"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	//LastActivity sets time of last auth request on user's behalf
	LastActivity time.Time `json:"last_activity"`

	//PermissionUsers sets level of permissions to control users
	//in the auth system
	PermissionUsers Permission `json:"users_permission"`
}

type Service int

const (
	//Unversal is a user that can perform ops on any service
	Unversal Service = iota + 1
	RedBackend
)

func (s Service) String() string {
	return [...]string{"unversal", "red_backend_auth"}[s]
}

// Permission divides user access level by data category.
// Higher levels include all lower ones
type Permission int

const (
	UsersDenied Permission = iota
	UsersRead
	UsersWrite
	UsersCreate
	UsersDelete
)

var (
	UserSystem = User{
		ServiceId:       Unversal,
		Login:           "SYSTEM_USER",
		PasswordHash:    "retina-misc1-monstrous-23",
		Name:            "SYSTEM_USER",
		SecName:         "",
		Email:           "",
		PermissionUsers: UsersDelete,
	}
)
