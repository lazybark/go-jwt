package redis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/lazybark/go-helpers/cli/clf"
	"github.com/lazybark/go-jwt/storage"
	"golang.org/x/crypto/bcrypt"
)

func (r *Redis) GenerateUserId() int {
	r.mutexUID.Lock()
	*r.lastUID++
	r.mutexUID.Unlock()
	return *r.lastUID
}

func (r Redis) UserAdd(u storage.User) (int, error) {
	//Check if such login exists
	exists, err := r.CheckExistense(fmt.Sprintf(keys["logins"], u.Login, int(u.ServiceId)))
	if err != nil {
		return 0, storage.ErrInternal
	}
	if exists {
		return 0, storage.ErrEntityExists
	}
	//Generate new ID
	id := r.GenerateUserId()

	//Add user data to user data list
	u.ID = id
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	u.LastActivity = time.Now()

	pwdHashed, err := r.hashAndSaltPassword([]byte(u.PasswordHash))
	if err != nil {
		return 0, storage.ErrInternal
	}

	u.PasswordHash = pwdHashed

	um, err := json.Marshal(u)
	if err != nil {
		return 0, storage.ErrInternal
	}
	err = r.db.Set(fmt.Sprintf(keys["users"], id), um, 0).Err()
	if err != nil {
		return 0, storage.ErrInternal
	}
	//Add user login to login list
	err = r.db.Set(fmt.Sprintf(keys["logins"], u.Login, int(u.ServiceId)), id, 0).Err()
	if err != nil {
		fmt.Println(clf.Red("ERROR SETTING USER LOGIN RECORD. MAIN RECORD WILL BE DELETED for ", id))
		//Remove user data from data list
		err = r.db.Del(fmt.Sprintf(keys["users"], id)).Err()
		if err != nil {
			fmt.Println(clf.Red("ERROR DELETING BAD USER RECORD for", id))
		}
		return 0, storage.ErrInternal
	}

	return id, nil
}

func (r Redis) UserGetData(login string, service int) (storage.User, error) {
	//Check if user exists (has login & id associated with it)
	u := storage.User{}
	uid := 0
	lBytes, err := r.GetKey(fmt.Sprintf(keys["logins"], login, service))
	if err != nil {
		return u, err
	}
	//Result should be user id > 0
	if err := json.Unmarshal(lBytes, &uid); err != nil || uid == 0 {
		fmt.Println(clf.Red(err))
		return u, storage.ErrInternal
	}

	uBytes, err := r.GetKey(fmt.Sprintf(keys["users"], uid))
	if err != nil {
		return u, err
	}
	if err := json.Unmarshal(uBytes, &u); err != nil {
		return u, err
	}

	return u, nil

}

func (r Redis) UserUpdateActivity(uid int) error {

	return nil
}

func (r *Redis) hashAndSaltPassword(pwd []byte) (string, error) {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), err
}
