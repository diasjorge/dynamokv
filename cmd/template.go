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
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/diasjorge/dynamokv/serializer"
	"github.com/diasjorge/dynamokv/table"
	"github.com/spf13/cobra"
)

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template TABLENAME TEMPLATEFILE [OUTPUTFILE]",
	Short: "Substitute placeholders with values from DynamoDB",
	Long: `Replace placeholders for their value in an AWS DynamoDB table.
Any key in between braces ("{{Key}}") is considered a placeholder.

Placeholders accept the following modifiers:
  {{Key}}
  Example: "{{Username}}" will be replaced by the value of the "Username" Key.
  {{RAW:Key}}
  Example: "{{RAW:Username}}" will be replaced by the value of the "Username" Key without applying deserialization.`,
	RunE: templateParse,
}

func init() {
	RootCmd.AddCommand(templateCmd)
	templateCmd.Flags().BoolVarP(&inplace, "inplace", "i", false, "Replace template file inline")
}

const modRaw = "RAW"

var inplace bool

func templateParse(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("Invalid arguments\n%s", cmd.UsageString())
	}

	tableName, templateFile := args[0], args[1]

	outputFile := ""

	if inplace {
		outputFile = templateFile
	}

	if len(args) == 3 {
		if inplace {
			return fmt.Errorf("OUTPUTFILE and inline flag are mutually exclusive\n%s", cmd.UsageString())
		}
		outputFile = args[2]
	}

	session := newSession(region, profile, endpointURL)

	return template(session, tableName, templateFile, outputFile)
}

func template(session *Session, tableName, templateFile, outputFile string) error {
	template, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return err
	}

	var replaceErrors = &[]error{}
	re := regexp.MustCompile(`{{(\w+?:)?.+?}}`)
	output := re.ReplaceAllFunc(template, generateReplaceFunc(session, tableName, replaceErrors))

	if len(*replaceErrors) > 0 {
		return errors.New("Processing template error")
	}

	if outputFile != "" {
		err := ioutil.WriteFile(outputFile, output, 0644)
		if err != nil {
			return err
		}
	} else {
		fmt.Println(string(output))
	}
	return nil
}

func generateReplaceFunc(session *Session, tableName string, errors *[]error) func([]byte) []byte {
	table := table.NewTable(session.DynamoDB, tableName)
	logger := log.New(os.Stderr, "", 0)

	return func(input []byte) []byte {
		re := regexp.MustCompile(`{{((?P<mod>\w+?):)?(?P<key>.+?)}}`)
		matches := re.FindSubmatch(input)

		var key, mod []byte
		for i, name := range re.SubexpNames() {
			switch name {
			case "mod":
				mod = matches[i]
			case "key":
				key = matches[i]
			}
		}
		parsedItem, err := table.Get(string(key))
		if err != nil {
			logger.Fatal(err)
		}
		deserialize := string(mod) != modRaw
		item, err := serializer.DeserializeItem(session.KMS, parsedItem, deserialize)
		if err != nil {
			logger.Fatal(err)
		}
		return []byte(item.Value)
	}
}
