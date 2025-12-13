package utils

import "encoding/json"

func MustMarshalJSON(v interface{}) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
