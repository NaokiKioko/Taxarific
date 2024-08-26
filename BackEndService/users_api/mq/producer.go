package mq

import (
	"encoding/json"
	"taxarific_users_api/models"
	"log"

	"github.com/IBM/sarama"
)


func AccountCreationEmail(user *models.User) {
	producer, err := sarama.NewSyncProducer([]string{"kafka:9092"}, nil)
	if err != nil {
		log.Fatalf("couldnt create producer: {%s}", err)
	}
	userJson, err := json.Marshal(&user)
	if err != nil {
		log.Fatalf("couldnt Marshal json: {%s}", err)
	}
	msg := &sarama.ProducerMessage{
		Topic: "accountCreation",
		Key:   sarama.StringEncoder("account_creation"),
		Value: sarama.StringEncoder(userJson),
	}
	producer.SendMessage(msg)
}
