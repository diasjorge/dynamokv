package parser

import (
	"io/ioutil"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

type Item struct {
	Key   string
	Value *ItemValue
}

type ItemValue struct {
	Value         string
	Serialization *Serialization
}

type Serialization struct {
	Type    string
	Options map[string]string
}

func newSerialization() *Serialization {
	return &Serialization{Type: "plain"}
}

type rawValue struct {
	RawSerialization interface{} `mapstructure:"serialization"`
	RawValue         interface{} `mapstructure:"value"`
}

func (rawValue *rawValue) parseSerialization() (*Serialization, error) {
	serialization := newSerialization()

	switch rawValue.RawSerialization.(type) {
	case string:
		serialization = &Serialization{
			Type: rawValue.RawSerialization.(string),
		}
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
		value = string(content[:])
	}
	return value, nil
}

func (rawValue *rawValue) Parse() (*ItemValue, error) {
	serialization, err := rawValue.parseSerialization()
	if err != nil {
		return nil, err
	}
	value, err := rawValue.parseValue()
	if err != nil {
		return nil, err
	}
	itemValue := &ItemValue{Serialization: serialization, Value: value}
	return itemValue, nil
}

// Parse returns Items from yaml file
func Parse(filename string) ([]*Item, error) {
	var keyVal map[string]interface{}
	var items []*Item

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, &keyVal); err != nil {
		return nil, err
	}
	for key, value := range keyVal {
		item := &Item{Key: key}
		switch value.(type) {
		case string:
			item.Value = &ItemValue{
				Value:         value.(string),
				Serialization: newSerialization(),
			}
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
