package utils

import (
	"encoding/json"
	"strings"
)

func GenerateCacheKey(prefix string, identifiers ...string) string {
	allParts := append([]string{prefix}, identifiers...)
	return strings.Join(allParts, ":")
}

func ConvertRedisDataType[T any](result []byte) (T, error) {
	var data T
	if err := json.Unmarshal(result, &data); err != nil {
		var zero T
		return zero, err
	}

	return data, nil
}
