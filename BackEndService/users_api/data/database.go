package data

import (
	"context"
	"taxarific_users_api/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func NewDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	// client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://database:27017"))
	if err != nil {
		return err
	}
	return nil
}

func userCollection() *mongo.Collection {
	return client.Database("Taxarific").Collection("users")
}

func employeeCollection() *mongo.Collection {
	return client.Database("Taxarific").Collection("employees")
}

func adminCollection() *mongo.Collection {
	return client.Database("Taxarific").Collection("admins")
}



// Helper functions
func GetObjectID(id string) (primitive.ObjectID, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return objId, nil
}