newmigrate:
	migrate create -ext sql -dir db/migration/$(db)_migration -seq $(name)
migrateup:
	migrate -path db/migration/$(db)_migration -database "mysql://root:123456@tcp(127.0.0.1:3306)/IZH_$(db)?parseTime=true" -verbose up
migratedown:
	migrate -path db/migration/$(db)_migration -database "mysql://root:123456@tcp(127.0.0.1:3306)/IZH_$(db)?parseTime=true" -verbose down
test:
	go test ./...
user:
	cd application/user/rpc;go run user.go;
applet:
	cd application/applet;go run applet.go;
article-applet:
	cd application/article/api;go run applet.go;
article:
	cd application/article/rpc;go run article.go;
startContainers:
	docker start mysql redis etcd
stopContainers:
	docker stop mysql redis etcd
userModelMock:
	mockgen -package mockuserModel -destination application/user/rpc/internal/mock/user_model.go github.com/GGjahoon/IZH/application/user/rpc/internal/model UserModel
proto:
	cd $(path);goctl rpc protoc ./$(name).proto --go_out=. --go-grpc_out=. --zrpc_out=./
api:
	cd $(path);goctl api go --dir=./ --api $(name).api
cachemodel:
	cd $(path);goctl model mysql datasource --dir ./internal/model --table $(table) --cache true --url "root:123456@tcp(127.0.0.1:3306)/IZH_$(table)"
model:
	cd $(path);goctl model mysql datasource --dir ./internal/model --table $(table) --url "root:123456@tcp(127.0.0.1:3306)/IZH_$(table)"
.PHONY:newmigrate migrateup migratedown test user applet startContainers stopContainers proto api cachemodel model