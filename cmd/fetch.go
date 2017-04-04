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

	"github.com/diasjorge/dynamokv/serializer"
	"github.com/diasjorge/dynamokv/table"
	"github.com/spf13/cobra"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch TABLENAME",
	Short: "Retrieve All Key Value Pairs",
	Long:  "Retrieve All Key Value Pairs from a DynamoDB table",
	RunE:  fetchParse,
}

func init() {
	RootCmd.AddCommand(fetchCmd)
	fetchCmd.Flags().BoolVarP(&export, "export", "", false, "Export variables")
	fetchCmd.Flags().BoolVarP(&deserialize, "deserialize", "", true, "Deserialize items")
}

func fetchParse(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("TABLENAME required")
	}

	tableName := args[0]

	session := newSession(region, profile, endpointURL)

	return fetch(session, tableName, export, deserialize)
}

func fetch(session *Session, tableName string, export, deserialize bool) error {
	table := table.NewTable(session.DynamoDB, tableName)

	parsedItems, err := table.Read()
	if err != nil {
		return err
	}

	items, err := serializer.DeserializeItems(session.KMS, parsedItems, deserialize)
	if err != nil {
		return err
	}

	for _, item := range items {
		printItem(item, export)
	}

	return nil
}
