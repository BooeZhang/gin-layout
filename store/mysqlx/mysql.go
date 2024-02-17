package mysqlx

import (
	"database/sql"
	"fmt"
	"github.com/BooeZhang/gin-layout/config"
	"github.com/BooeZhang/gin-layout/internal/model"
	"github.com/BooeZhang/gin-layout/pkg/log"
	"github.com/BooeZhang/gin-layout/pkg/log/sqlhook"
	"os"
	"sync"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	mysqlDB *gorm.DB
	once    sync.Once
)

func GetDB() *gorm.DB {
	return mysqlDB
}

// Close 关闭数据库
func Close() error {
	db, err := mysqlDB.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

// InitMysql 初始化 mysql
func InitMysql(cf *config.MySQLConfig) {
	DialToMysql(cf)
}

// DialToMysql 连接 mysql
func DialToMysql(op *config.MySQLConfig) {
	var err error
	var dbIns *gorm.DB
	once.Do(func() {
		err := createDB(op)
		if err != nil {
			log.Errorf("Database %s creation failure", op.Database)
			os.Exit(1)
		}
		dsn := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s`,
			op.Username,
			op.Password,
			op.Host,
			op.Database,
			true,
			"Local")
		dbIns, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: sqlhook.New(op.LogLevel),
		})

		var sqlDB *sql.DB

		sqlDB, err = dbIns.DB()

		// SetMaxOpenConns sets the maximum number of open connections to the database.
		sqlDB.SetMaxOpenConns(op.MaxOpenConnections)

		// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
		sqlDB.SetConnMaxLifetime(op.MaxConnectionLifeTime)

		// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
		sqlDB.SetMaxIdleConns(op.MaxIdleConnections)
		err = sqlDB.Ping()
		if err != nil {
			log.Error("MySQL启动异常", zap.Error(err))
			os.Exit(0)
		}

		mysqlDB = dbIns
	})

	if mysqlDB == nil || err != nil {
		log.Errorf("mysqlx connection failure, error: %s", err)
		os.Exit(1)
	}
}

// createDB 创建数据库
func createDB(opts *config.MySQLConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/", opts.Username, opts.Password, opts.Host)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(db)

	createSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;", opts.Database)
	if err = db.Ping(); err != nil {
		return err
	}
	re, err := db.Exec(createSql)
	_, err = re.RowsAffected()
	return err
}

// CreateSuperUser 创建超级用户
func CreateSuperUser(db *gorm.DB, cf *config.MySQLConfig) {
	superUser := &model.SysUser{}
	err := db.Where("account = ?", cf.SuperUser).Find(superUser).Error
	if err != nil {
		log.Errorf("创建超级用户失败：%s", err)
		os.Exit(1)
	}

	if superUser.ID == 0 {
		superUser.Account = cf.SuperUser
		superUser.Password = cf.SuperUserPwd
		superUser.IsActive = true
		superUser.Password, _ = superUser.Encrypt()
		result := db.Create(&superUser)
		if result.Error != nil {
			log.Errorf("创建超级用户失败：%s", err)
			os.Exit(1)
		}
	}
}
