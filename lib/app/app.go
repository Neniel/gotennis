package app

import (
	"context"
	"fmt"
	"os"

	"github.com/Neniel/gotennis/lib/config"
	"github.com/Neniel/gotennis/lib/entity"
	"github.com/Neniel/gotennis/lib/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IApp interface {
	GetSystemMongoDBClient() *mongo.Client
	GetMongoDBClients() map[string]*CustomerMongoDB
}

type CustomerMongoDB struct {
	ID            string
	DatabaseName  string
	MongoDBClient *mongo.Client
}

type App struct {
	SystemMongoDBClient     *mongo.Client
	CustomersMongoDBClients map[string]*CustomerMongoDB
}

func (a *App) GetMongoDBClients() map[string]*CustomerMongoDB {
	return a.CustomersMongoDBClients
}

func (a *App) GetSystemMongoDBClient() *mongo.Client {
	return a.SystemMongoDBClient
}

func NewApp(ctx context.Context) IApp {
	// Store all mongodb clients for each of the customers
	mongoDBClients := make(map[string]*CustomerMongoDB)

	c, err := config.LoadConfiguration()
	if err != nil {
		log.Logger.Error(fmt.Errorf("error while loading configurations: %w", err).Error())
		os.Exit(1)
	}

	systemMongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(c.SystemDataSource.URI))
	if err != nil {
		log.Logger.Error(fmt.Errorf("error while connecting to system database: %w", err).Error())
		os.Exit(1)
	}

	if err := systemMongoClient.Ping(ctx, nil); err != nil {
		log.Logger.Error(fmt.Errorf("error while checking system database connection: %w", err).Error())
		os.Exit(1)
	}

	tier := os.Getenv("TIER")

	if tier == "diamond" {
		_id, err := primitive.ObjectIDFromHex(os.Getenv("CUSTOMER_ID"))
		if err != nil {
			log.Logger.Error(fmt.Errorf("error while reading CUSTOMER_ID environment variable: %w", err).Error())
			os.Exit(1)
		}

		var customer entity.Customer
		err = systemMongoClient.Database("system").Collection("customers").FindOne(ctx, bson.D{{Key: "_id", Value: _id}}).Decode(&customer)
		if err != nil {
			log.Logger.Error(fmt.Errorf("error while fetching database of customer '%s': %w", _id, err).Error())
			os.Exit(1)
		}

		customerMongoDBClient, err := mongo.Connect(ctx, options.Client().ApplyURI(customer.MongoDBConnectionString))
		if err != nil {
			log.Logger.Error(fmt.Errorf("error while connecting to '%s' customer database: %w", customer.Name, err).Error())
			os.Exit(1)
		}
		if err := customerMongoDBClient.Ping(ctx, nil); err != nil {
			log.Logger.Error(fmt.Errorf("error while checking customer database connection: %w", err).Error())
			os.Exit(1)
		}

		mongoDBClients[customer.ID.Hex()] = &CustomerMongoDB{
			ID:            customer.ID.Hex(),
			DatabaseName:  customer.DatabaseName,
			MongoDBClient: customerMongoDBClient,
		}

	} else {
		customersCursor, err := systemMongoClient.Database("system").Collection("customers").Find(ctx, bson.D{}) // TODO filter out diamond clients
		if err != nil {
			log.Logger.Error(fmt.Errorf("error while fetching to customers from database: %w", err).Error())
			os.Exit(1)
		}

		customers := make([]entity.Customer, 0)
		err = customersCursor.All(ctx, &customers)
		if err != nil {
			log.Logger.Error(fmt.Errorf("error while parsing customers: %w", err).Error())
			os.Exit(1)
		}

		mongoDBClients = loadCustomersClients(ctx, customers)
	}

	/*
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
	*/

	/*
		redisClient := redis.NewClient(&redis.Options{
			Addr:     c.Redis.Address,
			Password: c.Redis.Password,
		})
		log.Logger.Info("Connected to Redis")
	*/
	return &App{
		SystemMongoDBClient:     systemMongoClient,
		CustomersMongoDBClients: mongoDBClients,
	}

}

func loadCustomersClients(ctx context.Context, customers []entity.Customer) map[string]*CustomerMongoDB {
	mongoDBClients := make(map[string]*CustomerMongoDB)
	for _, customer := range customers {
		customerMongoDBClient, err := mongo.Connect(ctx, options.Client().ApplyURI(customer.MongoDBConnectionString))
		if err != nil {
			log.Logger.Warn(fmt.Errorf("error while connecting to '%s' customer database: %w", customer.Name, err).Error())
			continue
		}

		if err := customerMongoDBClient.Ping(ctx, nil); err != nil {
			log.Logger.Warn(fmt.Errorf("error while checking '%s' customer database connection: %w", customer.Name, err).Error())
			continue
		}

		if _, ok := mongoDBClients[customer.ID.Hex()]; !ok {
			mongoDBClients[customer.ID.Hex()] = &CustomerMongoDB{
				ID:            customer.ID.Hex(),
				DatabaseName:  customer.DatabaseName,
				MongoDBClient: customerMongoDBClient,
			}
		}
	}

	return mongoDBClients
}
