package model

import (
	"flag"
	"os"
	"testing"

	"github.com/GGjahoon/IZH/application/user/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var testModel UserModel
var configFile = flag.String("f", "user.yaml", "the config file")

func TestMain(m *testing.M) {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	//create a connection to mysql
	sqlConn := sqlx.NewMysql(c.DataSource)
	testModel = NewUserModel(sqlConn, c.CacheRedis)
	os.Exit(m.Run())
}
