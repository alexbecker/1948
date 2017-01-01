package local

import (
	"net/http"
)

const Index = "/index.html"

var AdditionalRegistrations = []func(*http.ServeMux){}
