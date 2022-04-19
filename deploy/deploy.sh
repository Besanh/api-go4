#!/bin/sh
echo "Copy environment file"
yes | cp -rf build/goautodial-go-api-env /root/go/env/goautodial-go-api-env
echo "Build go application"
GOOS=linux GOARCH=amd64 go build -o goautodial-go-api main.go
go install
echo "Restart service"
systemctl restart goautodial-go-api