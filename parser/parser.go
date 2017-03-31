package parser

import (
	"io/ioutil"
	"strings"

	"github.com/diasjorge/dynamokv/models"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

type rawValue struct {
	RawSerialization interface{} `mapstructure:"serialization"`
	RawValue         interface{} `mapstructure:"value"`
}

func (rawValue *rawValue) parseSerialization() (*models.Serialization, error) {
	serialization := models.NewSerialization()

	switch rawValue.RawSerialization.(type) {
	case string:
		serialization.Type = rawValue.RawSerialization.(string)
	case interface{}:
		if err := mapstructure.Decode(rawValue.RawSerialization, &serialization); err != nil {
			return nil, err
		}
	}
	return serialization, nil
}

func (rawValue *rawValue) parseValue() (string, error) {
	var value string

	switch rawValue.RawValue.(type) {
	case string:
		value = rawValue.RawValue.(string)
	case interface{}:
		var valueOptions map[string]string
		if err := mapstructure.Decode(rawValue.RawValue, &valueOptions); err != nil {
			return "", err
		}
		content, err := ioutil.ReadFile(valueOptions["file"])
		if err != nil {
			return "", err
		}
		value = strings.Trim(string(content[:]), "\n")
	}
	return value, nil
}

func (rawValue *rawValue) Parse() (*models.ParsedItemValue, error) {
	serialization, err := rawValue.parseSerialization()
	if err != nil {
		return nil, err
	}
	value, err := rawValue.parseValue()
	if err != nil {
		return nil, err
	}
	itemValue := &models.ParsedItemValue{Serialization: serialization, Value: value}
	return itemValue, nil
}

// Parse returns Items from yaml file
func Parse(filename string) ([]*models.ParsedItem, error) {
	var keyVal map[string]interface{}
	var items []*models.ParsedItem

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, &keyVal); err != nil {
		return nil, err
	}
	for key, value := range keyVal {
		item := models.NewParsedItem()
		item.Key = key
		switch value.(type) {
		case string:
			item.Value.Value = value.(string)
		case interface{}:
			var rawValue rawValue
			if err := mapstructure.Decode(value, &rawValue); err != nil {
				return nil, err
			}
			parsedValue, err := rawValue.Parse()
			if err != nil {
				return nil, err
			}
			item.Value = parsedValue
		}
		items = append(items, item)
	}
	return items, nil
}
