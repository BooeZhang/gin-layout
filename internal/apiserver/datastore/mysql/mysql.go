package mysql

import (
	"database/sql"
	"fmt"
	"github.com/BooeZhang/gin-layout/internal/apiserver/datastore"
	"github.com/BooeZhang/gin-layout/internal/apiserver/datastore/datainterface"
	"github.com/BooeZhang/gin-layout/pkg/log/sqlhook"
	"os"
	"sync"

	"github.com/BooeZhang/gin-layout/internal/pkg/config"
	"github.com/BooeZhang/gin-layout/pkg/log"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	mysqlFactory *_datastore
	once         sync.Once
)

type _datastore struct {
	db *gorm.DB
}

func (ds *_datastore) GetDB() *gorm.DB {
	return ds.db
}

func (ds *_datastore) SysUser() datainterface.SysUserData {
	return newSysUser(ds)
}

// Close 关闭数据库
func (ds *_datastore) Close() error {
	db, err := ds.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

// GetMysqlFactoryOr 使用给定的配置创建 mysql 工厂。
func GetMysqlFactoryOr(opts *config.MySQLConfig) (datastore.Factory, error) {
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
			Logger: sqlhook.New(opts.LogLevel),
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

		mysqlFactory = &_datastore{dbIns}
	})

	if mysqlFactory == nil || err != nil {
		return nil, fmt.Errorf("failed to get mysql store fatory, mysqlFactory: %+v, error: %w", mysqlFactory, err)
	}

	return mysqlFactory, nil
}

