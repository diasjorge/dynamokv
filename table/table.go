package table

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/diasjorge/dynamokv/models"
)

type Table struct {
	svc  *dynamodb.DynamoDB
	Name *string
}

func NewTable(svc *dynamodb.DynamoDB, name string) *Table {
	return &Table{
		svc:  svc,
		Name: aws.String(name),
	}
}

func (table *Table) Create() error {
	_, err := table.svc.DescribeTable(&dynamodb.DescribeTableInput{TableName: table.Name})
	if err == nil {
		return nil
	}
	_, err = table.svc.CreateTable(&dynamodb.CreateTableInput{TableName: table.Name,
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Key"),
				KeyType:       aws.String("HASH"),
			},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Key"),
				AttributeType: aws.String("S"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(2),
			WriteCapacityUnits: aws.Int64(1),
		},
	})
	if err != nil {
		return err
	}
	if err = table.svc.WaitUntilTableExists(&dynamodb.DescribeTableInput{TableName: table.Name}); err != nil {
		return err
	}
	return nil
}

func (table *Table) Write(items []*models.Item) error {
	writeRequests := []*dynamodb.WriteRequest{}
	for _, item := range items {
		writeRequests = append(writeRequests, &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: map[string]*dynamodb.AttributeValue{
					"Key": {
						S: aws.String(item.Key),
					},
					"Value": {
						S: aws.String(item.Value),
					},
					"Serialization": {
						S: aws.String(item.Serialization),
					},
				},
			},
		})
	}
	_, err := table.svc.BatchWriteItem(&dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			*table.Name: writeRequests,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (table *Table) Read() ([]*models.Item, error) {
	params := &dynamodb.ScanInput{
		TableName: table.Name,
		AttributesToGet: []*string{
			aws.String("Key"),
			aws.String("Value"),
			aws.String("Serialization"),
		},
		ConsistentRead: aws.Bool(true),
	}
	items := []*models.Item{}

	err := table.svc.ScanPages(
		params,
		func(resp *dynamodb.ScanOutput, lastPage bool) bool {
			for _, item := range resp.Items {
				key, ok := item["Key"]
				if !ok {
					continue
				}
				value, ok := item["Value"]
				if !ok {
					continue
				}
				serializationType := "plain"
				serialization, ok := item["Serialization"]
				if ok {
					serializationType = *serialization.S
				}
				items = append(items, &models.Item{
					Key:           *key.S,
					Value:         *value.S,
					Serialization: serializationType,
				})
			}
			return true
		},
	)

	if err != nil {
		return nil, err
	}
	return items, nil
}
