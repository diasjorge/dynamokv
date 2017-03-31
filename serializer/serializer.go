package serializer

import (
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/diasjorge/dynamokv/models"
)

func SerializeItems(svc *kms.KMS, parsedItems []*models.ParsedItem) ([]*models.Item, error) {
	result := []*models.Item{}
	for _, parsedItem := range parsedItems {
		item, err := SerializeItem(svc, parsedItem)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, nil
}

func SerializeItem(svc *kms.KMS, parsedItem *models.ParsedItem) (*models.Item, error) {
	value, err := serialize(svc, parsedItem.Value)
	if err != nil {
		return nil, err
	}
	return &models.Item{
		Key:           parsedItem.Key,
		Value:         value,
		Serialization: parsedItem.Value.Serialization.Type,
	}, nil
}

func DeserializeItems(svc *kms.KMS, parsedItems []*models.ParsedItem, deserializeItem bool) ([]*models.Item, error) {
	items := []*models.Item{}
	for _, parsedItem := range parsedItems {
		item, err := DeserializeItem(svc, parsedItem, deserializeItem)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func DeserializeItem(svc *kms.KMS, parsedItem *models.ParsedItem, deserializeItem bool) (*models.Item, error) {
	value := parsedItem.Value.Value
	if deserializeItem {
		deserializedValue, err := deserialize(svc, parsedItem)
		if err != nil {
			return nil, err
		}
		value = deserializedValue
	}
	return &models.Item{
		Key:           parsedItem.Key,
		Value:         value,
		Serialization: parsedItem.Value.Serialization.Type,
	}, nil
}

func serialize(svc *kms.KMS, value *models.ParsedItemValue) (string, error) {
	switch value.Serialization.Type {
	case "plain":
		return value.Value, nil
	case "base64":
		return encodeBase64([]byte(value.Value)), nil
	case "kms":
		params := &kms.EncryptInput{
			KeyId:     aws.String(value.Serialization.Options["key"]),
			Plaintext: []byte(value.Value),
		}
		resp, err := svc.Encrypt(params)
		if err != nil {
			return "", err
		}
		return encodeBase64(resp.CiphertextBlob), nil
	default:
		return "", fmt.Errorf("Unknown serialization type %s", value.Serialization.Type)
	}

}

func deserialize(svc *kms.KMS, item *models.ParsedItem) (string, error) {
	switch item.Value.Serialization.Type {
	case "plain":
		return item.Value.Value, nil
	case "base64":
		decoded, err := decodeBase64(item.Value.Value)
		if err != nil {
			return "", err
		}
		return string(decoded), nil
	case "kms":
		decoded, err := decodeBase64(item.Value.Value)
		if err != nil {
			return "", err
		}
		params := &kms.DecryptInput{
			CiphertextBlob: decoded,
		}
		resp, err := svc.Decrypt(params)
		if err != nil {
			return "", err
		}
		return string(resp.Plaintext), nil
	default:
		return "", fmt.Errorf("Unknown serialization type %s", item.Value.Serialization.Type)

	}
}

func encodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func decodeBase64(data string) ([]byte, error) {
	res, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return res, nil
}
