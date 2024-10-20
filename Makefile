all: local/server local/static local

STATIC_DEPS:=

COMPRESSIBLE=$(wildcard local/static/*.html) \
			 $(wildcard local/static/*.css) \
			 $(wildcard local/static/*.js) \
			 $(wildcard local/static/**/*.html) \
			 $(wildcard local/static/**/*.css) \
			 $(wildcard local/static/**/*.js)
.PHONY: gz
gz: $(addsuffix .gz,$(COMPRESSIBLE))

%.gz: %
	gzip -fk9 $<

include local/Makefile

GOSRCS=$(wildcard local/**/*.go) $(wildcard local/plugins/**/*.go)
space=$(eval) $(eval)
GOPATHS=$(subst $(space),:,$(abspath go local/go $(wildcard local/plugins/*/go)))
local/server: $(GOSRCS)
	GOPATH="$$GOPATH:$(GOPATHS)" CGO_ENABLED=0 cd server; go build -o ../local/server

gotest: $(GOSRCS)
	cd server
	GOPATH="$$GOPATH:$(GOPATHS)" go test -v ./...

.PHONY: local/static
local/static: $(STATIC_DEPS)
	mkdir -p local/static
	python compile_templates.py
