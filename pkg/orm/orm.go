package orm

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	*gorm.DB
}

type Config struct {
	DSN          string //数据库的url
	MaxOpenConns int
	MaxIdleConns int
	MaxLifeTime  int
}

// ormLog need to implement the methods of gorm log
type ormLog struct {
	LogLevel logger.LogLevel
}

func (l *ormLog) LogMode(level logger.LogLevel) logger.Interface {
	l.LogLevel = level
	return l
}
func (l *ormLog) Info(ctx context.Context, format string, v ...interface{}) {
	if l.LogLevel < logger.Info {
		return
	}
	logx.WithContext(ctx).Infof(format, v...)
}
func (l *ormLog) Warn(ctx context.Context, format string, v ...interface{}) {
	if l.LogLevel < logger.Warn {
		return
	}
	logx.WithContext(ctx).Infof(format, v...)
}
func (l *ormLog) Error(ctx context.Context, format string, v ...interface{}) {
	if l.LogLevel < logger.Error {
		return
	}
	logx.WithContext(ctx).Errorf(format, v...)
}
func (l *ormLog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	logx.WithContext(ctx).WithDuration(elapsed).Infof("[%.3fms][rows:%v] %s", float64(elapsed.Nanoseconds())/1e6, rows, sql)
}

func NewMysql(config *Config) (*DB, error) {
	if config.MaxOpenConns == 0 {
		config.MaxOpenConns = 100
	}
	if config.MaxIdleConns == 0 {
		config.MaxIdleConns = 10
	}
	if config.MaxLifeTime == 0 {
		config.MaxLifeTime = 3600
	}

	db, err := gorm.Open(mysql.Open(config.DSN), &gorm.Config{
		Logger: &ormLog{},
	})
	if err != nil {
		logx.Errorf("gorm Open :%s error :%v", config.DSN, err)
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	//设置相关的最大连接数、maxlifetime
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(config.MaxLifeTime))

	//使用Hook，将metrics 和 trace能力添加到gorm中（埋点的过程）

	err = db.Use(NewCustomPlugin())
	if err != nil {
		//fmt.Println("db use plugin failed")
		return nil, err
	}

	return &DB{db}, nil
}

func MustNewMysql(config *Config) *DB {
	db, err := NewMysql(config)
	if err != nil {
		panic(err)
	}
	return db
}
