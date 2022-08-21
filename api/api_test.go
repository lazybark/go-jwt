package api

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/lazybark/go-helpers/mock"
	"github.com/lazybark/go-jwt/config"
	"github.com/lazybark/go-jwt/storage/redis"
	"github.com/stretchr/testify/require"
)

func TestUsersApi(t *testing.T) {
	conf := config.Config{
		Host:   "localhost",
		Secret: "sample_secret_this_is",
	}

	rdb, err := redis.NewRedisStorage("localhost:6379", "", false, 5)
	require.NoError(t, err)

	a := New(rdb, conf)
	err = a.db.Flush()
	require.NoError(t, err)
	err = a.db.Init()
	require.NoError(t, err)

	//Create initial request data
	r := http.Request{
		Header: map[string][]string{"Cookie": {""}},
		PostForm: url.Values{
			"login":      []string{"SYSTEM_USER"},
			"pwd":        []string{"retina-misc1-monstrous-23"},
			"service_id": []string{"1"},
		},
	}
	w := mock.New()

	//Testing login under existing user
	a.ResponseUserLogin(&r, w)
	token := string(*w.Data)
	r.Header["Cookie"] = []string{"token=" + token}
	w.Flush()

	//Creating new user currectly
	r.PostForm["login"] = []string{"scnd"}
	r.PostForm["pwd"] = []string{"some"}
	r.PostForm["service_id"] = []string{"2"}
	r.PostForm["perm_users"] = []string{"3"}

	a.ResponseUserAdd(&r, w)
	w.AssertAndFlush(t, `{"success":true,"result":"2"}`)

	//Creating new user currectly: same login but different system
	r.PostForm["service_id"] = []string{"3"}
	a.ResponseUserAdd(&r, w)
	w.AssertAndFlush(t, `{"success":true,"result":"3"}`)

	//Creating new user wrong: no login, no pass, no system_id, login exists in this system
	r.PostForm["login"] = []string{""}
	a.ResponseUserAdd(&r, w)
	w.AssertAndFlush(t, `{"success":false,"code":"400","error":"bad_form_fields"}`)

	r.PostForm["login"] = []string{"scnd"}
	r.PostForm["pwd"] = []string{""}
	a.ResponseUserAdd(&r, w)
	w.AssertAndFlush(t, `{"success":false,"code":"400","error":"bad_form_fields"}`)

	r.PostForm["pwd"] = []string{"some"}
	r.PostForm["service_id"] = []string{""}
	a.ResponseUserAdd(&r, w)
	w.AssertAndFlush(t, `{"success":false,"code":"400","error":"bad_form_fields"}`)

	r.PostForm["service_id"] = []string{"3"}
	a.ResponseUserAdd(&r, w)
	w.AssertAndFlush(t, `{"success":false,"code":"400","error":"entity_exists"}`)

	//Logging under created user
	r.PostForm["service_id"] = []string{"2"}
	a.ResponseUserLogin(&r, w)
	r.Header["Cookie"] = []string{"token=" + string(*w.Data)}
	w.Flush()

	//Creating user under created user correctly
	r.PostForm["login"] = []string{"thrd"}
	r.PostForm["perm_users"] = []string{"1"}
	a.ResponseUserAdd(&r, w)
	w.AssertAndFlush(t, `{"success":true,"result":"4"}`)

	//Creating user under created user wrong: not our system_id, perms above ours
	r.PostForm["login"] = []string{"thrd4"}
	r.PostForm["service_id"] = []string{"3"}
	a.ResponseUserAdd(&r, w)
	w.AssertAndFlush(t, `{"success":false,"code":"403","error":"forbidden"}`)

	r.PostForm["service_id"] = []string{"2"}
	r.PostForm["perm_users"] = []string{"4"}
	a.ResponseUserAdd(&r, w)
	w.AssertAndFlush(t, `{"success":false,"code":"403","error":"forbidden"}`)

	//Logging under created user (with low perms) and making attempt to create another user
	r.PostForm["login"] = []string{"thrd"}
	r.PostForm["perm_users"] = []string{""}
	a.ResponseUserLogin(&r, w)
	r.Header["Cookie"] = []string{"token=" + string(*w.Data)}
	w.Flush()

	r.PostForm["login"] = []string{"thrd5"}
	r.PostForm["perm_users"] = []string{"1"}
	a.ResponseUserAdd(&r, w)
	w.AssertAndFlush(t, `{"success":false,"code":"403","error":"forbidden"}`)

	//Testing login under non-existent user
	r.PostForm["login"] = []string{"NON_EXITENT_GUY"}
	r.PostForm["pwd"] = []string{"retina-misc1-monstrous-23"}
	r.PostForm["service_id"] = []string{"1"}
	a.ResponseUserLogin(&r, w)
	w.AssertAndFlush(t, `{"success":false,"code":"403","error":"entity_not_exist"}`)

	//Login under existing user but wrong system_id
	r.PostForm["login"] = []string{"SYSTEM_USER"}
	r.PostForm["service_id"] = []string{"2"}
	a.ResponseUserLogin(&r, w)
	w.AssertAndFlush(t, `{"success":false,"code":"403","error":"entity_not_exist"}`)

}
