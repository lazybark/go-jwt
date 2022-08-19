package api

import (
	"fmt"
	"net/http"

	"github.com/lazybark/go-jwt/storage"
)

// UserAdd call to storage UserAdd method, checks results and writes []byte answer to client
func (a *Api) ResponseUserAdd(req *http.Request, w http.ResponseWriter) {
	//Check if user data fits our requirements
	login := req.PostForm.Get("login")
	pwd := req.PostForm.Get("pwd")
	if login == "" || pwd == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(ApiError, ErrorBadFormCode, ErrorBadForm)))
		return
	}

	newUser := storage.User{
		Login:        login,
		PasswordHash: pwd,
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
