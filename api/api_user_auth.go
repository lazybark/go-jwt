package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/lazybark/go-helpers/cli/clf"
	"github.com/lazybark/go-jwt/storage"
	"golang.org/x/crypto/bcrypt"
)

func (a *Api) GenerateHMACToken(u storage.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo": "bar",
		"iss": storage.Unversal.String(),
		"sub": u.ServiceId.String(),
		"iat": jwt.NewNumericDate(time.Now()),
		"nbf": jwt.NewNumericDate(time.Now()),
		"exp": jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	})

	tokenString, err := token.SignedString([]byte(a.conf.Secret))

	return tokenString, err
}

// UserAdd call to storage UserAdd method, checks results and writes []byte answer to client
func (a *Api) ResponseUserLogin(req *http.Request, w http.ResponseWriter) {
	//Check if user data fits our requirements
	login := req.Form.Get("login")
	pwd := req.Form.Get("pwd")
	if login == "" || pwd == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(ApiError, ErrorBadFormCode, ErrorBadForm)))
		return
	}

	userData, err := a.db.UserGetData(login)
	if err != nil {
		fmt.Println(clf.Red(err))
		if err == storage.ErrEntityNotExist {
			w.Write([]byte(fmt.Sprintf(ApiError, ErrorNotExistCode, ErrorNotExist)))
		}
		w.Write([]byte(fmt.Sprintf(ApiError, ErrorInternalCode, ErrorInternal)))
		return
	}

	//Check if password is correct
	correct, err := a.ComparePasswords(userData.PasswordHash, pwd)
	if !correct {
		w.Write([]byte(fmt.Sprintf(ApiError, ErrorWrongCredsCode, ErrorWrongCreds)))
		return
	}
	if err != nil {
		fmt.Println(clf.Red(err))
		w.Write([]byte(fmt.Sprintf(ApiError, ErrorInternalCode, ErrorInternal)))
		return
	}

	token, err := a.GenerateHMACToken(userData)
	if err != nil {
		fmt.Println(clf.Red(err))
		w.Write([]byte(fmt.Sprintf(ApiError, ErrorInternalCode, ErrorInternal)))
		return
	}

	w.Write([]byte(token))

}

func (a *Api) ComparePasswords(hashedPwd string, plainPwd string) (bool, error) {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	plainPwdBytes := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwdBytes)
	if err != nil {
		log.Println(err)
		return false, err
	}

	return true, nil
}
