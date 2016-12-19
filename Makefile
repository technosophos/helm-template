HELM_HOME ?= $(helm home)

.PHONY: install
install: build
	cp tpl $(HELM_HOME)/plugins/template/
	cp plugin.yaml $(HELM_HOME)/plugins/template/

.PHONY: build
build:
	go build -o tpl ./main.go

.PHONY: bootstrap
bootstrap:
	glide install --strip-vendor
