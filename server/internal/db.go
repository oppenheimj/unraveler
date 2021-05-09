package internal

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GraphDocument is a mongo graph document
type GraphDocument struct {
	E        map[string][]int   `bson:"edges"`
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	Ka       float64            `json:"ka"`
	Kr       float64            `json:"kr"`
	Kn       float64            `json:"kn"`
	MaxIters int                `bson:"max_iters"`
	MinError float64            `bson:"min_error"`
}

// FindGraph finds one graph document in mongo
func FindGraph(filter bson.D) *GraphDocument {
	clientOptions := options.Client().ApplyURI("mongodb://gd-mongo:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	collection := client.Database("gd_data").Collection("graphs")

	var result GraphDocument
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Found graph", result.ID)
	fmt.Println(result)
	disconnectClient(client)

	return &result
}

// UpdateGraph updates one graph document in mongo
func UpdateGraph(filter, update bson.D) {
	clientOptions := options.Client().ApplyURI("mongodb://gd-mongo:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	collection := client.Database("gd_data").Collection("graphs")

	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	disconnectClient(client)
}

func disconnectClient(client *mongo.Client) {
	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}
