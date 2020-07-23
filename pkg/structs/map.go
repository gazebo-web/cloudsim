package structs

import "encoding/json"

// ToMap converts the given struct into map[string]interface{}.
// NOTE: The struct need to have json labels.
func ToMap(input interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	data, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
