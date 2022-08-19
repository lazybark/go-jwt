package api

import "github.com/golang-jwt/jwt/v4"

type ApiAnswer struct {
	Success bool        `json:"success"`
	Result  interface{} `json:"result"`
}

type JWTClaims struct {
	Login     string `json:"login"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	LastName  string `json:"last_name"`
	UsersPerm string `json:"users_permission"`
	ServiceID int    `json:"service_id"`
	jwt.RegisteredClaims
}

var (
	//ApiError represents error with status, code and human-readable error
	ApiError = `{"success":false,"code":"%v","error":"%s"}`
	//ApiStringResult represents one-string result.
	//Useful in case of one-string result from a method
	ApiStringResult = `{"success":true,"result":"%v"}`

	ErrorInternal     = "err_internal"
	ErrorInternalCode = 500

	ErrorBadRequest     = "bad_request"
	ErrorBadRequestCode = 400

	ErrorBadForm     = "bad_form_fields"
	ErrorBadFormCode = 400

	ErrorExists     = "entity_exists"
	ErrorExistsCode = 200

	ErrorNotExist     = "entity_not_exist"
	ErrorNotExistCode = 403

	ErrorWrongCreds     = "wrong_credentials"
	ErrorWrongCredsCode = 403

	ErrorUnauthed     = "unathorized"
	ErrorUnauthedCode = 401
)
