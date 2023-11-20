newmigrate:
	migrate create -ext sql -dir db/migration/$(path) -seq $(name)
usermigrateup:
	migrate -path db/migration/user_migration -database "mysql://root:123456@tcp(127.0.0.1:3306)/IZH_user?parseTime=true" -verbose up
usermigratedown:
	migrate -path db/migration/user_migration -database "mysql://root:123456@tcp(127.0.0.1:3306)/IZH_user?parseTime=true" -verbose down
test:
	go test ./...
.PHONY:newmigrate usermigrateup usermigratedown