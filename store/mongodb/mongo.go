package mongodb

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"go.uber.org/zap"

	"github.com/BooeZhang/gin-layout/config"
	"github.com/BooeZhang/gin-layout/pkg/log"
)

var (
	dbClient *mongo.Client
	db       *mongo.Database
)

func InitMongo(cf *config.Mongo) {
	DialToMongo(cf)
}

// DialToMongo 根据配置连接到mongo
func DialToMongo(op *config.Mongo) {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOpt := options.Client()
	dbClient, err = mongo.Connect(ctx, clientOpt.ApplyURI(op.Uri).SetMaxPoolSize(uint64(op.PoolLimit)))
	if err != nil {
		panic(err)
	}
	err = dbClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Error("mongodb connect failed, err: ", zap.Error(err))
		os.Exit(1)
	}
	connStr, _ := connstring.Parse(op.Uri)
	db = dbClient.Database(connStr.Database)
}

// EnsureIndex 添加索引
func EnsureIndex(ms *mongo.Client, collection string, ensureIndex []string, unique bool) {
	var (
		indexs []mongo.IndexModel
	)
	for _, name := range ensureIndex {
		indexs = append(indexs, mongo.IndexModel{
			Keys:    bson.D{{name, 1}},
			Options: options.Index().SetUnique(unique),
		})
	}
	opts := options.CreateIndexes().SetMaxTime(2 * time.Second)
	names, err := ms.Database("").Collection(collection).Indexes().CreateMany(context.TODO(), indexs, opts)
	if err != nil {
		panic(err)
	}
	fmt.Println(names)
}

// GetClient 获取mongo连接会话
func GetClient() *mongo.Client {
	return dbClient
}

// GetDBSession 获取指定数据库
func GetDBSession() *mongo.Database {
	return db
}
