package table

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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

func (table *Table) Write(items map[string]string) error {
	writeRequests := []*dynamodb.WriteRequest{}
	for key, value := range items {
		writeRequests = append(writeRequests, &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: map[string]*dynamodb.AttributeValue{
					"Key": {
						S: aws.String(key),
					},
					"Value": {
						S: aws.String(value),
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
