package mongodb

import (
	"context"

	"github.com/Neniel/tennis/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDbReader struct {
	MongodbClient *mongo.Client
	Categories    *mongo.Collection
}

func NewMongoDbReader(client *mongo.Client) *MongoDbReader {
	return &MongoDbReader{
		MongodbClient: client,
		Categories:    client.Database("tennis").Collection("categories"),
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
		return nil, primitive.ErrInvalidHex
	}

	var result entity.Category
	err = mdbr.Categories.FindOne(context.Background(), bson.D{{Key: "_id", Value: _id}}).Decode(&result)
	if err != nil {
		return nil, mongo.ErrNoDocuments
	}
	return &result, nil
}

type MongoDbWriter struct {
	mongodbClient *mongo.Client
	Categories    *mongo.Collection
}

func NewMongoDbWriter(client *mongo.Client) *MongoDbWriter {
	return &MongoDbWriter{
		mongodbClient: client,
		Categories:    client.Database("tennis").Collection("categories"),
	}

}

func (mdbw *MongoDbWriter) AddCategory(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	category.ID = primitive.NewObjectID()
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
		return primitive.ErrInvalidHex
	}
	_, err = mdbw.Categories.DeleteOne(ctx, bson.D{{Key: "_id", Value: _id}})
	if err != nil {
		return err
	}

	return nil
}
