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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/diasjorge/dynamokv/serializer"
	"github.com/diasjorge/dynamokv/table"
	"github.com/spf13/cobra"
)

// storeCmd represents the store command
var storeCmd = &cobra.Command{
	Use:   "store TABLENAME CONFIG_FILE",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: store,
	// func(cmd *cobra.Command, args []string) {
	//      // TODO: Work your own magic here
	//      fmt.Println("store called")
	// },
}

func init() {
	RootCmd.AddCommand(storeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// storeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// storeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

// ```yml
// YOUR_KEY: VALUE GOES HERE
// ANOTHER_KEY_SERIALIZED:
//   serialization: 'base64'
//   value: |
//     SOME LONG STRING
//     WITH MULTIPLE LINES
// ANOTHER_KEY_ENCRYPTED:
//   serialization:
//     kms:
//       key: 'KMS_KEY_ID'
//   value: YOUR SECRET VALUE
// ANOTHER_KEY_FROM_FILE:
//   value:
//     file: 'path_to_file'
// ```

func readConfigFile(filename string) (map[string]string, error) {
	return map[string]string{
		"keyname":  "keyvalue",
		"keyname1": "keyvalue1",
	}, nil
}

func store(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return cmd.Usage()
	}

	tableName := args[0]
	configFile := args[1]

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	dynamodb := dynamodb.New(sess, &aws.Config{
		Endpoint: aws.String(endpointURL),
	})
	items, err := readConfigFile(configFile)
	if err != nil {
		return err
	}

	items, err := serializer.SerializeItems(kms, parsedItems)
	if err != nil {
		return err
	}

	table := table.NewTable(dynamodb, tableName)
	if err := table.Create(); err != nil {
		return err
	}
	if err := table.Write(items); err != nil {
		return err
	}
	return nil
}
