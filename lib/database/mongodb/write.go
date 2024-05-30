package mongodb

import (
	"context"
	"fmt"
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
	DB            *mongo.Database
}

func NewMongoDbWriter(client *mongo.Client, databaseName string) *MongoDbWriter {
	return &MongoDbWriter{
		mongodbClient: client,
		DB:            client.Database(databaseName),
	}

}

func (mdbw *MongoDbWriter) AddCategory(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	category.ID = primitive.NewObjectID()
	category.CreatedAt = time.Now().UTC()

	_, err := mdbw.DB.Collection("categories").InsertOne(ctx, category)
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

	_, err = mdbw.DB.Collection("categories").ReplaceOne(ctx, bson.M{"_id": category.ID}, updatedCatgory)
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

	_, err = mdbw.DB.Collection("categories").DeleteOne(ctx, bson.D{{Key: "_id", Value: _id}})
	if err != nil {
		return err
	}

	return nil
}

func (mdbw *MongoDbWriter) AddPlayer(ctx context.Context, player *entity.Player) (*entity.Player, error) {
	player.ID = primitive.NewObjectID()
	player.CreatedAt = time.Now().UTC()
	player.TemporaryAccessCode = fmt.Sprintf("%v", rand.Uint32())

	session, err := mdbw.DB.Client().StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}

		if _, err := mdbw.DB.Collection("players").InsertOne(sc, player); err != nil {
			if err := session.AbortTransaction(sc); err != nil {
				return err
			}
			return err
		}

		if _, err := mdbw.DB.Collection("users").InsertOne(sc, &entity.User{
			ID:                  player.ID,
			GovernmentID:        player.GovernmentID,
			Email:               player.Email,
			Alias:               player.Alias,
			TemporaryAccessCode: player.TemporaryAccessCode,
			CreatedAt:           player.CreatedAt,
		}); err != nil {
			if err := session.AbortTransaction(sc); err != nil {
				return err
			}
			return err
		}

		if err := session.CommitTransaction(sc); err != nil {
			return err
		}

		return nil

	})

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

	_, err = mdbw.DB.Collection("players").ReplaceOne(ctx, bson.M{"_id": player.ID}, updatedPlayer)
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

	_, err = mdbw.DB.Collection("players").DeleteOne(ctx, bson.D{{Key: "_id", Value: _id}})
	if err != nil {
		return err
	}

	return nil
}

func (mdbw *MongoDbWriter) AddTournament(ctx context.Context, tournament *entity.Tournament) (*entity.Tournament, error) {
	tournament.ID = primitive.NewObjectID()
	_, err := mdbw.DB.Collection("tournaments").InsertOne(ctx, tournament)
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

	_, err = mdbw.DB.Collection("tournaments").ReplaceOne(ctx, bson.M{"_id": tournament.ID}, updatedTournament)
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

	_, err = mdbw.DB.Collection("tournaments").DeleteOne(ctx, bson.D{{Key: "_id", Value: _id}})
	if err != nil {
		return err
	}

	return nil
}

func (mdbw *MongoDbWriter) AddTenant(ctx context.Context, tenant *entity.Tenant) (*entity.Tenant, error) {
	tenant.ID = primitive.NewObjectID()
	tenant.CreatedAt = time.Now().UTC()

	_, err := mdbw.DB.Collection("tenants").InsertOne(ctx, tenant)
	if err != nil {
		return nil, err
	}

	return tenant, nil
}

func (mdbw *MongoDbWriter) DeleteTenant(ctx context.Context, id string) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = mdbw.DB.Collection("tenants").DeleteOne(ctx, bson.D{{Key: "_id", Value: _id}})
	if err != nil {
		return err
	}

	return nil
}
