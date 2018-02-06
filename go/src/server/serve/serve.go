package serve

import (
	"crypto/tls"
	"local"
	"log"
	"net/http"
	"os"
	"reflect"
	"server/handlers"
	"server/middleware"
	"time"
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

		certfile := os.Getenv("CERTFILE")
		keyfile := os.Getenv("KEYFILE")

		// Reload the certificate daily to handle renewals without a restart.
		// Note this can cause a broken connection if the reload occurs during handshake.
		cert, err := tls.LoadX509KeyPair(certfile, keyfile)
		if err != nil {
			log.Fatalf("Error loading certificate: %v", err)
		}
		certLoaded := time.Now()
		getCertificate := func(helloInfo *tls.ClientHelloInfo) (*tls.Certificate, error) {
			if time.Since(certLoaded).Hours() > 24 {
				log.Println("Checking for new certificate")
				oldCert := cert
				cert, err = tls.LoadX509KeyPair(certfile, keyfile)
				if err != nil {
					log.Fatalf("Error reloading certificate: %v", err)
				}

				if reflect.DeepEqual(oldCert.Certificate, cert.Certificate) {
					log.Println("Certificate unchanged")
				} else {
					log.Println("Certificate reloaded")
				}
				certLoaded = time.Now()
			}
			return &cert, err
		}

		server := &http.Server{
			Addr:    port_https,
			Handler: mux,
			TLSConfig: &tls.Config{
				GetCertificate: getCertificate,
			},
		}

		// Call ListenAndServerTLS without certificate paths because we use getCertificate.
		err = server.ListenAndServeTLS("", "")
	} else {
		err = http.ListenAndServe(port_http, mux)
	}

	if err != nil {
		log.Fatal(err)
	}
}
