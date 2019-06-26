CURDIR=$(pwd)
OLDGOPATH="$GOPATH"
export GOPATH="$CURDIR"

gofmt -w ./src

go install -race main

export GOPATH="$OLDGOPATH"
