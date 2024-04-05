TEST?=$$(go list ./... | grep -v 'vendor'|grep -v 'examples')
TESTTIMEOUT=180m
HOSTNAME=lumen.com
NAMESPACE=lumentech
NAME=lumen
BINARY=terraform-provider-${NAME}
VERSION=2.2.0
OS=$(shell go env GOOS)
OS_ARCH=$(shell go env GOARCH)
INSTALL_PATH=~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/linux_$(OS_ARCH)

ifeq ($(OS), darwin)
	INSTALL_PATH=~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/darwin_$(OS_ARCH)
endif

ifeq ($(OS), "windows")
	INSTALL_PATH=%APPDATA%/Terraform/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/windows_$(OS_ARCH)
endif

default: install

tools:
	go get github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource@v2.7.0

build: tools
	go build -o ${BINARY}

release:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64
	GOOS=freebsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_freebsd_386
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_freebsd_amd64
	GOOS=freebsd GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_freebsd_arm
	GOOS=linux GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_linux_386
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64
	GOOS=linux GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_linux_arm
	GOOS=openbsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_openbsd_386
	GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_openbsd_amd64
	GOOS=solaris GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_solaris_amd64
	GOOS=windows GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_windows_386
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64

install: build setup
	mkdir -p $(INSTALL_PATH)
	mv ${BINARY} $(INSTALL_PATH)

fmt:
	echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./lumen

fmtcheck:
	bash -x scripts/gofmtcheck.sh

test: fmtcheck
	go test -i $(TEST) || exit 1                                                   
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    

testacc:
	echo "==> Running acceptance tests..."
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

clean:
	bash -x scripts/cleanup.sh

setup:
	bash -x scripts/setup.sh
