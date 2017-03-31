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
	"fmt"
	"strings"

	"github.com/diasjorge/dynamokv/models"
	"github.com/diasjorge/dynamokv/serializer"
	"github.com/diasjorge/dynamokv/table"
	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set TABLENAME KEY VALUE",
	Short: "Set Value for Key",
	RunE:  set,
}

type serializationFlag struct {
	stype   string
	options map[string]string
}

func (serialization *serializationFlag) String() string {
	return "plain"
}

func (serialization *serializationFlag) Type() string {
	return "serialization"
}

func (serialization *serializationFlag) Set(value string) error {
	// "kms::key:alias/value"
	typeOptions := strings.Split(value, "::")

	if len(typeOptions) < 1 {
		return fmt.Errorf("invalid serialization format. Expected: type::option:optionValue,option2:optionValue2")
	}

	serialization.stype = typeOptions[0]

	if len(typeOptions) > 1 {
		options := strings.Split(typeOptions[1], ",")
		serialization.options = map[string]string{}
		for _, option := range options {
			keyVal := strings.Split(option, ":")
			if len(keyVal) != 2 {
				return fmt.Errorf("invalid serialization format. Expected: type::option:optionValue,option2:optionValue2")
			}
			serialization.options[keyVal[0]] = keyVal[1]
		}
	}
	return nil
}

var serializationF serializationFlag

func init() {
	RootCmd.AddCommand(setCmd)
	setCmd.Flags().VarP(&serializationF, "serialization", "", "type::option:optionValue,*")
}

func set(cmd *cobra.Command, args []string) error {
	if len(args) != 3 {
		return fmt.Errorf("Invalid arguments\n%s", cmd.UsageString())
	}
	tableName, key, value := args[0], args[1], args[2]

	session := newSession()

	parsedItem := models.NewParsedItem()
	parsedItem.Key = key
	parsedItem.Value.Value = value

	if serializationF.stype != "" {
		parsedItem.Value.Serialization.Type = serializationF.stype
		parsedItem.Value.Serialization.Options = serializationF.options
	}

	item, err := serializer.SerializeItem(session.KMS, parsedItem)
	if err != nil {
		return err
	}

	table := table.NewTable(session.DynamoDB, tableName)
	if err := table.Create(); err != nil {
		return err
	}

	err = table.Set(item)
	if err != nil {
		return err
	}
	return nil
}
