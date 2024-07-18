package core

import (
	"context"
	"database/sql"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"

	"github.com/BooeZhang/gin-layout/config"
	"github.com/BooeZhang/gin-layout/store/mongodb"
	"github.com/BooeZhang/gin-layout/store/mysqlx"
	"github.com/BooeZhang/gin-layout/store/redisx"
)

type StoreImpl struct {
	mysql  *gorm.DB
	mgo    *mongo.Database
	rds    redis.UniversalClient
	c      config.Config
	closed bool
}

func NewStorageWithConfig(c config.Config) *StoreImpl {
	// mongodb.InitMongo(c.MongoConfig)
	redisx.InitRedis(c.RedisConfig)
	mysqlx.InitMysql(c.MySQLConfig)

	return &StoreImpl{
		mysql: mysqlx.GetDB(),
		mgo:   mongodb.GetDBSession(),
		rds:   redisx.GetRedis(),
	}
}

func (s *StoreImpl) GetMySQL() *gorm.DB {
	return s.mysql
}

func (s *StoreImpl) GetMongo() *mongo.Database {
	return s.mgo
}

func (s *StoreImpl) GetRedis() redis.UniversalClient {
	return s.rds
}

func (s *StoreImpl) WithMongoTransaction(ctx context.Context, fn func(ctx mongo.SessionContext) (interface{}, error),
	opts ...*options.TransactionOptions) (interface{}, error) {
	cli := mongodb.GetClient()
	session, err := cli.StartSession()
	if err != nil {
		return nil, err
	}
	// Defers ending the session after the transaction is committed or ended
	defer session.EndSession(ctx)

	// Inserts multiple documents into a collection within a transaction,
	// then commits or ends the transaction
	result, err := session.WithTransaction(context.TODO(), fn, opts...)

	return result, err
}

func (s *StoreImpl) Close() bool {
	if s.closed == true {
		return true
	}

	err := s.rds.Close()
	if err != nil {
		return false
	}

	var db *sql.DB
	db, err = s.mysql.DB()
	if err != nil {
		return false
	}
	err = db.Close()
	if err != nil {
		return false
	}

	if s.mgo != nil {
		if err = s.mgo.Client().Disconnect(context.Background()); err != nil {
			return false
		}
	}

	s.closed = true

	return true
}
