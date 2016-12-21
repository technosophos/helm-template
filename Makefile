HELM_HOME ?= $(helm home)
HAS_GLIDE := $(shell command -v glide;)

.PHONY: install
install: bootstrap build
	mkdir -p $(HELM_HOME)/plugins/template
	cp tpl $(HELM_HOME)/plugins/template/
	cp plugin.yaml $(HELM_HOME)/plugins/template/

.PHONY: build
build:
	go build -o tpl ./main.go


.PHONY: bootstrap
bootstrap:
ifndef HAS_GLIDE
	go get -u github.com/Masterminds/glide
endif
	glide install --strip-vendor
