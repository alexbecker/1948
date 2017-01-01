package handlers

import (
	"server/auth"
	"server/middleware"
	"net/http"
)

func challenge(w http.ResponseWriter, role string) {
	w.Header().Set("WWW-Authenticate", "Basic realm=\""+role+"\"")
	w.WriteHeader(http.StatusUnauthorized)
}

func AuthHandler(role string, handler http.Handler) http.Handler {
	return middleware.LoggedHandler(func(w http.ResponseWriter, req *http.Request) error {
		username, err := auth.CheckAuthCookie(req)
		if err != nil {
			username, password, ok := req.BasicAuth()
			if !ok {
				challenge(w, role)
				return nil
			}

			ok, err = auth.CheckPassword(username, password)
			if err != nil {
				return err
			}
			if !ok {
				challenge(w, role)
				return nil
			}

			auth.SetAuthCookie(w, username)
			handler.ServeHTTP(w, req)
			return nil
		}

		roles, err := auth.GetRoles(username)
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
