newmigrate:
	migrate create -ext sql -dir db/migration/$(path) -seq $(name)
usermigrateup:
	migrate -path db/migration/user_migration -database "mysql://root:123456@tcp(127.0.0.1:3306)/IZH_user?parseTime=true" -verbose up
usermigratedown:
	migrate -path db/migration/user_migration -database "mysql://root:123456@tcp(127.0.0.1:3306)/IZH_user?parseTime=true" -verbose down
statusproto:
	protoc --proto_path=pkg/xcode/types pkg/xcode/types/status.proto --go_out=pkg/xcode/types --go_opt=paths=source_relative
test:
	go test ./...
user:
	cd application/user/rpc;go run user.go;
applet:
	cd application/applet;go run applet.go;
startContainers:
	docker start mysql redis etcd
stopContainers:
	docker stop mysql redis etcd
userModelMock:
	mockgen -package mockuserModel -destination application/user/rpc/internal/mock/user_model.go github.com/GGjahoon/IZH/application/user/rpc/internal/model UserModel
.PHONY:newmigrate usermigrateup usermigratedown user startContainers stopContainers 