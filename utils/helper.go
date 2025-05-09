package utils

import "encoding/json"

func Jparse(str struct{}) ([]byte, error) {
	bs, err := json.Marshal(str)
	if err != nil {
		return nil, err
	}
	return bs, nil
}
