package common

import (
	"encoding/json"
)

func InterfaceToObject(data, result interface{}) error {
	dataB, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(dataB, result)
	if err != nil {
		return err
	}

	return nil
}
