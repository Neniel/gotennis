package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Neniel/gotennis/entity"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type cacheSyncronizer struct {
	mongoClient *mongo.Client
	redisClient *redis.Client
}

const databaseName string = "tennis"

var collections []string = []string{"categories"}

func (c *cacheSyncronizer) StartSync() {
	for _, collection := range collections {
		changeStream, err := c.getChangeStream(databaseName, collection)
		if err != nil {
			panic(err)
		}

		go func(cs *mongo.ChangeStream, collectionName string) {
			// Manejar señales de interrupción para salir graciosamente
			signalChan := make(chan os.Signal, 1)
			signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				<-signalChan
				fmt.Println("\nRecibida señal de interrupción. Saliendo...")
				os.Exit(0)
			}()

			// Escuchar cambios en MongoDB y actualizar el caché en Redis
			for cs.Next(context.Background()) {
				var changeEvent bson.M
				if err := cs.Decode(&changeEvent); err != nil {
					log.Println("Error al decodificar el evento de cambio:", err)
					continue
				}

				fmt.Println("Cambio detectado:", changeEvent)

				// Obtener el ID del documento afectado
				var documentID string

				if id, ok := changeEvent["documentKey"].(bson.M)["_id"]; ok {
					documentID = id.(primitive.ObjectID).Hex()
				} else {
					log.Println("No se pudo obtener el ID del documento afectado.")
					continue
				}

				// Actuar en consecuencia del operationType
				operationType := changeEvent["operationType"]
				if operationType == "delete" {
					c.redisClient.HDel("categories", documentID)
					continue
				}

				if operationType == "insert" || operationType == "update" {
					// Obtener el documento actualizado desde MongoDB
					var updatedDocument entity.Category
					_id, _ := primitive.ObjectIDFromHex(documentID)
					if err := c.mongoClient.Database(databaseName).Collection(collectionName).FindOne(context.Background(), bson.M{"_id": _id}).Decode(&updatedDocument); err != nil {
						log.Println("Error al obtener el documento actualizado:", err)
						continue
					}

					// Actualizar el caché en Redis
					// Convertir el documento a una cadena JSON
					jsonString, err := json.Marshal(updatedDocument)
					if err != nil {
						log.Println(err.Error())
						continue
					}

					// Guardar en Redis utilizando el ID del documento como clave
					err = c.redisClient.HSet(collectionName, updatedDocument.ID.Hex(), jsonString).Err()
					if err != nil {
						log.Println(err.Error())
						continue
					}

					fmt.Println("Caché de Redis actualizado:", updatedDocument.ID)
				}

			}
		}(changeStream, collection)
	}
}

func (c *cacheSyncronizer) getChangeStream(databaseName string, collectionName string) (*mongo.ChangeStream, error) {
	return c.mongoClient.Database(databaseName).Collection(collectionName).Watch(context.Background(), mongo.Pipeline{})
}
