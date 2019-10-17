#!/bin/bash

set -e -u

# install dependencies for coverage if missing
if [[ -z "$(which 'goveralls' 2>/dev/null)" ]]; then
  go get -v github.com/mattn/goveralls
fi

coverprofile=''
for pkg in $(go list ./...); do
    covfile="$(echo $pkg | sed -e 's/github.com\/kamichidu\///' | sed -e 's/\//_/g').cov"
    go test -cover -covermode count -coverprofile "$covfile" "$pkg"
    if [[ -z "$coverprofile" ]]; then
        coverprofile="${covfile}"
    else
        coverprofile="${coverprofile},${covfile}"
    fi
done
goveralls -coverprofile "$coverprofile"
