package database

import (
	"context"

	"github.com/Neniel/gotennis/lib/database/mongodb"
	redisdb "github.com/Neniel/gotennis/lib/database/redis"
	"github.com/Neniel/gotennis/lib/entity"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

type Database interface {
	DBReader
	DBWriter
}

type DBReader interface {
	GetCategories(context.Context) ([]entity.Category, error)
	GetCategory(context.Context, string) (*entity.Category, error)

	GetPlayers(context.Context) ([]entity.Player, error)
	GetPlayer(context.Context, string) (*entity.Player, error)
	IsAvailable(context.Context, string, string) (bool, error)
}

type DBWriter interface {
	AddCategory(context.Context, *entity.Category) (*entity.Category, error)
	UpdateCategory(context.Context, *entity.Category) (*entity.Category, error)
	DeleteCategory(context.Context, string) error

	AddPlayer(context.Context, *entity.Player) (*entity.Player, error)
	UpdatePlayer(context.Context, *entity.Player) (*entity.Player, error)
	DeletePlayer(context.Context, string) error
}

func NewDatabaseReader(client interface{}) DBReader {
	mongoClient, isMongoClient := client.(*mongo.Client)
	if isMongoClient {
		return mongodb.NewMongoDbReader(mongoClient)
	}

	/*redisClient, isRedisClient := client.(*redis.Client)
	if isRedisClient {
		return redisdb.NewRedisReader(redisClient)
	}*/

	panic("client is not supported")
}

func NewDatabaseWriter(client interface{}) DBWriter {
	mongoClient, isMongoClient := client.(*mongo.Client)
	if isMongoClient {
		return mongodb.NewMongoDbWriter(mongoClient)
	}

	redisClient, isRedisClient := client.(*redis.Client)
	if isRedisClient {
		return redisdb.NewRedisWriter(redisClient)
	}

	panic("client is not supported")
}
