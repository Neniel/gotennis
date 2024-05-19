package app

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Neniel/gotennis/logger"
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
	runWith := os.Getenv("RUN_WITH")
	switch runWith {
	case "k8s":
		return NewKubernetesApp(ctx)
	case "localhost", "docker":
		return NewCommonApp(ctx)
	default:
		return nil
	}
}

func NewCommonApp(ctx context.Context) IApp {
	bsMongoURI, err := os.ReadFile(os.Getenv("MONGODB_URI_FILE"))
	if err != nil {
		logger.Fatal(err.Error())
	}
	mongoURI := strings.Replace(string(bsMongoURI), "\n", "", -1)

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		logger.Fatal(err.Error())
	}
	if err := mongoClient.Ping(ctx, nil); err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info("Connected to MongoDB")

	db := mongoClient.Database(util.DBName)

	logger.Info("Preparing 'players' collection")
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
		logger.Fatal(err.Error())
	}

	logger.Info(fmt.Sprintf("Indexes %v created!", indexNames))

	bsRedisAddress, err := os.ReadFile(os.Getenv("REDIS_ADDRESS_FILE"))
	if err != nil {
		logger.Fatal(err.Error())
	}
	redisAddress := strings.Replace(string(bsRedisAddress), "\n", "", -1)

	bsRedisPassword, err := os.ReadFile(os.Getenv("REDIS_PASSWORD_FILE"))
	if err != nil {
		logger.Fatal(err.Error())
	}
	redisPassword := strings.Replace(string(bsRedisPassword), "\n", "", -1)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: redisPassword,
	})
	logger.Info("Connected to Redis")

	return &App{
		DBClients: &DBClients{
			MongoDB: mongoClient,
			Redis:   redisClient,
		},
	}

}

func NewKubernetesApp(ctx context.Context) IApp {
	mongoURI := os.Getenv("MONGODB_CONNECTION_STRING")

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		logger.Fatal(err.Error())
	}
	if err := mongoClient.Ping(ctx, nil); err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info("Connected to MongoDB")

	db := mongoClient.Database(util.DBName)

	logger.Info("Preparing 'players' collection")
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
		logger.Fatal(err.Error())
	}

	logger.Info(fmt.Sprintf("Indexes %v created!", indexNames))

	redisAddress := os.Getenv("REDIS_SERVER_ADDRESS")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: redisPassword,
	})
	logger.Info("Connected to Redis")

	return &App{
		DBClients: &DBClients{
			MongoDB: mongoClient,
			Redis:   redisClient,
		},
	}

}
