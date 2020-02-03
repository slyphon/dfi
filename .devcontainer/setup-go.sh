#!/bin/bash

# Install gocode-gomod
go get -x -d github.com/stamblerre/gocode 2>&1
go build -o gocode-gomod github.com/stamblerre/gocode
mv gocode-gomod $GOPATH/bin/

gguv() {
  go get -u -v "$@"
}

gguv github.com/mdempsky/gocode
gguv github.com/uudashr/gopkgs/cmd/gopkgs
gguv github.com/ramya-rao-a/go-outline
gguv github.com/acroca/go-symbols
gguv github.com/godoctor/godoctor
gguv golang.org/x/tools/cmd/guru
gguv golang.org/x/tools/cmd/gorename
gguv github.com/rogpeppe/godef
gguv github.com/zmb3/gogetdoc
gguv github.com/haya14busa/goplay/cmd/goplay
gguv github.com/sqs/goreturns
gguv github.com/josharian/impl
gguv github.com/davidrjenni/reftools/cmd/fillstruct
gguv github.com/fatih/gomodifytags
gguv github.com/cweill/gotests/...
gguv golang.org/x/tools/cmd/goimports
gguv golang.org/x/lint/golint
gguv golang.org/x/tools/gopls@latest
gguv honnef.co/go/tools/...
gguv github.com/golangci/golangci-lint/cmd/golangci-lint
gguv github.com/mgechev/revive
pushd /tmp
gguv github.com/derekparker/delve/cmd/dlv
popd
