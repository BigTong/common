package terminator

import (
	"encoding/json"

	"common/parser"
)

func MapToInputConfig(val interface{}) (*parser.InputConfig, error) {
	data, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}

	ret := &parser.InputConfig{}
	err = json.Unmarshal(data, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func MapToElemConfig(val interface{}) (*parser.ElemConfig, error) {
	data, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}

	ret := &parser.ElemConfig{}
	err = json.Unmarshal(data, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func InterfaceToMap(val interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]interface{})
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
