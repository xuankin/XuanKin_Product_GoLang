package service

import (
	"encoding/json"
	"gorm.io/datatypes"
)

func toMap(data datatypes.JSON) map[string]interface{} {
	var res map[string]interface{}
	if len(data) > 0 {
		json.Unmarshal(data, &res)
	}
	return res
}
func toJson(data interface{}) datatypes.JSON {
	bytes, _ := json.Marshal(data)
	return datatypes.JSON(bytes)
}
