# golang grpc demo
a demo for basic grpc service
## install packages
```
go mod tidy
```
## start
```
docker-compose up
```

## build
```
SET CGO_ENABLED=0&SET GOOS=linux&SET GOARCH=amd64&go build main.go
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
```