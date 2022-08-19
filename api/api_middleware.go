package api

import (
	"fmt"
	"net/http"
)

func parseFormMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf(ApiError, ErrorBadFormCode, ErrorBadForm)))
			return
		}
		next.ServeHTTP(w, r)
	})
}
