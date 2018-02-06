#!/bin/sh
go build -ldflags "-s -w" -o dist/scron scron.go
if [ $? -eq 0 ];then
    rm -rf dist/scron-upx
    tools/upx -9 -o dist/scron-upx dist/scron
fi
