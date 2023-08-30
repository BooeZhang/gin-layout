package mongo

import (
	"context"
	"fmt"
	"github.com/BooeZhang/gin-layout/config"
	"github.com/BooeZhang/gin-layout/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"os"
	"time"
)

var (
	db *mongo.Client
)

func InitMongo(cf *config.MongoConf) {
	DialToMongo(cf)
}

// DialToMongo 根据配置连接到mongo
func DialToMongo(op *config.MongoConf) {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOpt := options.Client()
	if len(op.Username) != 0 && len(op.Password) != 0 {
		clientOpt.SetAuth(options.Credential{
			Username: op.Username,
			Password: op.Password,
		})
	}
	db, err = mongo.Connect(ctx, clientOpt.ApplyURI(fmt.Sprintf("mongodb://%s", op.Host)).SetMaxPoolSize(uint64(op.PoolLimit)))
	if err != nil {
		panic(err)
	}
	err = db.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Error("mongo connect failed, err: ", zap.Error(err))
		os.Exit(1)
	}
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

// GetSession 获取mongo连接会话
func GetSession() *mongo.Client {
	return db
}

// GetDBSession 获取指定数据库
func GetDBSession(dbName string) *mongo.Database {
	if len(dbName) == 0 && dbName == "" {
		return db.Database(config.GetConfig().MongoConfig.Database)
	}
	return db.Database(dbName)
}
