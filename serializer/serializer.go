package serializer

import (
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/diasjorge/dynamokv/models"
	"github.com/diasjorge/dynamokv/parser"
)

func SerializeItems(svc *kms.KMS, items []*parser.Item) ([]*models.Item, error) {
	result := []*models.Item{}
	for _, item := range items {
		value, err := serialize(svc, item.Value)
		if err != nil {
			return nil, err
		}
		result = append(result, &models.Item{
			Key:           item.Key,
			Value:         value,
			Serialization: item.Value.Serialization.Type,
		})
	}
	return result, nil
}

func DeserializeItems(svc *kms.KMS, items []*models.Item, deserilizeItems bool) ([]*models.Item, error) {
	if deserilizeItems {
		for _, item := range items {
			value, err := deserialize(svc, item)
			if err != nil {
				return nil, err
			}
			item.Value = value
		}
	}
	return items, nil
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

func serialize(svc *kms.KMS, value *parser.ItemValue) (string, error) {
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

func deserialize(svc *kms.KMS, item *models.Item) (string, error) {
	switch item.Serialization {
	case "plain":
		return item.Value, nil
	case "base64":
		decoded, err := decodeBase64(item.Value)
		if err != nil {
			return "", err
		}
		return string(decoded), nil
	case "kms":
		decoded, err := decodeBase64(item.Value)
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
		return "", fmt.Errorf("Unknown serialization type %s", item.Serialization)

	}
}
