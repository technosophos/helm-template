HELM_HOME ?= $(shell helm home)
HAS_GLIDE := $(shell command -v glide;)
VERSION := $(shell sed -n -e 's/version:[ "]*\([^"]*\).*/\1/p' plugin.yaml)
DIST := $(CURDIR)/_dist

.PHONY: install
install: bootstrap build
	mkdir -p $(HELM_HOME)/plugins/template
	cp tpl $(HELM_HOME)/plugins/template/
	cp plugin.yaml $(HELM_HOME)/plugins/template/

.PHONY: build
build:
	go build -o tpl ./main.go

.PHONY: dist
dist:
	mkdir -p $(DIST)
	GOOS=linux GOARCH=amd64 go build -o tpl ./main.go
	tar -zcvf $(DIST)/helm-template-linux-$(VERSION).tgz tpl README.md LICENSE.txt plugin.yaml
	GOOS=darwin GOARCH=amd64 go build -o tpl ./main.go
	tar -zcvf $(DIST)/helm-template-macos-$(VERSION).tgz tpl README.md LICENSE.txt plugin.yaml


.PHONY: bootstrap
bootstrap:
ifndef HAS_GLIDE
	go get -u github.com/Masterminds/glide
endif
	glide install --strip-vendor
