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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/spf13/cobra"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch TABLENAME",
	Short: "Retrieve Key Value Pairs from a dynamodb table",
	RunE:  fetch,
}

func init() {
	RootCmd.AddCommand(fetchCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fetchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fetchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

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

	svc := dynamodb.New(sess, &aws.Config{Endpoint: aws.String(endpointURL)})

	params := &dynamodb.ScanInput{
		TableName: aws.String(tableName), // Required
		AttributesToGet: []*string{
			aws.String("Key"),
			aws.String("Value"),
		},
		ConsistentRead: aws.Bool(true),
	}
	resp, err := svc.Scan(params)

	if err != nil {
		return err
	}

	// Pretty-print the response data.
	// fmt.Println(resp)

	for _, item := range resp.Items {
		key := *item["Key"].S
		value := *item["Value"].S
		fmt.Printf("%s=%s\n", key, value)
	}

	return nil
}
