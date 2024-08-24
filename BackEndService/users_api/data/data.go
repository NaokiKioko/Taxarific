package data

import (
	"context"
	"errors"
	"fmt"
	"os"
	"taxarific_users_api/auth"
	"taxarific_users_api/models"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func NewDB() error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}
	uri := os.Getenv("MONGO_CON_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	defer client.Disconnect(context.TODO())
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		return err
	}
	fmt.Println("Successfully connected to Atlas")
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

func Userlogin(login *models.PostLoginJSONRequestBody) (string, error) {
	var user models.User
	err := userCollection().FindOne(context.Background(), bson.M{"email": login.Email}).Decode(&user)
	if err != nil {
		return "", errors.New("user not found")
	}
	if err := auth.CheckPassword(user.Password, login.Password); err != nil {
		return "", errors.New("invalid password")
	}
	token, err := auth.GenerateJWTToken(user.Id.Hex(), "user")
	if err != nil {
		return "", err
	}
	return token, nil
}

func AdminLogin(login *models.PostLoginJSONRequestBody) (string, error) {
	var admin models.Admin
	err := adminCollection().FindOne(context.Background(), bson.M{"email": login.Email}).Decode(&admin)
	if err != nil {
		return "", errors.New("admin not found")
	}
	if err := auth.CheckPassword(admin.Password, login.Password); err != nil {
		return "", errors.New("invalid password")
	}
	token, err := auth.GenerateJWTToken(admin.Id.Hex(), "admin")
	if err != nil {
		return "", err
	}
	return token, nil
}

func EmployeeLogin(login *models.PostLoginJSONRequestBody) (string, error) {
	var employee models.Employee
	err := employeeCollection().FindOne(context.Background(), bson.M{"email": login.Email}).Decode(&employee)
	if err != nil {
		return "", errors.New("employee not found")
	}
	if err := auth.CheckPassword(employee.Password, login.Password); err != nil {
		return "", errors.New("invalid password")
	}
	token, err := auth.GenerateJWTToken(employee.Id.Hex(), "employee")
	if err != nil {
		return "", err
	}
	return token, nil
}

// Users
func CreateUser(user *models.PostUserJSONRequestBody) (string, error) {
	insertedId, err := userCollection().InsertOne(context.Background(), &user)
	if err != nil {
		return "", err
	}
	token, err := auth.GenerateJWTToken(insertedId.InsertedID.(primitive.ObjectID).Hex(), "user")
	return token, nil
}

// func GetUser(id string) (models.User, error) {
// 	objId, err := GetObjectID(id)
// 	if err != nil {
// 		return models.User{}, err
// 	}
// 	var user models.User
// 	err = userCollection().FindOne(context.Background(), bson.M{"_id": objId}).Decode(&user)
// 	if err != nil {
// 		return models.User{}, err
// 	}
// 	return user, nil
// }

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

// func UpdateUser(id string, user *models.User) error {
// 	objId, err := GetObjectID(id)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = userCollection().UpdateOne(context.Background(), bson.M{"_id": objId}, bson.M{"$set": user})
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func PutUser(id string, user *models.User) error {
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

// Employees
func CreateEmployee(employee *models.PostAdminEmployeeJSONRequestBody) error {
	_, err := employeeCollection().InsertOne(context.Background(), &employee)
	if err != nil {
		return err
	}
	return nil
}

func GetEmployee(id string) (*models.Employee, error) {
	objId, err := GetObjectID(id)
	if err != nil {
		return &models.Employee{}, err
	}
	var employee models.Employee
	err = employeeCollection().FindOne(context.Background(), bson.M{"_id": objId}).Decode(&employee)
	if err != nil {
		return &models.Employee{}, err
	}
	return &employee, nil
}

func GetEmployees() (*[]models.Employee, error) {
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
	return &employees, nil
}

func PutEmployee(id string, employee *models.Employee) error {
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

// Admins
func CreateAdmin(admin *models.PostAdminJSONRequestBody) error {
	_, err := adminCollection().InsertOne(context.Background(), &admin)
	if err != nil {
		return err
	}
	return nil
}

func GetAdmins() (*[]models.Admin, error) {
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
	return &admins, nil
}

func GetAdmin(id string) (*models.Admin, error) {
	objId, err := GetObjectID(id)
	if err != nil {
		return &models.Admin{}, err
	}
	var admin models.Admin
	err = adminCollection().FindOne(context.Background(), bson.M{"_id": objId}).Decode(&admin)
	if err != nil {
		return &models.Admin{}, err
	}
	return &admin, nil
}

func PutAdmin(id string, admin *models.Admin) error {
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
func GetObjectID(id string) (*primitive.ObjectID, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &primitive.NilObjectID, err
	}
	return &objId, nil
}
