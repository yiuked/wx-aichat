package config

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"log"
	"os"
	"time"
	"wx-aichat/lib"
)

// 全局配置
var (
	WxToken    string
	WxAppid    string
	WxSecret   string
	WxMediaId  string
	WXOpenId   string
	WxAuthor   string
	OpenAIKey  string
	HOST       string
	DataSource string
	Mgo        *MgoClient
	Limit      map[string]*UserLimit
	Cache      *lib.Cache
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
	WxAppid = os.Getenv("WX_APPID")
	WxSecret = os.Getenv("WX_SECRET")
	WxMediaId = os.Getenv("WX_MEDIA_ID")
	WXOpenId = os.Getenv("WX_OPEN_ID")
	WxAuthor = os.Getenv("WX_AUTHOR")
	DataSource = os.Getenv("DATA_SOURCE")
	OpenAIKey = os.Getenv("OPEN_AI_KEY")
	HOST = os.Getenv("HOST")
	err := initMongoDB()
	if err != nil {
		log.Println(DataSource)
		panic(err)
	}
	Limit = make(map[string]*UserLimit)
	// 初始化缓存
	Cache = lib.NewCache(10 * time.Second)
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
