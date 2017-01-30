HELM_HOME ?= $(shell helm home)
HAS_GLIDE := $(shell command -v glide;)

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
	GOOS=linux GOARCH=amd64 go build -o tpl ./main.go
	tar -zcvf helm-template-linux.tgz tpl README.md LICENSE.txt plugin.yaml
	GOOS=darwin GOARCH=amd64 go build -o tpl ./main.go
	tar -zcvf helm-template-macos.tgz tpl README.md LICENSE.txt plugin.yaml


.PHONY: bootstrap
bootstrap:
ifndef HAS_GLIDE
	go get -u github.com/Masterminds/glide
endif
	glide install --strip-vendor
