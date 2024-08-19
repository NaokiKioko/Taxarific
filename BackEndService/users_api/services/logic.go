package services

import (
	"errors"
	"reflect"
	"taxarific_users_api/data"
	"taxarific_users_api/models"

	"go.mongodb.org/mongo-driver/internal/uuid"
)

func GetEmployee(id string) (*models.Employee, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	result, err := data.GetEmployee(id)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func GetEmployees() (*[]models.Employee, error) {
	result, err := data.GetEmployees()
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func GetAdmin(id string) (*models.Admin, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	result, err := data.GetAdmin(id)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func GetAdmins() (*[]models.Admin, error) {
	result, err := data.GetAdmins()
	if err != nil {
		return nil, err 
	}
	return &result, nil
}

func GetUser(id string) (*models.User, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	result, err := data.GetUser(id)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func GetUsers() (*[]models.User, error) {
	result, err := data.GetUsers()
	if err != nil {
		return nil, err
	}
	return &result, nil 
} 

func PutUser(id uuid.UUID, user models.User) error {
	if id == "" {
		return errors.New("id is required")
	}
	if reflect.DeepEqual(user, models.User{}) {
		return errors.New("user is required")
	}
	return data.PutUser(id, user)
}