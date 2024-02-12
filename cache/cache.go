package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Neniel/gotennis/app"
	"github.com/Neniel/gotennis/entity"
	"github.com/Neniel/gotennis/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CacheMicroservice struct {
	App *app.App
}

func (ms *CacheMicroservice) StartSync() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	for _, collection := range util.Collections {
		changeStream, err := ms.App.DBClients.MongoDB.Database(util.DBName).Collection(collection).Watch(ctx, mongo.Pipeline{})
		if err != nil {
			cancel()
			panic(err)
		}

		wg.Add(1)
		go func(cs *mongo.ChangeStream, wg *sync.WaitGroup, collectionName string) {
			defer wg.Done()
			// Manejar señales de interrupción para salir graciosamente
			signalChan := make(chan os.Signal, 1)
			signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				<-signalChan
				fmt.Println("\nRecibida señal de interrupción. Saliendo...")
				os.Exit(0)
			}()

			// Escuchar cambios en MongoDB y actualizar el caché en Redis
			log.Printf("Listening for changes on '%s' collection", collectionName)
			for cs.Next(ctx) {
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
					ms.App.DBClients.Redis.HDel("categories", documentID)
					continue
				}

				if operationType == "insert" || operationType == "update" {
					// Obtener el documento actualizado desde MongoDB
					var updatedDocument entity.Category
					_id, _ := primitive.ObjectIDFromHex(documentID)
					if err := ms.App.DBClients.MongoDB.Database(util.DBName).Collection(collectionName).FindOne(ctx, bson.M{"_id": _id}).Decode(&updatedDocument); err != nil {
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
					err = ms.App.DBClients.Redis.HSet(collectionName, updatedDocument.ID.Hex(), jsonString).Err()
					if err != nil {
						log.Println(err.Error())
						continue
					}

					fmt.Println("Caché de Redis actualizado:", updatedDocument.ID)
				}

			}
		}(changeStream, &wg, collection)
	}
	wg.Wait()
}
