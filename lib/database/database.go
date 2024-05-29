package database

import (
	"context"

	"github.com/Neniel/gotennis/lib/database/mongodb"
	"github.com/Neniel/gotennis/lib/entity"

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

	GetTournaments(context.Context) ([]entity.Tournament, error)
	GetTournament(context.Context, string) (*entity.Tournament, error)
}

type DBWriter interface {
	AddCategory(context.Context, *entity.Category) (*entity.Category, error)
	UpdateCategory(context.Context, *entity.Category) (*entity.Category, error)
	DeleteCategory(context.Context, string) error

	AddPlayer(context.Context, *entity.Player) (*entity.Player, error)
	UpdatePlayer(context.Context, *entity.Player) (*entity.Player, error)
	DeletePlayer(context.Context, string) error

	AddTournament(context.Context, *entity.Tournament) (*entity.Tournament, error)
	UpdateTournament(context.Context, *entity.Tournament) (*entity.Tournament, error)
	DeleteTournament(context.Context, string) error
}

func NewDatabaseReader(client interface{}, databaseName string) DBReader {
	mongoClient, isMongoClient := client.(*mongo.Client)
	if isMongoClient {
		return mongodb.NewMongoDbReader(mongoClient, databaseName)
	}

	/*redisClient, isRedisClient := client.(*redis.Client)
	if isRedisClient {
		return redisdb.NewRedisReader(redisClient)
	}*/

	panic("client is not supported")
}

func NewDatabaseWriter(client interface{}, databaseName string) DBWriter {
	mongoClient, isMongoClient := client.(*mongo.Client)
	if isMongoClient {
		return mongodb.NewMongoDbWriter(mongoClient, databaseName)
	}

	/*
		redisClient, isRedisClient := client.(*redis.Client)
		if isRedisClient {
			return redisdb.NewRedisWriter(redisClient)
		}
	*/

	panic("client is not supported")
}
