package dslutil

import "encoding/json"

//除了用户自身谁也不知道参数的类型是int,interface,还是struct，所以没法做到用户无感。
func GetParam(input *map[string][]byte, name string, v interface{}) error {
	return json.Unmarshal((*input)[name], v)
}

func SetParam(input *map[string][]byte, name string, v interface{}) error {
	valueByte, err := json.Marshal(v)
	if err != nil {
		return err
	}
	(*input)[name] = valueByte
	return err
}
