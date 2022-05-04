package mysql

import (
	"database/sql"
	"fmt"
	"github.com/BooeZhang/gin-layout/internal/pkg/options"
	"github.com/BooeZhang/gin-layout/pkg/log"
	"github.com/BooeZhang/gin-layout/pkg/logger"
	"github.com/BooeZhang/gin-layout/store"
	"os"
	"sync"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type datastore struct {
	db *gorm.DB
}

func (ds *datastore) Close() error {
	_db, err := ds.db.DB()
	if err != nil {
		return fmt.Errorf("get gorm db instance failed: %w", err)
	}
	return _db.Close()
}

func (ds *datastore) GetDB() *gorm.DB {
	return ds.db
}

var (
	mysqlFactory store.Factory
	once         sync.Once
)

// GetMysqlFactoryOr 使用给定的配置创建 mysql 工厂。
func GetMysqlFactoryOr(opts *options.MySQLOptions) (store.Factory, error) {
	if opts == nil && mysqlFactory == nil {
		return nil, fmt.Errorf("failed to get mysql store fatory")
	}

	var err error
	var dbIns *gorm.DB
	once.Do(func() {
		dsn := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s`,
			opts.Username,
			opts.Password,
			opts.Host,
			opts.Database,
			true,
			"Local")
		dbIns, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.New(opts.LogLevel),
		})

		var sqlDB *sql.DB

		sqlDB, err = dbIns.DB()

		// SetMaxOpenConns sets the maximum number of open connections to the database.
		sqlDB.SetMaxOpenConns(opts.MaxOpenConnections)

		// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
		sqlDB.SetConnMaxLifetime(opts.MaxConnectionLifeTime)

		// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
		sqlDB.SetMaxIdleConns(opts.MaxIdleConnections)
		err = sqlDB.Ping()
		if err != nil {
			log.Error("MySQL启动异常", zap.Error(err))
			os.Exit(0)
		}

		mysqlFactory = &datastore{dbIns}
	})

	if mysqlFactory == nil || err != nil {
		return nil, fmt.Errorf("failed to get mysql store fatory, mysqlFactory: %+v, error: %w", mysqlFactory, err)
	}
	err = migrateDatabase(dbIns)
	if err != nil {
		log.Error(err.Error())
	}
	return mysqlFactory, nil
}

func GetMysqlFactory() store.Factory {
	return mysqlFactory
}
