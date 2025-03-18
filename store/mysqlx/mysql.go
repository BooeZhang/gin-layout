package mysqlx

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"gin-layout/config"
	"gin-layout/internal/model"
)

var (
	mysqlDB *gorm.DB
	once    sync.Once
)

func DialToMysql(op *config.MySQL) {
	var dbIns *gorm.DB
	once.Do(func() {
		err := createDB(op)
		if err != nil {
			log.Fatal().Err(err).Msgf("---> [MYSQL] Database %s creation failure", op.Database)
		}
		dsn := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s`,
			op.Username,
			op.Password,
			op.Host,
			op.Database,
			true,
			"Local")

		lv := logger.Info
		if !op.EchoSql {
			lv = logger.Silent
		}

		dbIns, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(lv)})
		if err != nil {
			log.Error().Err(err).Msg("--->[MYSQL] mysql connection open failure")
			os.Exit(1)
		}

		var sqlDB *sql.DB
		sqlDB, err = dbIns.DB()
		if err != nil {
			log.Error().Err(err).Msg("--->[MYSQL] mysql get sql db failure")
			os.Exit(1)

		}
		sqlDB.SetMaxOpenConns(op.MaxOpenConnections)
		sqlDB.SetConnMaxLifetime(op.MaxConnectionLifeTime)
		sqlDB.SetMaxIdleConns(op.MaxIdleConnections)
		err = sqlDB.Ping()
		if err != nil {
			log.Error().Err(err).Msg("---> [MYSQL] mysql connection failure")
			os.Exit(0)
		}

		mysqlDB = dbIns
	})

	if mysqlDB == nil {
		log.Error().Msgf("---> [MYSQL] failed to get mysql store: %+v", mysqlDB)
		os.Exit(1)
	}
}

// createDB 创建数据库
func createDB(opts *config.MySQL) error {
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

func GetDB() *gorm.DB {
	return mysqlDB
}

func Close() {
	if mysqlDB != nil {
		db, err := mysqlDB.DB()
		if err != nil {
			log.Error().Err(err).Msg("---> [MYSQL] close db failure")
			return
		}
		err = db.Close()
		if err != nil {
			log.Error().Err(err).Msg("---> [MYSQL] close db failure")
		}
		log.Info().Msg("---> [MYSQL] close db failure")
	}
}

func CreateSuperUser(db *gorm.DB, cf *config.MySQL) {
	superUser := &model.User{}
	err := db.Where("account = ?", cf.SuperUser).Find(superUser).Error
	if err != nil {
		log.Error().Err(err).Msg("创建超级用户失败")
		return
	}

	if superUser.ID == 0 {
		superUser.Account = cf.SuperUser
		superUser.Password = cf.SuperUserPwd
		superUser.IsActive = true
		superUser.Password, _ = superUser.Encrypt()
		result := db.Create(&superUser)
		if result.Error != nil {
			log.Error().Err(err).Msg("创建超级用户失败")
			os.Exit(1)
		}
	}
}
