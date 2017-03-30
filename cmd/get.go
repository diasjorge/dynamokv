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

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get TABLENAME KEY",
	Short: "Retrieve Value of Key",
	RunE:  get,
}

func init() {
	RootCmd.AddCommand(getCmd)

	getCmd.Flags().BoolVarP(&export, "export", "", false, "Export variables")
	getCmd.Flags().BoolVarP(&deserialize, "deserialize", "", true, "Deserialize items")
}

func get(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return errors.New("TABLENAME KEY required")
	}

	tableName := args[0]
	key := args[1]

	session := newSession()

	table := table.NewTable(session.DynamoDB, tableName)

	item, err := table.Get(key)
	if err != nil {
		return err
	}

	err = serializer.DeserializeItem(session.KMS, item, deserialize)
	if err != nil {
		return err
	}

	printItem(item, export)
	return nil
}
