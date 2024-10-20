package handlers

import (
	"log"
	"net/http"
)

func RedirectHTTPS(w http.ResponseWriter, req *http.Request) {
	dest := req.URL
	dest.Scheme = "https"
	dest.Host = req.Host
	http.Redirect(w, req, dest.String(), http.StatusMovedPermanently)
}

func AlwaysRedirectHTTPS(port string) {
	err := http.ListenAndServe(port, http.HandlerFunc(RedirectHTTPS))
	if err != nil {
		log.Fatal(err)
	}
}
