package util

import (
	"context"
	"encoding/json"

	"github.com/lockp111/go-easyzap"
)

func ParseToString(data interface{}) (string, error) {
	dataBytes, err := ConvertToBytes(data)
	if err == nil {
		return string(dataBytes), nil
	}
	return "", err
}

func ConvertToBytes(data interface{}) ([]byte, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		easyzap.Error(context.TODO(), err, "Error while converting data to bytes")
	}
	return dataBytes, err
}
