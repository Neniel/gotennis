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

func SyncMongoCache(mongoClient *mongo.Client, redisClient *redis.Client) {
	// Configurar el cambio de transmisión (change stream) en MongoDB
	changeStream, err := mongoClient.Database("tennis").Collection("categories").Watch(context.Background(), mongo.Pipeline{})
	if err != nil {
		log.Fatal("Error al configurar el cambio de transmisión en MongoDB:", err)
	}

	// Manejar señales de interrupción para salir graciosamente
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		fmt.Println("\nRecibida señal de interrupción. Saliendo...")
		os.Exit(0)
	}()

	// Escuchar cambios en MongoDB y actualizar el caché en Redis
	for changeStream.Next(context.Background()) {
		var changeEvent bson.M
		if err := changeStream.Decode(&changeEvent); err != nil {
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
			redisClient.HDel("categories", documentID)
			continue
		}

		if operationType == "insert" || operationType == "update" {
			// Obtener el documento actualizado desde MongoDB
			var updatedDocument entity.Category
			_id, _ := primitive.ObjectIDFromHex(documentID)
			if err := mongoClient.Database("tennis").Collection("categories").FindOne(context.Background(), bson.M{"_id": _id}).Decode(&updatedDocument); err != nil {
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
			err = redisClient.HSet("categories", updatedDocument.ID.Hex(), jsonString).Err()
			if err != nil {
				log.Println(err.Error())
				continue
			}

			fmt.Println("Caché de Redis actualizado:", updatedDocument.ID)
		}
	}
}
