proto:
	goctl rpc protoc ./user.proto --go_out=. --go-grpc_out=. --zrpc_out=./
model:
	goctl model mysql datasource --dir ./internal/model --table user --cache true --url "root:123456@tcp(127.0.0.1:3306)/IZH_user"
user:
	go run user.go
.PHONY:proto user model