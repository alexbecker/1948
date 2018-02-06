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

GOSRCS=go/src/server/main.go $(wildcard go/src/server/**/*.go) $(wildcard local/go/src/**/*.go) $(wildcard local/plugins/*/go/src/**/*.go)
space=$(eval) $(eval)
GOPATHS=$(subst $(space),:,$(abspath go local/go $(wildcard local/plugins/*/go)))
local/server: $(GOSRCS)
	GOPATH="$$GOPATH:$(GOPATHS)" go install -v server
	cp go/bin/server local/server

gotest: $(GOSRCS)
	GOPATH="$$GOPATH:$(GOPATHS)" go test -v ./...

.PHONY: static
local/static: $(STATIC_DEPS)
	mkdir -p local/static
	python compile_templates.py
