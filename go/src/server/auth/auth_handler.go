package auth

import (
	"net/http"
	"server/middleware"
)

func challenge(w http.ResponseWriter, role string) {
	w.Header().Set("WWW-Authenticate", "Basic realm=\""+role+"\"")
	w.WriteHeader(http.StatusUnauthorized)
}

func AuthHandler(role string, handler http.Handler) http.Handler {
	return middleware.LoggedHandler(func(w http.ResponseWriter, req *http.Request) error {
		username, err := CheckAuthCookie(req)
		if err != nil {
			username, password, ok := req.BasicAuth()
			if !ok {
				challenge(w, role)
				return nil
			}

			ok, err = CheckPassword(username, password)
			if err != nil {
				return err
			}
			if !ok {
				challenge(w, role)
				return nil
			}

			SetAuthCookie(w, username)
			handler.ServeHTTP(w, req)
			return nil
		}

		roles, err := GetRoles(username)
		if err != nil {
			return err
		}

		if !roles[role] {
			http.Error(w, "403 forbidden", http.StatusForbidden)
			return nil
		}

		handler.ServeHTTP(w, req)
		return nil
	})
}
