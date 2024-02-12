package database

import (
	"context"

	"categories/database/mongodb"
	redisdb "categories/database/redis"

	"github.com/Neniel/gotennis/entity"

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
}

type DBWriter interface {
	AddCategory(context.Context, *entity.Category) (*entity.Category, error)
	UpdateCategory(context.Context, *entity.Category) (*entity.Category, error)
	DeleteCategory(context.Context, string) error
}

func NewDatabaseReader(client interface{}) DBReader {
	mongoClient, isMongoClient := client.(*mongo.Client)
	if isMongoClient {
		return mongodb.NewMongoDbReader(mongoClient)
	}

	redisClient, isRedisClient := client.(*redis.Client)
	if isRedisClient {
		return redisdb.NewRedisReader(redisClient)
	}

	panic("DB_TYPE is not supported")
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

	panic("DB_TYPE is not supported")
}
