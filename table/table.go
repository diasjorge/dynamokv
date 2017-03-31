package table

import (
	"fmt"

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

func (table *Table) Read() ([]*models.ParsedItem, error) {
	params := &dynamodb.ScanInput{
		TableName: table.Name,
		AttributesToGet: []*string{
			aws.String("Key"),
			aws.String("Value"),
			aws.String("Serialization"),
		},
	}
	items := []*models.ParsedItem{}

	err := table.svc.ScanPages(
		params,
		func(resp *dynamodb.ScanOutput, lastPage bool) bool {
			for _, dynamodbItem := range resp.Items {
				item, err := models.NewParsedItemFromDynamoDB(dynamodbItem)
				if err != nil {
					continue
				}
				items = append(items, item)
			}
			return true
		},
	)

	if err != nil {
		return nil, err
	}
	return items, nil
}

func (table *Table) Get(key string) (*models.ParsedItem, error) {
	params := &dynamodb.QueryInput{
		TableName: table.Name,
		AttributesToGet: []*string{
			aws.String("Key"),
			aws.String("Value"),
			aws.String("Serialization"),
		},
		KeyConditions: map[string]*dynamodb.Condition{
			"Key": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(key),
					},
				},
			},
		},
	}

	resp, err := table.svc.Query(params)
	if err != nil {
		return nil, err
	}

	if *resp.Count != 1 {
		return nil, fmt.Errorf("error querying for Item with Key \"%v\": %v occurrences found", key, *resp.Count)
	}

	return models.NewParsedItemFromDynamoDB(resp.Items[0])
}
