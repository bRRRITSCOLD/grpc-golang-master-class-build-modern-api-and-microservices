package mongodb_client

import (
	"context"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/* Used to create a singleton object of MongoDB client.
Initialized and exposed through  GetMongoClient().*/
var clientInstance *mongo.Client

//Used during creation of singleton client object in GetMongoClient().
var clientInstanceError error

//Used to execute client creation procedure only once.
var mongoOnce sync.Once

//I have used below constants just to hold required database config's.
const (
	MONGO_DB_USERNAME = "MONGO_DB_USERNAME"
	MONGO_DB_PASSWORD = "MONGO_DB_PASSWORD"
	MONGO_DB_HOST     = "MONGO_DB_HOST"
	MONGO_DB_PORT     = "MONGO_DB_PORT"
	MONGO_DB_DATABASE = "MONGO_DB_DATABASE"
)

var (
	client *mongo.Client
)

type MongoDBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func getMongoDBConfig() *MongoDBConfig {
	return &MongoDBConfig{
		Host:     os.Getenv(MONGO_DB_HOST),
		Port:     os.Getenv(MONGO_DB_PORT),
		Database: os.Getenv(MONGO_DB_DATABASE),
		Username: os.Getenv(MONGO_DB_USERNAME),
		Password: os.Getenv(MONGO_DB_PASSWORD),
	}
}

func Init() error {
	config := getMongoDBConfig()

	// Set client options
	cOpts := options.Client().
		ApplyURI("mongodb://" + config.Host + ":" + config.Port).
		SetAuth(options.Credential{
			Username: config.Username,
			Password: config.Password,
		})

	// Connect to MongoDB
	c, err := mongo.Connect(context.TODO(), cOpts)
	if err != nil {
		clientInstanceError = err
	}

	// Check the connection
	err = c.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}

	client = c

	return nil
}

//GetMongoClient - Return mongodb connection to work with
func GetClient() *mongo.Client {
	return client
}

func GetDatabase(database string) *mongo.Database {
	return GetClient().Database(database)
}

func GetCollection(database string, collection string) *mongo.Collection {
	return GetDatabase(database).Collection(collection)
}

// const (
// 	USERS_MYSQL_DB_USERNAME = "USERS_MYSQL_DB_USERNAME"
// 	USERS_MYSQL_DB_PASSWORD = "USERS_MYSQL_DB_PASSWORD"
// 	USERS_MYSQL_DB_HOST     = "USERS_MYSQL_DB_HOST"
// 	USERS_MYSQL_DB_DATABASE = "USERS_MYSQL_DB_DATABASE"
// )

// var (
// 	Client *sqlx.DB
// )

// func Init() {
// 	username := os.Getenv(USERS_MYSQL_DB_USERNAME)
// 	password := os.Getenv(USERS_MYSQL_DB_PASSWORD)
// 	host := os.Getenv(USERS_MYSQL_DB_HOST)
// 	database := os.Getenv(USERS_MYSQL_DB_DATABASE)

// 	datasourceName := fmt.Sprintf(
// 		"%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true",
// 		username,
// 		password,
// 		host,
// 		database,
// 	)

// 	var openErr error
// 	Client, openErr = sqlx.Open("mysql", datasourceName)
// 	errors_utils.PanicOnError(openErr)

// 	pingErr := Client.Ping()
// 	errors_utils.PanicOnError(pingErr)

// 	Client.Mapper = reflectx.NewMapperFunc("mysql", strings.ToLower)

// 	// mysql.SetLogger()
// 	log.Println("database successfully configured")
// }
