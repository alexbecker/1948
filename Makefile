all: go/bin/server static local

STATIC_DEPS:=

COMPRESSIBLE=$(wildcard static/*.html) $(wildcard static/*.css) $(wildcard static/*.js) $(wildcard static/**/*.html) $(wildcard static/**/*.css) $(wildcard static/**/*.js)
.PHONY: gz
gz: $(addsuffix .gz,$(COMPRESSIBLE))

%.gz: %
	gzip -fk9 $<

include local/Makefile

GOSRCS=go/src/server/main.go $(wildcard go/src/server/**/*.go) $(wildcard local/go/src/**/*.go) $(wildcard local/plugins/*/go/src/**/*.go)
space=$(eval) $(eval)
GOPATHS=$(subst $(space),:,$(abspath go local/go $(wildcard local/plugins/*/go)))
go/bin/server: $(GOSRCS)
	GOPATH="$$GOPATH:$(GOPATHS)" go install -v server

gotest: $(GOSRCS)
	GOPATH="$$GOPATH:$(GOPATHS)" go test -v ./...

.PHONY: static
static: $(STATIC_DEPS)
	mkdir -p static
	python compile_templates.py
