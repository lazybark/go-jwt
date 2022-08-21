package mock

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/lazybark/go-jwt/storage"
)

// Mock is a mock storage for test purposes. It works
// relatively slow on big data sets and is meant to be used
// only in case small local tests
type Mock struct {
	lastUID  *int
	mutexUID *sync.Mutex
	mutexU   *sync.Mutex
	Users    map[int][]byte
}

func NewMockStorage() (storage.Storage, error) {
	return Mock{Users: make(map[int][]byte), mutexUID: &sync.Mutex{}, mutexU: &sync.Mutex{}}, nil
}

func (r *Mock) GenerateUserId() int {
	r.mutexUID.Lock()
	*r.lastUID++
	r.mutexUID.Unlock()
	return *r.lastUID
}

func (m Mock) Init() error {
	return nil
}

func (m Mock) Flush() error {
	return nil
}

func (m Mock) UserAdd(u storage.User) (int, error) {
	u.ID = m.GenerateUserId()
	uBytes, err := json.Marshal(u)
	if err != nil {
		return 0, err
	}
	m.Users[u.ID] = uBytes

	return u.ID, nil
}
func (m Mock) UserGetData(login string, service int) (storage.User, error) {
	u := storage.User{}

	//We just unmarshal every user here until we find right one.
	//No worries about efficiency in tests on this step
	for _, us := range m.Users {
		err := json.Unmarshal(us, &u)
		if err != nil {
			return u, err
		}
		if u.ServiceId == storage.Service(service) && u.Login == login {
			return u, nil
		}
	}

	return u, storage.ErrEntityNotExist
}
func (m Mock) UserUpdateActivity(uid int) error {
	u := storage.User{}

	m.mutexU.Lock()
	defer m.mutexU.Unlock()

	err := json.Unmarshal(m.Users[uid], &u)
	if err != nil {
		return err
	}

	u.LastLogin = fmt.Sprint(time.Now().Unix())
	uBytes, err := json.Marshal(u)
	if err != nil {
		return err
	}
	m.Users[uid] = uBytes

	return nil
}
func (m Mock) UserGetParam(uid string, param string) (string, error) {
	u := storage.User{}
	id, err := strconv.Atoi(uid)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(m.Users[id], &u)
	if err != nil {
		return "", err
	}

	mp := u.TransfromToHashSet()

	return mp[param].(string), nil
}
