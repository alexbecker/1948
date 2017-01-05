package serve

import (
	"local"
	"log"
	"net/http"
	"os"
	"server/handlers"
	"server/middleware"
)

func makeServeMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/" || (req.URL.Path == "/index.html" && local.Index != "/index.html") {
			req.URL.Path = local.Index
			mux.ServeHTTP(w, req)
		} else {
			handlers.ServeEncFiles(w, req, "static/")
		}
	})

	for _, register := range local.AdditionalRegistrations {
		register(mux)
	}

	return middleware.MiddlewareServeMux(mux)
}

func Serve() {
	mux := makeServeMux()

	port_http := ":" + os.Getenv("PORT_HTTP")
	port_https := ":" + os.Getenv("PORT_HTTPS")
	var err error
	if port_https != ":" {
		// spawn another process to listen over HTTP
		go handlers.AlwaysRedirectHTTPS(port_http)

		err = http.ListenAndServeTLS(port_https, os.Getenv("CERTFILE"), os.Getenv("KEYFILE"), mux)
	} else {
		err = http.ListenAndServe(port_http, mux)
	}

	if err != nil {
		log.Fatal(err)
	}
}
