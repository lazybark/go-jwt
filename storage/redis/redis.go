package redis

import (
	"crypto/tls"
	"errors"
	"strconv"
	"sync"

	"github.com/go-redis/redis"
	"github.com/lazybark/go-jwt/storage"
)

type Redis struct {
	db       *redis.Client
	lastUID  *int
	mutexUID *sync.Mutex
}

var keys = map[string]string{
	"logins":   "login-%s-%v", //user-login-serviceID
	"users":    "user-%v",     //user-id
	"last_uid": "last_uid",    //last user id issued
}

func NewRedisStorage(RedisHost string, RedisPassword string, RedisTLS bool, DB int) (storage.Storage, error) {
	options := &redis.Options{Addr: RedisHost, Password: RedisPassword, DB: DB}
	if RedisTLS {
		options.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	db := redis.NewClient(options)
	if err := db.Ping().Err(); err != nil {
		return nil, err
	}

	r := Redis{db: db, lastUID: new(int), mutexUID: &sync.Mutex{}}

	//If we have any prev key stored in db - restore it
	bts, err := r.GetKey(keys["last_uid"])
	if err != nil && err != storage.ErrEntityNotExist {
		return nil, err
	}
	if len(bts) > 0 {
		bk, err := strconv.Atoi(string(bts))
		if err != nil {
			return nil, err
		}
		r.lastUID = &bk
	}

	return r, nil
}

func (r Redis) Flush() error {
	*r.lastUID = 0
	return r.db.FlushDB().Err()
}

func (r Redis) Init() error {
	_, err := r.UserAdd(storage.UserSystem)
	if err != nil {
		return err
	}
	return err
}

func (r *Redis) CheckKeyExistense(key string) (bool, error) {
	//If error is redis.Nil - no such entity
	//Any other error means actual error
	//No error means entity exists
	err := r.db.Get(key).Err()
	if err == nil {
		return true, nil
	}
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	return false, err
}

func (r *Redis) GetKey(key string) ([]byte, error) {
	bytes, err := r.db.Get(key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, storage.ErrEntityNotExist
		}
		return nil, err
	}
	if len(bytes) == 0 {
		return nil, storage.ErrEntityNotExist
	}
	return bytes, nil
}

func (r *Redis) GetSet(key string) (map[string]string, error) {
	bytes, err := r.db.HGetAll(key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, storage.ErrEntityNotExist
		}
		return nil, err
	}
	return bytes, nil
}
