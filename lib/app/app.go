package app

import (
	"context"
	"fmt"
	"os"

	"github.com/Neniel/gotennis/config"
	"github.com/Neniel/gotennis/log"
	"github.com/Neniel/gotennis/util"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IApp interface {
	GetMongoDBClient() *mongo.Client
	GetRedisClient() *redis.Client
}

type App struct {
	DBClients *DBClients
}

func (a *App) GetMongoDBClient() *mongo.Client {
	return a.DBClients.MongoDB
}

func (a *App) GetRedisClient() *redis.Client {
	return a.DBClients.Redis
}

type DBClients struct {
	MongoDB *mongo.Client
	Redis   *redis.Client
}

func NewApp(ctx context.Context) IApp {
	c, err := config.LoadConfiguration()
	if err != nil {
		log.Logger.Error(err.Error())
		os.Exit(1)
	}

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(c.MongoDB.URI))
	if err != nil {
		log.Logger.Error(err.Error())
		os.Exit(1)
	}
	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Logger.Error(err.Error())
		os.Exit(1)
	}
	log.Logger.Info("Connected to MongoDB")

	db := mongoClient.Database(util.DBName)

	log.Logger.Info("Preparing 'players' collection")
	playersColl := db.Collection(util.CollNamePlayers)
	indexOptionsGovernmentID := options.Index().SetUnique(true).SetPartialFilterExpression(bson.M{"government_id": bson.M{"$type": "string"}})
	indexOptionsEmail := options.Index().SetUnique(true).SetPartialFilterExpression(bson.M{"email": bson.M{"$type": "string"}})
	indexOptionsAlias := options.Index().SetUnique(true).SetPartialFilterExpression(bson.M{"alias": bson.M{"$type": "string"}})

	indexModelGovernmentID := mongo.IndexModel{
		Keys:    bson.D{{"government_id", 1}},
		Options: indexOptionsGovernmentID,
	}

	indexModelEmail := mongo.IndexModel{
		Keys:    bson.D{{"email", 1}},
		Options: indexOptionsEmail,
	}

	indexModelAlias := mongo.IndexModel{
		Keys:    bson.D{{"alias", 1}},
		Options: indexOptionsAlias,
	}

	indexNames, err := playersColl.Indexes().CreateMany(ctx, []mongo.IndexModel{indexModelGovernmentID, indexModelEmail, indexModelAlias})
	if err != nil {
		log.Logger.Error(err.Error())
		os.Exit(1)
	}

	log.Logger.Info(fmt.Sprintf("Indexes %v created!", indexNames))

	redisClient := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Address,
		Password: c.Redis.Password,
	})
	log.Logger.Info("Connected to Redis")

	return &App{
		DBClients: &DBClients{
			MongoDB: mongoClient,
			Redis:   redisClient,
		},
	}

}
