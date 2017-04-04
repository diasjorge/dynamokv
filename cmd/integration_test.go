package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

const (
	testTableName = "TEST_TABLE"
	testRegion    = "eu-west-1"
)

var testEndpointURL = os.Getenv("DYNAMODB_URL")

func writeConfig(config string) string {
	file, err := ioutil.TempFile("", "test")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.Write([]byte(config))
	if err != nil {
		panic(err)
	}
	path, err := filepath.Abs(file.Name())
	if err != nil {
		panic(err)
	}
	return path
}

func storeTestConfig(session *Session) {
	config := `
KEY: VALUE
SERIALIZED_KEY:
  serialization: 'base64'
  value: VALUE
`
	configPath := writeConfig(config)

	deleteTable()

	store(session, testTableName, configPath)
}

func captureStdout(f func()) []byte {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	return out
}

func deleteTable() {
	session := newSession(testRegion, "", testEndpointURL)
	_, err := session.DynamoDB.DescribeTable(&dynamodb.DescribeTableInput{TableName: aws.String(testTableName)})
	if err == nil {
		session.DynamoDB.DeleteTable(&dynamodb.DeleteTableInput{
			TableName: aws.String(testTableName),
		})
		session.DynamoDB.WaitUntilTableNotExists(&dynamodb.DescribeTableInput{
			TableName: aws.String(testTableName),
		})
	}
}

func TestMain(m *testing.M) {
	if testEndpointURL == "" {
		fmt.Println("DYNAMODB_URL not set. Skipping Integration tests.")
		os.Exit(0)
	}
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestFetch(t *testing.T) {
	session := newSession(testRegion, "", testEndpointURL)

	storeTestConfig(session)

	out := captureStdout(func() {
		fetch(session, testTableName, false, true)
	})
	expectedOut := "KEY='VALUE'\nSERIALIZED_KEY='VALUE'\n"

	assert.Equal(t, expectedOut, string(out))
}

func TestFetchNoDeserialize(t *testing.T) {
	session := newSession(testRegion, "", testEndpointURL)

	storeTestConfig(session)

	out := captureStdout(func() {
		fetch(session, testTableName, false, false)
	})
	expectedOut := "KEY='VALUE'\nSERIALIZED_KEY='VkFMVUU='\n"

	assert.Equal(t, expectedOut, string(out))
}

func TestFetchExport(t *testing.T) {
	session := newSession(testRegion, "", testEndpointURL)

	storeTestConfig(session)

	out := captureStdout(func() {
		fetch(session, testTableName, true, true)
	})
	expectedOut := "export KEY='VALUE'\nexport SERIALIZED_KEY='VALUE'\n"

	assert.Equal(t, expectedOut, string(out))
}

func TestGet(t *testing.T) {
	session := newSession(testRegion, "", testEndpointURL)

	deleteTable()
	set(session, testTableName, "SINGLE_KEY", "SINGLE_VALUE", "base64", map[string]string{})

	out := captureStdout(func() {
		get(session, testTableName, "SINGLE_KEY", false, true)
	})
	expectedOut := "SINGLE_KEY='SINGLE_VALUE'\n"

	assert.Equal(t, expectedOut, string(out))
}

func TestGetNoDeserialize(t *testing.T) {
	session := newSession(testRegion, "", testEndpointURL)

	deleteTable()
	set(session, testTableName, "SINGLE_KEY", "SINGLE_VALUE", "base64", map[string]string{})

	out := captureStdout(func() {
		get(session, testTableName, "SINGLE_KEY", false, false)
	})
	expectedOut := "SINGLE_KEY='U0lOR0xFX1ZBTFVF'\n"

	assert.Equal(t, expectedOut, string(out))
}

func TestGetExport(t *testing.T) {
	session := newSession(testRegion, "", testEndpointURL)

	deleteTable()
	set(session, testTableName, "SINGLE_KEY", "SINGLE_VALUE", "base64", map[string]string{})

	out := captureStdout(func() {
		get(session, testTableName, "SINGLE_KEY", true, true)
	})
	expectedOut := "export SINGLE_KEY='SINGLE_VALUE'\n"

	assert.Equal(t, expectedOut, string(out))
}
