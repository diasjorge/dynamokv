package models

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Item struct {
	Key           string
	Value         string
	Serialization string
}

func NewItemFromDynamoDB(item map[string]*dynamodb.AttributeValue) (*Item, error) {
	key, ok := item["Key"]
	if !ok {
		return nil, fmt.Errorf("Missing Key attribute for item: %v", item)
	}
	value, ok := item["Value"]
	if !ok {
		return nil, fmt.Errorf("Missing Value attribute for item: %v", item)
	}
	serializationType := "plain"
	serialization, ok := item["Serialization"]
	if ok {
		serializationType = *serialization.S
	}
	return &Item{
		Key:           *key.S,
		Value:         *value.S,
		Serialization: serializationType,
	}, nil
}
