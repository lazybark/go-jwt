package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/lazybark/go-helpers/cli/clf"
	"github.com/lazybark/go-jwt/storage"
	"golang.org/x/crypto/bcrypt"
)

func (a *Api) GenerateHMACToken(u storage.User) (string, error) {
	claims := &JWTClaims{
		ID:          u.ID,
		Login:       u.Login,
		Name:        u.Name,
		LastName:    u.LastName,
		BytePayload: string(u.BytePayload),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    storage.ServiceUnversal.String(),
			Subject:   u.ServiceId.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(a.conf.Secret))

	return tokenString, err
}

// UserAdd call to storage UserAdd method, checks results and writes []byte answer to client
func (a *Api) ResponseUserLogin(req *http.Request, w http.ResponseWriter) {
	//Check if user data fits our requirements
	login := req.PostForm.Get("login")
	pwd := req.PostForm.Get("pwd")
	srv := req.PostForm.Get("service_id")
	if login == "" || pwd == "" || srv == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(ApiError, ErrorBadFormCode, ErrorBadForm)))
		return
	}

	serviceId, err := strconv.Atoi(srv)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(ApiError, ErrorBadFormCode, ErrorBadForm)))
		return
	}

	userData, err := a.db.UserGetData(login, serviceId)
	if err != nil {
		if err == storage.ErrEntityNotExist {
			w.Write([]byte(fmt.Sprintf(ApiError, ErrorNotExistCode, ErrorNotExist)))
			return
		}
		fmt.Println(clf.Red(err))
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

	err = a.db.UserUpdateActivity(userData.ID)
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
		return false, err
	}

	return true, nil
}

func (a *Api) CheckUsersControlPermission(uid string, perm storage.PermissionUsers, targPerm storage.PermissionUsers) (bool, error) {
	str, err := a.db.UserGetParam(uid, "permission_users")
	if err != nil {
		return false, err
	}
	if str == "" {
		return false, err
	}
	p, err := strconv.Atoi(str)
	if err != nil {
		return false, err
	}
	up := storage.PermissionUsers(p)
	//Check if user has the right permisson at least
	if !perm.Check(up) {
		return false, nil
	}
	//And if the desired permission is not higher that user's
	return targPerm.Check(up), nil
}

func (a *Api) CheckUsersControlService(uid string, srv int) (bool, error) {
	str, err := a.db.UserGetParam(uid, "service_id")
	if err != nil {
		return false, err
	}
	if str == "" {
		return false, err
	}
	p, err := strconv.Atoi(str)
	if err != nil {
		return false, err
	}

	//User can control only its service users (except 1, 1 is the universal service controller)
	fmt.Println(p == 1 || p == srv)
	return p == 1 || p == srv, nil
}
