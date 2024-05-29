package mongodb

import (
	"context"
	"errors"

	"github.com/Neniel/gotennis/lib/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDbReader struct {
	MongodbClient  *mongo.Client
	MongoDBClients map[string]*mongo.Client
	Categories     *mongo.Collection
	Players        *mongo.Collection
	Tournaments    *mongo.Collection
}

func NewMongoDbReader(client *mongo.Client, databaseName string) *MongoDbReader {
	return &MongoDbReader{
		MongodbClient: client,
		Categories:    client.Database(databaseName).Collection("categories"),
		Players:       client.Database(databaseName).Collection("players"),
		Tournaments:   client.Database(databaseName).Collection("tournaments"),
	}

}

func (mdbr *MongoDbReader) GetCategories(ctx context.Context) ([]entity.Category, error) {
	cursor, err := mdbr.Categories.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	var output []entity.Category
	if err := cursor.All(ctx, &output); err != nil {
		return nil, err
	}

	return output, nil
}

func (mdbr *MongoDbReader) GetCategory(ctx context.Context, id string) (*entity.Category, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var result entity.Category
	err = mdbr.Categories.FindOne(context.Background(), bson.D{{Key: "_id", Value: _id}}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (mdbr *MongoDbReader) GetPlayers(ctx context.Context) ([]entity.Player, error) {
	cursor, err := mdbr.Players.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	var output []entity.Player
	if err := cursor.All(ctx, &output); err != nil {
		return nil, err
	}

	return output, nil
}

func (mdbr *MongoDbReader) GetPlayer(ctx context.Context, id string) (*entity.Player, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var result entity.Player
	err = mdbr.Players.FindOne(ctx, bson.D{{Key: "_id", Value: _id}}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (mdbr *MongoDbReader) IsAvailable(ctx context.Context, field string, value string) (bool, error) {
	result := mdbr.Players.FindOne(context.TODO(), bson.D{{Key: field, Value: value}})
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return true, nil
	}

	if result.Err() != nil {
		return false, result.Err()
	}

	return false, nil
}

func (mdbr *MongoDbReader) GetTournaments(ctx context.Context) ([]entity.Tournament, error) {
	cursor, err := mdbr.Tournaments.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	var output []entity.Tournament
	if err := cursor.All(ctx, &output); err != nil {
		return nil, err
	}

	return output, nil
}

func (mdbr *MongoDbReader) GetTournament(ctx context.Context, id string) (*entity.Tournament, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var result entity.Tournament
	err = mdbr.Tournaments.FindOne(ctx, bson.D{{Key: "_id", Value: _id}}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
