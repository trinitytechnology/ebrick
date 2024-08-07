package utils

import (
	"encoding/json"
	"unsafe"
)

func UnmarshalJSONByte[T any](result []byte) (T, error) {
	var data T
	if err := json.Unmarshal(result, &data); err != nil {
		var zero T
		return zero, err
	}

	return data, nil
}

func UnmarshalJSON[T any](result string) (T, error) {
	var data T
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		var zero T
		return zero, err
	}

	return data, nil
}

func MarshalJSON(in any) string {
	bs, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}
	return BinaryString(bs)
}
func MarshalJSONBytes(in any) []byte {
	bs, err := json.Marshal(in)
	if err != nil {
		return []byte{}
	}
	return bs
}

func BinaryString(bs []byte) string {
	return unsafe.String(unsafe.SliceData(bs), len(bs))
}
