package cmd

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/diasjorge/dynamokv/models"
)

var export, deserialize bool
var endpointURL, region, profile string

// commandError is an error used to signal different error situations in command handling.
type commandError struct {
	s         string
	userError bool
}

func (c commandError) Error() string {
	return c.s
}

func (c commandError) isUserError() bool {
	return c.userError
}

func newUserError(a ...interface{}) commandError {
	return commandError{s: fmt.Sprintln(a...), userError: true}
}

type Session struct {
	Session  *session.Session
	DynamoDB *dynamodb.DynamoDB
	KMS      *kms.KMS
}

func newSession(region, profile, endpointURL string) *Session {
	config := aws.NewConfig().WithRegion(region)
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:  *config,
		Profile: profile,
	}))

	dynamodbSvc := dynamodb.New(sess, &aws.Config{Endpoint: aws.String(endpointURL)})

	kmsSvc := kms.New(sess)

	return &Session{
		Session:  sess,
		DynamoDB: dynamodbSvc,
		KMS:      kmsSvc,
	}
}

func printItem(item *models.Item, export bool) {
	format := "%s='%s'\n"
	if export {
		format = "export " + format
	}
	fmt.Printf(format, item.Key, escape(item.Value))
}

func escape(value string) string {
	return strings.Replace(value, "'", "'\\''", -1)
}
