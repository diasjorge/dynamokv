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
	"github.com/diasjorge/dynamokv/parser"
	"github.com/diasjorge/dynamokv/serializer"
	"github.com/diasjorge/dynamokv/table"
	"github.com/spf13/cobra"
)

// storeCmd represents the store command
var storeCmd = &cobra.Command{
	Use:   "store TABLENAME CONFIG_FILE",
	Short: "Store Key Value pairs and optionally serialize them into a dynamodb table",
	RunE:  storeParse,
}

func init() {
	RootCmd.AddCommand(storeCmd)
}

func storeParse(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return cmd.Usage()
	}

	tableName := args[0]
	configFile := args[1]

	session := newSession(region, profile, endpointURL)

	return store(session, tableName, configFile)
}

func store(session *Session, tableName, configFile string) error {
	parsedItems, err := parser.Parse(configFile)
	if err != nil {
		return err
	}

	items, err := serializer.SerializeItems(session.KMS, parsedItems)
	if err != nil {
		return err
	}

	table := table.NewTable(session.DynamoDB, tableName)
	if err := table.Create(); err != nil {
		return err
	}
	if err := table.Write(items); err != nil {
		return err
	}
	return nil
}
