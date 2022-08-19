package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
	"github.com/lazybark/go-helpers/cli/clf"
	"github.com/lazybark/go-jwt/storage"
)

// UserAdd call to storage UserAdd method, checks results and writes []byte answer to client
func (a *Api) ResponseUserAdd(req *http.Request, w http.ResponseWriter) {
	//Parse user jwt token
	cookie, err := req.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf(ApiError, ErrorUnauthedCode, ErrorUnauthed)))
			return
		}
		fmt.Println(clf.Red(err))
		w.Write([]byte(fmt.Sprintf(ApiError, ErrorBadRequestCode, ErrorBadRequest)))
		return
	}

	claims := &JWTClaims{}

	tkn, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.conf.Secret), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf(ApiError, ErrorUnauthedCode, ErrorUnauthed)))
			return
		}
		fmt.Println(clf.Red(err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(ApiError, ErrorBadRequestCode, ErrorBadRequest)))
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf(ApiError, ErrorUnauthedCode, ErrorUnauthed)))
		return
	}

	//Now check perms
	ok, err := a.CheckUsersControlPermission(fmt.Sprint(claims.ID), storage.UsersCreate)
	if err != nil || !ok {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(fmt.Sprintf(ApiError, ErrorForbiddenCode, ErrorForbidden)))
		return
	}

	//Check if user data fits our requirements
	login := req.PostForm.Get("login")
	pwd := req.PostForm.Get("pwd")
	name := req.PostForm.Get("name")
	lastName := req.PostForm.Get("last_name")
	email := req.PostForm.Get("email")
	permUsers := req.PostForm.Get("perm_users")
	service := req.PostForm.Get("service_id")
	if login == "" || pwd == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(ApiError, ErrorBadFormCode, ErrorBadForm)))
		return
	}

	permUsersId := 0
	if permUsers != "" {
		permUsersId, err = strconv.Atoi(permUsers)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf(ApiError, ErrorBadFormCode, ErrorBadForm)))
			return
		}
	}

	serviceId := 0
	if service != "" {
		serviceId, err = strconv.Atoi(service)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf(ApiError, ErrorBadFormCode, ErrorBadForm)))
			return
		}
	}

	newUser := storage.User{
		Login:           login,
		PasswordHash:    pwd,
		Name:            name,
		LastName:        lastName,
		Email:           email,
		ServiceId:       storage.Service(serviceId),
		PermissionUsers: storage.PermissionUsers(permUsersId),
	}

	uid, err := a.db.UserAdd(newUser)
	if err != nil {
		if err == storage.ErrEntityExists {
			w.Write([]byte(fmt.Sprintf(ApiError, ErrorBadFormCode, storage.ErrEntityExists)))
			return
		}
	}
	w.Write([]byte(fmt.Sprintf(ApiStringResult, uid)))
}
