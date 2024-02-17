package store

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type Storage interface {
	GetMySQL() *gorm.DB
	GetMongo() *mongo.Database
	GetRedis() redis.UniversalClient
	Close() bool
	WithMongoTransaction(ctx context.Context, fn func(ctx mongo.SessionContext) (interface{}, error), opts ...*options.TransactionOptions) (interface{}, error)
}
