package app

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	DBClients *DBClients
}

type DBClients struct {
	MongoDB *mongo.Client
	Redis   *redis.Client
}

func NewApp(ctx context.Context) *App {
	bsMongoURI, err := os.ReadFile(os.Getenv("MONGODB_URI_FILE"))
	if err != nil {
		log.Fatalln(err.Error())
	}
	mongoURI := strings.Replace(string(bsMongoURI), "\n", "", -1)

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("Connected to MongoDB")

	bsRedisAddress, err := os.ReadFile(os.Getenv("REDIS_ADDRESS_FILE"))
	if err != nil {
		log.Fatalln(err.Error())
	}
	redisAddress := strings.Replace(string(bsRedisAddress), "\n", "", -1)

	bsRedisPassword, err := os.ReadFile(os.Getenv("REDIS_PASSWORD_FILE"))
	if err != nil {
		log.Fatalln(err.Error())
	}
	redisPassword := strings.Replace(string(bsRedisPassword), "\n", "", -1)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: redisPassword,
	})
	log.Println("Connected to Redis")

	return &App{
		DBClients: &DBClients{
			MongoDB: mongoClient,
			Redis:   redisClient,
		},
	}

}
