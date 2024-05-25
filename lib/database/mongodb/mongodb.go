package mongodb

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/Neniel/gotennis/lib/entity"
	"github.com/Neniel/gotennis/lib/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDbReader struct {
	MongodbClient *mongo.Client
	Categories    *mongo.Collection
	Players       *mongo.Collection
	Tournaments   *mongo.Collection
}

func NewMongoDbReader(client *mongo.Client) *MongoDbReader {
	return &MongoDbReader{
		MongodbClient: client,
		Categories:    client.Database("tennis").Collection("categories"),
		Players:       client.Database("tennis").Collection("players"),
		Tournaments:   client.Database("tennis").Collection("tournaments"),
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

type MongoDbWriter struct {
	mongodbClient *mongo.Client
	Categories    *mongo.Collection
	Players       *mongo.Collection
	Tournaments   *mongo.Collection
}

func NewMongoDbWriter(client *mongo.Client) *MongoDbWriter {
	return &MongoDbWriter{
		mongodbClient: client,
		Categories:    client.Database("tennis").Collection("categories"),
		Players:       client.Database("tennis").Collection("players"),
		Tournaments:   client.Database("tennis").Collection("tournaments"),
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
