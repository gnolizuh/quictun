#!/bin/bash

LDFLAGS="-s -w"
GCFLAGS=""

OSES=(linux darwin windows freebsd)
ARCHS=(amd64 386)
for os in ${OSES[@]}; do
	for arch in ${ARCHS[@]}; do
		suffix=""
		if [ "$os" == "windows" ]
		then
			suffix=".exe"
		fi
		env CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o quictun_client_${os}_${arch}${suffix} github.com/gnolizuh/quictun/client
		env CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -ldflags "$LDFLAGS" -gcflags "$GCFLAGS" -o quictun_server_${os}_${arch}${suffix} github.com/gnolizuh/quictun/server
		tar -zcf quictun-${os}-${arch}.tar.gz quictun_client_${os}_${arch}${suffix} quictun_server_${os}_${arch}${suffix}
	done
done
