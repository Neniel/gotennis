package mongodb

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/Neniel/gotennis/lib/entity"
	"github.com/Neniel/gotennis/lib/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDbWriter struct {
	mongodbClient *mongo.Client
	Categories    *mongo.Collection
	Players       *mongo.Collection
	Tournaments   *mongo.Collection
}

func NewMongoDbWriter(client *mongo.Client, databaseName string) *MongoDbWriter {
	return &MongoDbWriter{
		mongodbClient: client,
		Categories:    client.Database(databaseName).Collection("categories"),
		Players:       client.Database(databaseName).Collection("players"),
		Tournaments:   client.Database(databaseName).Collection("tournaments"),
	}

}

func (mdbw *MongoDbWriter) AddCategory(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	category.ID = primitive.NewObjectID()
	category.CreatedAt = time.Now().UTC()

	_, err := mdbw.Categories.InsertOne(ctx, category)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (mdbw *MongoDbWriter) UpdateCategory(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	category.UpdatedAt = util.ToPtr(time.Now().UTC())

	updatedCatgory, err := bson.Marshal(&category)
	if err != nil {
		return nil, err
	}

	_, err = mdbw.Categories.ReplaceOne(ctx, bson.M{"_id": category.ID}, updatedCatgory)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (mdbw *MongoDbWriter) DeleteCategory(ctx context.Context, id string) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = mdbw.Categories.DeleteOne(ctx, bson.D{{Key: "_id", Value: _id}})
	if err != nil {
		return err
	}

	return nil
}

func (mdbw *MongoDbWriter) AddPlayer(ctx context.Context, player *entity.Player) (*entity.Player, error) {
	player.ID = primitive.NewObjectID()
	player.CreatedAt = time.Now().UTC()
	player.TemporaryAccessCode = util.ToPtr(rand.Uint32())

	_, err := mdbw.Players.InsertOne(ctx, player)
	if err != nil {
		return nil, err
	}
	return player, nil
}

func (mdbw *MongoDbWriter) UpdatePlayer(ctx context.Context, player *entity.Player) (*entity.Player, error) {
	player.UpdatedAt = util.ToPtr(time.Now().UTC())

	updatedPlayer, err := bson.Marshal(&player)
	if err != nil {
		return nil, err
	}

	_, err = mdbw.Players.ReplaceOne(ctx, bson.M{"_id": player.ID}, updatedPlayer)
	if err != nil {
		if e, ok := err.(mongo.WriteException); ok {
			for _, ee := range e.WriteErrors {
				if strings.Contains(ee.Message, "government_id_1") {
					return nil, &util.AppError{Message: "government_id has already been assigned to another player"}
				}

				if strings.Contains(ee.Message, "email_1") {
					return nil, &util.AppError{Message: "email has already been assigned to another player"}
				}

				if strings.Contains(ee.Message, "alias_1") {
					return nil, &util.AppError{Message: "alias has already been assigned to another player"}
				}
			}
		}

		return nil, err
	}

	return player, nil
}

func (mdbw *MongoDbWriter) DeletePlayer(ctx context.Context, id string) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = mdbw.Players.DeleteOne(ctx, bson.D{{Key: "_id", Value: _id}})
	if err != nil {
		return err
	}

	return nil
}

func (mdbw *MongoDbWriter) AddTournament(ctx context.Context, tournament *entity.Tournament) (*entity.Tournament, error) {
	tournament.ID = primitive.NewObjectID()
	_, err := mdbw.Tournaments.InsertOne(ctx, tournament)
	if err != nil {
		return nil, err
	}

	return tournament, nil
}

func (mdbw *MongoDbWriter) UpdateTournament(ctx context.Context, tournament *entity.Tournament) (*entity.Tournament, error) {
	updatedTournament, err := bson.Marshal(&tournament)
	if err != nil {
		return nil, err
	}

	_, err = mdbw.Tournaments.ReplaceOne(ctx, bson.M{"_id": tournament.ID}, updatedTournament)
	if err != nil {
		return nil, err
	}

	return tournament, nil
}

func (mdbw *MongoDbWriter) DeleteTournament(ctx context.Context, id string) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = mdbw.Tournaments.DeleteOne(ctx, bson.D{{Key: "_id", Value: _id}})
	if err != nil {
		return err
	}

	return nil
}

func (mdbw *MongoDbWriter) AddCustomer(ctx context.Context, customer *entity.Customer) (*entity.Customer, error) {
	customer.ID = primitive.NewObjectID()
	_, err := mdbw.mongodbClient.Database("system").Collection("customers").InsertOne(ctx, customer)
	if err != nil {
		return nil, err
	}

	return customer, nil
}
