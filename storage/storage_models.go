package storage

import (
	"fmt"
	"strconv"
)

var (
	ErrEntityExists   = fmt.Errorf("entity_exists")
	ErrEntityNotExist = fmt.Errorf("entity_not_exist")
	ErrInternal       = fmt.Errorf("internal_storage_error")
)

type User struct {
	ID int `json:"user_id"`
	//ServiceId represents service that uses ths user data
	ServiceId    Service `json:"service_id"`
	Login        string  `json:"login"`
	PasswordHash string  `json:"pwd_hash"`
	Name         string  `json:"name"`
	LastName     string  `json:"last_name"`
	Email        string  `json:"email"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`

	//LastActivity sets time of last auth request on user's behalf
	LastLogin string `json:"last_activity"`

	//PermissionUsers sets level of permissions to control users
	//in the auth system
	PermissionUsers PermissionUsers `json:"permission_users"`

	//BytePayload stores any data that client system may want to
	//store in the ACL system. It can be info relevant to
	//client system needs: perms, roles, etc.
	BytePayload []byte `json:"byte_payload"`
}

func (u User) TransfromToHashSet() map[string]interface{} {
	fields := make(map[string]interface{})
	fields["user_id"] = u.ID
	fields["service_id"] = int(u.ServiceId)
	fields["login"] = u.Login
	fields["pwd_hash"] = u.PasswordHash
	fields["name"] = u.Name
	fields["last_name"] = u.LastName
	fields["email"] = u.Email
	fields["created_at"] = u.CreatedAt
	fields["updated_at"] = u.UpdatedAt
	fields["last_login"] = u.LastLogin
	fields["permission_users"] = int(u.PermissionUsers)
	fields["byte_payload"] = u.BytePayload

	return fields
}

// TransfromFromMap translates map fields into user parameters.
// Method does not check map fields existence, so
// map should always have ALL pre-defined fields.
func (u *User) TransfromFromMap(fields map[string]string) error {
	//Ints go first
	uid, err := strconv.Atoi(fields["user_id"])
	if err != nil {
		return err
	}
	sid, err := strconv.Atoi(fields["service_id"])
	if err != nil {
		return err
	}
	puid, err := strconv.Atoi(fields["permission_users"])
	if err != nil {
		return err
	}
	u.ID = uid
	u.ServiceId = Service(sid)
	u.PermissionUsers = PermissionUsers(puid)

	//Timings
	u.CreatedAt = fields["created_at"]
	u.UpdatedAt = fields["updated_at"]
	u.LastLogin = fields["last_login"]

	//Strings
	u.Login = fields["login"]
	u.PasswordHash = fields["pwd_hash"]
	u.Name = fields["name"]
	u.LastName = fields["last_name"]
	u.Email = fields["email"]
	u.BytePayload = []byte(fields["byte_payload"])

	return nil
}

type Service int

const (
	//Unversal is a user that can perform ops on any service
	services_start Service = iota

	ServiceUnversal
	ServiceTesting
	ServiceRedBackend

	services_end
)

func (s Service) String() string {
	if s > services_end || s < services_start {
		return "unknown"
	}
	return [...]string{"", "unversal", "lazybark_test", "red_backend_auth", ""}[s]
}

// PermissionUsers divides user access level to control other users.
// Higher levels include all lower ones
type PermissionUsers int

const (
	UsersDenied PermissionUsers = iota
	UsersRead
	UsersWrite
	UsersCreate
	UsersDelete
)

func (p PermissionUsers) Check(cp PermissionUsers) bool {
	return cp >= p
}

var (
	UserSystem = User{
		ServiceId:       ServiceUnversal,
		Login:           "SYSTEM_USER",
		PasswordHash:    "retina-misc1-monstrous-23",
		Name:            "SYSTEM_USER",
		LastName:        "",
		Email:           "",
		PermissionUsers: UsersDelete,
	}
)
