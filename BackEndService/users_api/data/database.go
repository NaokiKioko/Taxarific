package data

import (
	"context"
	"os"
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
	uri := "mongodb+srv://" + os.Getenv("MONGO_USER") + ":" + os.Getenv("MONGO_PASS") + "@taxarific.wxqpl.mongodb.net/"
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
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

// TODO Users
func CreateUser(user models.User) (string, error) {
	result, err := userCollection().InsertOne(context.Background(), &user)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func GetUser(id string) (models.User, error) {
	objId, err := GetObjectID(id)
	if err != nil {
		return models.User{}, err
	}
	var user models.User
	err = userCollection().FindOne(context.Background(), bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func GetUsers() ([]models.User, error) {
	var users []models.User
	cursor, err := userCollection().Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var user models.User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func UpdateUser(id string, user models.User) error {
	objId, err := GetObjectID(id)
	if err != nil {
		return err
	}
	_, err = userCollection().UpdateOne(context.Background(), bson.M{"_id": objId}, bson.M{"$set": user})
	if err != nil {
		return err
	}
	return nil
}

func PutUser(id string, user models.User) error {
	objId, err := GetObjectID(id)
	if err != nil {
		return err
	}
	_, err = userCollection().UpdateOne(context.Background(), bson.M{"_id": objId}, bson.M{"$set": user})
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(id string) error {
	objId, err := GetObjectID(id)
	if err != nil {
		return err
	}
	_, err = userCollection().DeleteOne(context.Background(), bson.M{"_id": objId})
	if err != nil {
		return err
	}
	return nil
}

// TODO Employees
func CreateEmployee(employee models.Employee) (string, error) {
	result, err := employeeCollection().InsertOne(context.Background(), &employee)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func GetEmployee(id string) (models.Employee, error) {
	objId, err := GetObjectID(id)
	if err != nil {
		return models.Employee{}, err
	}
	var employee models.Employee
	err = employeeCollection().FindOne(context.Background(), bson.M{"_id": objId}).Decode(&employee)
	if err != nil {
		return models.Employee{}, err
	}
	return employee, nil
}

func GetEmployees() ([]models.Employee, error) {
	var employees []models.Employee
	cursor, err := employeeCollection().Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var employee models.Employee
		err := cursor.Decode(&employee)
		if err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, nil
}

func PutEmployee(id string, employee models.Employee) error {
	objId, err := GetObjectID(id)
	if err != nil {
		return err
	}
	_, err = employeeCollection().UpdateOne(context.Background(), bson.M{"_id": objId}, bson.M{"$set": employee})
	if err != nil {
		return err
	}
	return nil
}

func DeleteEmployee(id string) error {
	objId, err := GetObjectID(id)
	if err != nil {
		return err
	}
	_, err = employeeCollection().DeleteOne(context.Background(), bson.M{"_id": objId})
	if err != nil {
		return err
	}
	return nil
}

// TODO Admins
func CreateAdmin(admin models.Admin) (string, error) {
	result, err := adminCollection().InsertOne(context.Background(), &admin)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func GetAdmins() ([]models.Admin, error) {
	var admins []models.Admin
	cursor, err := adminCollection().Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var admin models.Admin
		err := cursor.Decode(&admin)
		if err != nil {
			return nil, err
		}
		admins = append(admins, admin)
	}
	return admins, nil
}

func GetAdmin(id string) (models.Admin, error) {
	objId, err := GetObjectID(id)
	if err != nil {
		return models.Admin{}, err
	}
	var admin models.Admin
	err = adminCollection().FindOne(context.Background(), bson.M{"_id": objId}).Decode(&admin)
	if err != nil {
		return models.Admin{}, err
	}
	return admin, nil
}

func PutAdmin(id string, admin models.Admin) error {
	objId, err := GetObjectID(id)
	if err != nil {
		return err
	}
	_, err = adminCollection().UpdateOne(context.Background(), bson.M{"_id": objId}, bson.M{"$set": admin})
	if err != nil {
		return err
	}
	return nil
}

func DeleteAdmin(id string) error {
	objId, err := GetObjectID(id)
	if err != nil {
		return err
	}
	_, err = adminCollection().DeleteOne(context.Background(), bson.M{"_id": objId})
	if err != nil {
		return err
	}
	return nil
}

// Helper functions
func GetObjectID(id string) (primitive.ObjectID, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return objId, nil
}
