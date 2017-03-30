// Copyright © 2017 Jorge Dias <jorge@mrdias.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/diasjorge/dynamokv/models"
	"github.com/diasjorge/dynamokv/serializer"
	"github.com/diasjorge/dynamokv/table"
	"github.com/spf13/cobra"
)

var export, deserialize bool

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch TABLENAME",
	Short: "Retrieve Key Value Pairs from a dynamodb table",
	RunE:  fetch,
}

func init() {
	RootCmd.AddCommand(fetchCmd)
	fetchCmd.Flags().BoolVarP(&export, "export", "", false, "Export variables")
	fetchCmd.Flags().BoolVarP(&deserialize, "deserialize", "", true, "Deserialize items")
}

func fetch(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("TABLENAME required")
	}

	tableName := args[0]

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String(region)},
		Profile: profile,
	}))

	dynamodbSvc := dynamodb.New(sess, &aws.Config{Endpoint: aws.String(endpointURL)})

	kmsSvc := kms.New(sess)

	table := table.NewTable(dynamodbSvc, tableName)

	items, err := table.Read()
	if err != nil {
		return err
	}

	items, err = serializer.DeserializeItems(kmsSvc, items, deserialize)
	if err != nil {
		return err
	}

	for _, item := range items {
		printItem(item)
	}

	return nil
}

func printItem(item *models.Item) {
	format := "%s='%s'\n"
	if export {
		format = "export " + format
	}
	fmt.Printf(format, item.Key, escape(item.Value))
}

func escape(value string) string {
	return strings.Replace(value, "'", "'\\''", -1)
}
