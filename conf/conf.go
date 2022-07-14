package conf

import (
	"chat/model"
	"context"
	"fmt"
	logging "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/ini.v1"
	"strings"
)

var (
	MongoDBClient *mongo.Client
	AppMode    string
	HttpPort    string
	Db    string
	DbHost    string
	DbPort    string
	DbUser    string
	DbPassWord    string
	DbName    string
	MongoDBName    string
	MongoDBAddr    string
	MongoDBPwd    string
	MongoDBPort    string
)

func Init()  {
	file, err := ini.Load("./conf/conf.ini")
	if err != nil {
		logging.Fatalf("ini load file failed, err:%s \n", err)
	}

	LoadServer(file)
	LoadMysql(file)
	LoadMongoDB(file)

	// mysql path
	path := strings.Join([]string{DbUser, ":", DbPassWord, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8&parseTime=true&interpolateParams=True"}, "")
	model.ConnectMysql(path)

	ConnectMongoDB()

}

func ConnectMongoDB()  {
	// 设置 mongodb 客户端连接信息
	// 无用户密码的
	//path := fmt.Sprintf("mongodb://%s:%s", MongoDBAddr, MongoDBPort)
	// 由用户密码的
	path := fmt.Sprintf("mongodb://root:123456@%s:%s", MongoDBAddr, MongoDBPort)
	fmt.Println(path)
	fmt.Println(MongoDBPort)
	clientOptions := options.Client().ApplyURI(path)

	var err error
	MongoDBClient, err = mongo.Connect(context.TODO(), clientOptions)
	// 这样赋值有问题，没明白？
	//MongoDBClient, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		logging.Fatalf("mongo connect failed, err:%s \n", err)
	}

	err = MongoDBClient.Ping(context.TODO(), nil)

	if err != nil {
		logging.Fatalf("mongo ping failed, err:%s \n", err)
	}


	logging.Println("mongo connect success")
	//MongoDBClient = MongoDBClient
	//fmt.Println(MongoDBClient)
}

func LoadServer(file *ini.File)  {
	AppMode = file.Section("service").Key("AppMode").String()
	HttpPort = file.Section("service").Key("HttpPort").String()
}

func LoadMysql(file *ini.File)  {
	Db = file.Section("mysql").Key("Db").String()
	DbHost = file.Section("mysql").Key("DbHost").String()
	DbPort = file.Section("mysql").Key("DbPort").String()
	DbUser = file.Section("mysql").Key("DbUser").String()
	DbPassWord = file.Section("mysql").Key("DbPassWord").String()
	DbName = file.Section("mysql").Key("DbName").String()
}


func LoadMongoDB(file *ini.File)  {
	MongoDBName = file.Section("MongoDB").Key("MongoDBName").String()
	MongoDBAddr = file.Section("MongoDB").Key("MongoDBAddr").String()
	MongoDBPwd = file.Section("MongoDB").Key("MongoDBPwd").String()
	MongoDBPort = file.Section("MongoDB").Key("MongoDBPort").String()
}