package config

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"os"
	"time"
)

// 全局配置
var (
	WxToken    string
	AIToken    string
	DataSource string
	Mgo        *MgoClient
	Limit      map[string]*UserLimit
)

type UserLimit struct {
	LastSt int64
	Cnt    int64
}

type MgoClient struct {
	Client *mongo.Client
	Db     *mongo.Database
}

func init() {
	WxToken = os.Getenv("WX_TOKEN")
	DataSource = os.Getenv("DATA_SOURCE")
	AIToken = os.Getenv("AI_TOKEN")
	err := initMongoDB()
	if err != nil {
		panic(err)
	}
	Limit = make(map[string]*UserLimit)
}

// initMongoDB 初始化mongodb
func initMongoDB() error {
	// 验证数据源
	cs, err := connstring.ParseAndValidate(DataSource)
	if err != nil {
		return err
	}

	Mgo = &MgoClient{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Client().ApplyURI(DataSource)
	cli, err := mongo.Connect(ctx, opts)
	if err != nil {
		return err
	}
	Mgo.Client = cli

	err = Mgo.Client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}

	Mgo.Db = Mgo.Client.Database(cs.Database)
	return nil
}
