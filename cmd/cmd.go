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

func newSession() *Session {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String(region)},
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
