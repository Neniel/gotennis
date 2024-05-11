package mongodb

import (
	"context"
	"errors"

	"github.com/Neniel/gotennis/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDbReader struct {
	MongodbClient *mongo.Client
	Categories    *mongo.Collection
	Players       *mongo.Collection
}

func NewMongoDbReader(client *mongo.Client) *MongoDbReader {
	return &MongoDbReader{
		MongodbClient: client,
		Categories:    client.Database("tennis").Collection("categories"),
		Players:       client.Database("tennis").Collection("players"),
	}

}

func (mdbr *MongoDbReader) GetCategories(ctx context.Context) ([]entity.Category, error) {
	cursor, err := mdbr.Categories.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	var output []entity.Category
	cursor.All(ctx, &output)
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
	cursor.All(ctx, &output)
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

func (mdbr *MongoDbReader) ExistsByGovernmentID(ctx context.Context, governmentID string) (bool, error) {
	result := mdbr.Players.FindOne(ctx, bson.D{{Key: "government_id", Value: governmentID}})
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return false, nil
	}

	if result.Err() != nil {
		return false, result.Err()
	}

	return true, nil
}

type MongoDbWriter struct {
	mongodbClient *mongo.Client
	Categories    *mongo.Collection
	Players       *mongo.Collection
}

func NewMongoDbWriter(client *mongo.Client) *MongoDbWriter {
	return &MongoDbWriter{
		mongodbClient: client,
		Categories:    client.Database("tennis").Collection("categories"),
		Players:       client.Database("tennis").Collection("players"),
	}

}

func (mdbw *MongoDbWriter) AddCategory(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	//category.ID = primitive.NewObjectID()
	//category.CreatedAt = time.Now().UTC()
	_, err := mdbw.Categories.InsertOne(ctx, category)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (mdbw *MongoDbWriter) UpdateCategory(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	return nil, nil
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
	_, err := mdbw.Players.InsertOne(ctx, player)
	if err != nil {
		return nil, err
	}
	return player, nil
}

func (mdbw *MongoDbWriter) UpdatePlayer(ctx context.Context, category *entity.Player) (*entity.Player, error) {
	return nil, nil
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
