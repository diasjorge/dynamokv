package serializer

import (
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/diasjorge/dynamokv/parser"
)

func SerializeItems(svc *kms.KMS, items []*parser.Item) (map[string]string, error) {
	result := make(map[string]string)
	for _, item := range items {
		value, err := serialize(svc, item.Value)
		if err != nil {
			return nil, err
		}
		result[item.Key] = value
	}
	return result, nil
}

func encodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
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
