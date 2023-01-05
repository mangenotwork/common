package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

// JsonFind 按路径寻找指定json值
// 用法参考  ./_examples/json/main.go
// @find : 寻找路径，与目录的url类似， 下面是一个例子：
// json:  {a:[{b:1},{b:2}]}
// find=/a/[0]  =>   {b:1}
// find=a/[0]/b  =>   1
func JsonFind(jsonStr, find string) (interface{}, error) {
	if !IsJson(jsonStr) {
		return nil, fmt.Errorf("不是标准的Json格式")
	}
	jxList := strings.Split(find, "/")
	jxLen := len(jxList)
	var (
		data  = AnyToMap(jsonStr)
		value interface{}
	)
	for i := 0; i < jxLen; i++ {
		l := len(jxList[i])
		if l > 2 && string(jxList[i][0]) == "[" && string(jxList[i][l-1]) == "]" {
			numStr := jxList[i][1 : l-1]
			dataList := AnyToArr(value)
			value = dataList[AnyToInt(numStr)]
			data = AnyToMap(value)
		} else {
			if IsHaveKey(data, jxList[i]) {
				value = data[jxList[i]]
				data = AnyToMap(value)
			} else {
				value = nil
			}
		}
	}
	return value, nil
}

// JsonFind2Json 寻找json,输出 json格式字符串
func JsonFind2Json(jsonStr, find string) (string, error) {
	value, err := JsonFind(jsonStr, find)
	if err != nil {
		return "", err
	}
	return MapToJson(value)
}

// JsonFind2Map 寻找json,输出 map[string]interface{}
func JsonFind2Map(jsonStr, find string) (map[string]interface{}, error) {
	value, err := JsonFind(jsonStr, find)
	if err != nil {
		return nil, err
	}
	return AnyToMap(value), nil
}

// JsonFind2Arr 寻找json,输出 []interface{}
func JsonFind2Arr(jsonStr, find string) ([]interface{}, error) {
	value, err := JsonFind(jsonStr, find)
	if err != nil {
		return nil, err
	}
	return AnyToArr(value), nil
}

// IsJson 是否是json格式
func IsJson(str string) bool {
	var tempMap map[string]interface{}
	err := json.Unmarshal([]byte(str), &tempMap)
	if err != nil {
		return false
	}
	return true
}

// IsHaveKey map[string]interface{} 是否存在 输入的key
func IsHaveKey(data map[string]interface{}, key string) bool {
	_, ok := data[key]
	return ok
}
