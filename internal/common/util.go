package common

import "encoding/json"

func MapToStruct(src interface{}, dest interface{}) error {
	tmp, err := json.Marshal(src)
	if err != nil {
		return err
	}

	err = json.Unmarshal(tmp, dest)
	if err != nil {
		return err
	}

	return nil
}
