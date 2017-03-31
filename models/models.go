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

type ParsedItem struct {
	Key   string
	Value *ParsedItemValue
}

type ParsedItemValue struct {
	Value         string
	Serialization *Serialization
}

type Serialization struct {
	Type    string
	Options map[string]string
}

func NewItem() *Item {
	return &Item{
		Serialization: "plain",
	}
}

func NewParsedItem() *ParsedItem {
	return &ParsedItem{
		Value: &ParsedItemValue{
			Serialization: NewSerialization(),
		},
	}
}

func NewSerialization() *Serialization {
	return &Serialization{Type: "plain"}
}

func NewParsedItemFromDynamoDB(dynamodbItem map[string]*dynamodb.AttributeValue) (*ParsedItem, error) {
	item := NewParsedItem()
	key, ok := dynamodbItem["Key"]
	if !ok {
		return nil, fmt.Errorf("Missing Key attribute for item: %v", dynamodbItem)
	}
	item.Key = *key.S
	value, ok := dynamodbItem["Value"]
	if !ok {
		return nil, fmt.Errorf("Missing Value attribute for item: %v", dynamodbItem)
	}
	item.Value.Value = *value.S
	serialization, ok := dynamodbItem["Serialization"]
	if ok {
		item.Value.Serialization.Type = *serialization.S
	}
	return item, nil
}

// func NewItemFromVars(key, value string, options map[string]string) (*Item, error) {
// 	return &Item{
// 		Key:           key,
// 		Value:         value,
// 		Serialization: serializationType,
// 	}, nil
// }
