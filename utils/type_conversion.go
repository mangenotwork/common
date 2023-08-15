package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/mangenotwork/common/log"

	"golang.org/x/text/encoding/simplifiedchinese"
)

// 类型转换

// AnyToString any -> string
func AnyToString(i interface{}) string {
	if i == nil {
		return ""
	}
	if reflect.ValueOf(i).Kind() == reflect.String {
		return i.(string)
	}
	var buf bytes.Buffer
	stringValue(reflect.ValueOf(i), 0, &buf)
	return buf.String()
}

func stringValue(v reflect.Value, indent int, buf *bytes.Buffer) {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Struct:
		buf.WriteString("{\n")
		for i := 0; i < v.Type().NumField(); i++ {
			ft := v.Type().Field(i)
			fv := v.Field(i)
			if ft.Name[0:1] == strings.ToLower(ft.Name[0:1]) {
				continue
			}
			if (fv.Kind() == reflect.Ptr || fv.Kind() == reflect.Slice) && fv.IsNil() {
				continue
			}
			buf.WriteString(strings.Repeat(" ", indent+2))
			buf.WriteString(ft.Name + ": ")
			if tag := ft.Tag.Get("sensitive"); tag == "true" {
				buf.WriteString("<sensitive>")
			} else {
				stringValue(fv, indent+2, buf)
			}
			buf.WriteString(",\n")
		}
		buf.WriteString("\n" + strings.Repeat(" ", indent) + "}")

	case reflect.Slice:
		nl, id, id2 := "", "", ""
		if v.Len() > 3 {
			nl, id, id2 = "\n", strings.Repeat(" ", indent), strings.Repeat(" ", indent+2)
		}
		buf.WriteString("[" + nl)
		for i := 0; i < v.Len(); i++ {
			buf.WriteString(id2)
			stringValue(v.Index(i), indent+2, buf)

			if i < v.Len()-1 {
				buf.WriteString("," + nl)
			}
		}
		buf.WriteString(nl + id + "]")

	case reflect.Map:
		buf.WriteString("{\n")
		for i, k := range v.MapKeys() {
			buf.WriteString(strings.Repeat(" ", indent+2))
			buf.WriteString(k.String() + ": ")
			stringValue(v.MapIndex(k), indent+2, buf)

			if i < v.Len()-1 {
				buf.WriteString(",\n")
			}
		}
		buf.WriteString("\n" + strings.Repeat(" ", indent) + "}")

	default:
		format := "%v"
		switch v.Interface().(type) {
		case string:
			format = "%q"
		}
		_, _ = fmt.Fprintf(buf, format, v.Interface())
	}
}

// JsonToMap json -> map
func JsonToMap(str string) (map[string]interface{}, error) {
	var tempMap map[string]interface{}
	err := json.Unmarshal([]byte(str), &tempMap)
	if err != nil {
		return nil, err
	}
	return tempMap, nil
}

// MapToJson map -> json
func MapToJson(m interface{}) (string, error) {
	jsonStr, err := json.Marshal(m)
	return string(jsonStr), err
}

// AnyToMap interface{} -> map[string]interface{}
func AnyToMap(data interface{}) map[string]interface{} {
	if v, ok := data.(map[string]interface{}); ok {
		return v
	}
	if reflect.ValueOf(data).Kind() == reflect.String {
		dataMap, err := JsonToMap(data.(string))
		if err == nil {
			return dataMap
		}
	}
	return nil
}

// AnyToInt interface{} -> int
func AnyToInt(data interface{}) int {
	var t2 int
	switch data.(type) {
	case uint:
		t2 = int(data.(uint))
		break
	case int8:
		t2 = int(data.(int8))
		break
	case uint8:
		t2 = int(data.(uint8))
		break
	case int16:
		t2 = int(data.(int16))
		break
	case uint16:
		t2 = int(data.(uint16))
		break
	case int32:
		t2 = int(data.(int32))
		break
	case uint32:
		t2 = int(data.(uint32))
		break
	case int64:
		t2 = int(data.(int64))
		break
	case uint64:
		t2 = int(data.(uint64))
		break
	case float32:
		t2 = int(data.(float32))
		break
	case float64:
		t2 = int(data.(float64))
		break
	case string:
		t2, _ = strconv.Atoi(data.(string))
		break
	default:
		t2 = data.(int)
		break
	}
	return t2
}

// AnyToInt64 interface{} -> int64
func AnyToInt64(data interface{}) int64 {
	return int64(AnyToInt(data))
}

// AnyToArr interface{} -> []interface{}
func AnyToArr(data interface{}) []interface{} {
	if v, ok := data.([]interface{}); ok {
		return v
	}
	return nil
}

// AnyToFloat64 interface{} -> float64
func AnyToFloat64(data interface{}) float64 {
	if v, ok := data.(float64); ok {
		return v
	}
	if v, ok := data.(float32); ok {
		return float64(v)
	}
	return 0
}

// AnyToStrings interface{} -> []string
func AnyToStrings(data interface{}) []string {
	listValue, ok := data.([]interface{})
	if !ok {
		return nil
	}
	keyStringValues := make([]string, len(listValue))
	for i, arg := range listValue {
		keyStringValues[i] = arg.(string)
	}
	return keyStringValues
}

// AnyToJson interface{} -> json string
func AnyToJson(data interface{}) (string, error) {
	jsonStr, err := json.Marshal(data)
	return string(jsonStr), err
}

// AnyToJsonB interface{} -> json string
func AnyToJsonB(data interface{}) ([]byte, error) {
	jsonStr, err := json.Marshal(data)
	return jsonStr, err
}

// IntToHex int -> hex
func IntToHex(i int) string {
	return fmt.Sprintf("%x", i)
}

// Int64ToHex int64 -> hex
func Int64ToHex(i int64) string {
	return fmt.Sprintf("%x", i)
}

// HexToInt hex -> int
func HexToInt(s string) int {
	n, err := strconv.ParseUint(s, 16, 8)
	if err != nil {
		panic("Parse Error")
	}
	n2 := uint8(n)
	return int(*(*int8)(unsafe.Pointer(&n2)))
}

// HexToInt64 hex -> int
func HexToInt64(s string) int64 {
	n, err := strconv.ParseUint(s, 16, 8)
	if err != nil {
		panic("Parse Error")
	}
	n2 := uint8(n)
	return int64(*(*int8)(unsafe.Pointer(&n2)))
}

// StrNumToInt64 string -> int64
func StrNumToInt64(str string) int64 {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

// StrNumToInt string -> int
func StrNumToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}

// StrNumToInt32 string -> int32
func StrNumToInt32(str string) int32 {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return int32(i)
}

// StrNumToFloat64 string -> float64
func StrNumToFloat64(str string) float64 {
	i, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return i
}

// StrNumToFloat32 string -> float32
func StrNumToFloat32(str string) float32 {
	i, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return 0
	}
	return float32(i)
}

// Uint8ToStr []uint8 -> string
func Uint8ToStr(bs []uint8) string {
	ba := make([]byte, 0)
	for _, b := range bs {
		ba = append(ba, b)
	}
	return string(ba)
}

// StrToByte string -> []byte
func StrToByte(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// ByteToStr []byte -> string
func ByteToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// BoolToByte bool -> []byte
func BoolToByte(b bool) []byte {
	if b == true {
		return []byte{1}
	}
	return []byte{0}
}

// ByteToBool []byte -> bool
func ByteToBool(b []byte) bool {
	if len(b) == 0 || bytes.Compare(b, make([]byte, len(b))) == 0 {
		return false
	}
	return true
}

// IntToByte int -> []byte
func IntToByte(i int) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(i))
	return b
}

// ByteToInt []byte -> int
func ByteToInt(b []byte) int {
	return int(binary.LittleEndian.Uint32(b))
}

// Int64ToByte int64 -> []byte
func Int64ToByte(i int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	return b
}

// ByteToInt64 []byte -> int64
func ByteToInt64(b []byte) int64 {
	return int64(binary.LittleEndian.Uint64(b))
}

// Float32ToByte float32 -> []byte
func Float32ToByte(f float32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, Float32ToUint32(f))
	return b
}

// Float32ToUint32 float32 -> uint32
func Float32ToUint32(f float32) uint32 {
	return math.Float32bits(f)
}

// ByteToFloat32 []byte -> float32
func ByteToFloat32(b []byte) float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(b))
}

// Float64ToByte float64 -> []byte
func Float64ToByte(f float64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, Float64ToUint64(f))
	return b
}

// Float64ToUint64 float64 -> uint64
func Float64ToUint64(f float64) uint64 {
	return math.Float64bits(f)
}

// ByteToFloat64 []byte -> float64
func ByteToFloat64(b []byte) float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(b))
}

// StructToMap  struct -> map[string]interface{}
func StructToMap(obj interface{}) map[string]interface{} {
	rt, rv := reflect.TypeOf(obj), reflect.ValueOf(obj)
	if rt != nil && rt.Kind() != reflect.Struct {
		return make(map[string]interface{}, 0)
	}
	out := make(map[string]interface{}, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		// Unexported fields, access not allowed
		if field.PkgPath != "" {
			continue
		}
		var fieldName string
		if tagVal, ok := field.Tag.Lookup("json"); ok {
			// Honor the special "-" in json attribute
			if strings.HasPrefix(tagVal, "-") {
				continue
			}
			fieldName = tagVal
		} else {
			fieldName = field.Name
		}
		val := valueToInterface(rv.Field(i))
		if val != nil {
			out[fieldName] = val
		}
	}
	return out
}

func valueToInterface(value reflect.Value) interface{} {
	if !value.IsValid() {
		return nil
	}
	switch value.Type().Kind() {
	case reflect.Struct:
		return StructToMap(value.Interface())

	case reflect.Ptr:
		if !value.IsNil() {
			return valueToInterface(value.Elem())
		}

	case reflect.Array:
	case reflect.Slice:
		arr := make([]interface{}, 0, value.Len())
		for i := 0; i < value.Len(); i++ {
			val := valueToInterface(value.Index(i))
			if val != nil {
				arr = append(arr, val)
			}
		}
		return arr

	case reflect.Map:
		m := make(map[string]interface{}, value.Len())
		for _, k := range value.MapKeys() {
			v := value.MapIndex(k)
			m[k.String()] = valueToInterface(v)
		}
		return m

	default:
		return value.Interface()
	}

	return nil
}

// EncodeByte encode byte
func EncodeByte(v interface{}) []byte {
	switch value := v.(type) {
	case int, int8, int16, int32:
		return IntToByte(value.(int))
	case int64:
		return Int64ToByte(value)
	case string:
		return StrToByte(value)
	case bool:
		return BoolToByte(value)
	case float32:
		return Float32ToByte(value)
	case float64:
		return Float64ToByte(value)
	}
	return []byte("")
}

// DecodeByte  decode byte
func DecodeByte(b []byte) (interface{}, error) {
	var values interface{}
	buf := bytes.NewBuffer(b)
	err := binary.Read(buf, binary.BigEndian, values)
	return values, err
}

// ByteToBit []byte -> []uint8 (bit)
func ByteToBit(b []byte) []uint8 {
	bits := make([]uint8, 0)
	for _, v := range b {
		bits = bits2Uint(bits, uint(v), 8)
	}
	return bits
}

// bits2Uint bits2Uint
func bits2Uint(bits []uint8, ui uint, l int) []uint8 {
	a := make([]uint8, l)
	for i := l - 1; i >= 0; i-- {
		a[i] = uint8(ui & 1)
		ui >>= 1
	}
	if bits != nil {
		return append(bits, a...)
	}
	return a
}

// BitToByte []uint8 -> []byte
func BitToByte(b []uint8) []byte {
	if len(b)%8 != 0 {
		for i := 0; i < len(b)%8; i++ {
			b = append(b, 0)
		}
	}
	by := make([]byte, 0)
	for i := 0; i < len(b); i += 8 {
		by = append(b, byte(bitsToUint(b[i:i+8])))
	}
	return by
}

// bitsToUint bitsToUint
func bitsToUint(bits []uint8) uint {
	v := uint(0)
	for _, i := range bits {
		v = v<<1 | uint(i)
	}
	return v
}

// StructToMapV2 Struct  ->  map
// hasValue=true表示字段值不管是否存在都转换成map
// hasValue=false表示字段为不为空或者不为0则转换成map
func StructToMapV2(obj interface{}, hasValue bool) (map[string]interface{}, error) {
	mp := make(map[string]interface{})
	value := reflect.ValueOf(obj).Elem()
	typeOf := reflect.TypeOf(obj).Elem()
	for i := 0; i < value.NumField(); i++ {
		switch value.Field(i).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if hasValue {
				if value.Field(i).Int() != 0 {
					mp[typeOf.Field(i).Name] = value.Field(i).Int()
				}
			} else {
				mp[typeOf.Field(i).Name] = value.Field(i).Int()
			}

		case reflect.String:
			if hasValue {
				if len(value.Field(i).String()) != 0 {
					mp[typeOf.Field(i).Name] = value.Field(i).String()
				}
			} else {
				mp[typeOf.Field(i).Name] = value.Field(i).String()
			}

		case reflect.Float32, reflect.Float64:
			if hasValue {
				if len(value.Field(i).String()) != 0 {
					mp[typeOf.Field(i).Name] = value.Field(i).Float()
				}
			} else {
				mp[typeOf.Field(i).Name] = value.Field(i).Float()
			}

		case reflect.Bool:
			if hasValue {
				if len(value.Field(i).String()) != 0 {
					mp[typeOf.Field(i).Name] = value.Field(i).Bool()
				}
			} else {
				mp[typeOf.Field(i).Name] = value.Field(i).Bool()
			}

		default:
			return mp, fmt.Errorf("数据类型不匹配")
		}
	}
	return mp, nil
}

// StructToMapV3 struct -> map
func StructToMapV3(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)
	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		data[obj1.Field(i).Name] = obj2.Field(i).Interface()
	}
	return data
}

// PanicToError panic -> error
func PanicToError(fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Panic error: %v", r)
		}
	}()
	fn()
	return
}

// P2E panic -> error
func P2E() {
	defer func() {
		if r := recover(); r != nil {
			log.Error("Panic error: ", r)
		}
	}()
}

// ByteToBinaryString  字节 -> 二进制字符串
func ByteToBinaryString(data byte) (str string) {
	var a byte
	for i := 0; i < 8; i++ {
		a = data
		data <<= 1
		data >>= 1
		switch a {
		case data:
			str += "0"
		default:
			str += "1"
		}
		data <<= 1
	}
	return str
}

// MapStrToAny map[string]string -> map[string]interface{}
func MapStrToAny(m map[string]string) map[string]interface{} {
	dest := make(map[string]interface{})
	for k, v := range m {
		dest[k] = interface{}(v)
	}
	return dest
}

// ByteToGBK   byte -> gbk byte
func ByteToGBK(strBuf []byte) []byte {
	if IsUtf8(strBuf) {
		if GBKBuf, err := simplifiedchinese.GBK.NewEncoder().Bytes(strBuf); err == nil {
			if IsUtf8(GBKBuf) {
				return GBKBuf
			}
		}
		if GB18030Buf, err := simplifiedchinese.GB18030.NewEncoder().Bytes(strBuf); err == nil {
			if IsUtf8(GB18030Buf) {
				return GB18030Buf
			}
		}
		if HZGB2312Buf, err := simplifiedchinese.HZGB2312.NewEncoder().Bytes(strBuf); err == nil {
			if IsUtf8(HZGB2312Buf) {
				return HZGB2312Buf
			}
		}
		return strBuf
	} else {
		return strBuf
	}
}

// Int64ToStr int64 -> string
func Int64ToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}
