# 常用实用方法


#### 类型转换 AnyToString any -> string
> func AnyToString(i interface{}) string

---

#### 类型转换 JsonToMap json -> map
> func JsonToMap(str string) (map[string]interface{}, error)

---

#### 类型转换 MapToJson map -> json
> func MapToJson(m interface{}) (string, error)

---

#### 类型转换 AnyToMap interface{} -> map[string]interface{}
> func AnyToMap(data interface{}) map[string]interface{}

---

#### 类型转换 AnyToInt interface{} -> int
> func AnyToInt(data interface{}) int

---

#### 类型转换 AnyToInt64 interface{} -> int64
> func AnyToInt64(data interface{}) int64

---

#### 类型转换 AnyToArr interface{} -> []interface{}
> func AnyToArr(data interface{}) []interface{}

---

#### 类型转换 AnyToFloat64 interface{} -> float64
> func AnyToFloat64(data interface{}) float64

---

#### 类型转换 AnyToStrings interface{} -> []string
> func AnyToStrings(data interface{}) []string

---

#### 类型转换 AnyToJson interface{} -> json string
> func AnyToJson(data interface{}) (string, error)

---

#### 类型转换 IntToHex int -> hex
> func IntToHex(i int) string

---

#### 类型转换 Int64ToHex int64 -> hex
> func Int64ToHex(i int64) string

---

#### 类型转换 HexToInt hex -> int
> func HexToInt(s string) int

---

#### 类型转换 HexToInt64 hex -> int
> func HexToInt64(s string) int64

---

#### 类型转换 StrNumToInt64 string -> int64
> func StrNumToInt64(str string) int64

---

#### 类型转换 StrNumToInt string -> int
> func StrNumToInt(str string) int

---

#### 类型转换 StrNumToInt32 string -> int32
> func StrNumToInt32(str string) int32

---

#### 类型转换 StrNumToFloat64 string -> float64
> func StrNumToFloat64(str string) float64

---

#### 类型转换 StrNumToFloat32 string -> float32
> func StrNumToFloat32(str string) float32

---

#### 类型转换 Uint8ToStr []uint8 -> string
> func Uint8ToStr(bs []uint8) string

---

#### 类型转换 StrToByte string -> []byte
> func StrToByte(s string) []byte

---

#### 类型转换 ByteToStr []byte -> string
> func ByteToStr(b []byte) string

---

#### 类型转换 BoolToByte bool -> []byte
> func BoolToByte(b bool) []byte

---

#### 类型转换 ByteToBool []byte -> bool
> func ByteToBool(b []byte) bool

---

#### 类型转换 IntToByte int -> []byte
> func IntToByte(i int) []byte

---

#### 类型转换 ByteToInt []byte -> int
> func ByteToInt(b []byte) int

---

#### 类型转换 func Int64ToByte(i int64) []byte
> func Int64ToByte(i int64) []byte

---

#### 类型转换 ByteToInt64 []byte -> int64
> func ByteToInt64(b []byte) int64

---

#### 类型转换 Float32ToByte float32 -> []byte
> func Float32ToByte(f float32) []byte

---

#### 类型转换 Float32ToUint32 float32 -> uint32
> func Float32ToUint32(f float32) uint32

---

#### 类型转换 ByteToFloat32 []byte -> float32
> func ByteToFloat32(b []byte) float32

---

#### 类型转换 Float64ToByte float64 -> []byte
> func Float64ToByte(f float64) []byte

---

#### 类型转换 Float64ToUint64 float64 -> uint64
> func Float64ToUint64(f float64) uint64

---

#### 类型转换 ByteToFloat64 []byte -> float64
> func ByteToFloat64(b []byte) float64

---

#### 类型转换 StructToMap  struct -> map[string]interface{}
> func StructToMap(obj interface{}) map[string]interface{}

---

#### EncodeByte encode byte 
> func EncodeByte(v interface{}) []byte

---

#### DecodeByte  decode byte
> func DecodeByte(b []byte) (interface{}, error)

---

#### 类型转换 ByteToBit []byte -> []uint8 (bit)
> func ByteToBit(b []byte) []uint8

---

#### 类型转换 BitToByte []uint8 -> []byte
> func BitToByte(b []uint8) []byte

---

#### 类型转换 StructToMapV2 Struct  ->  map
> func StructToMapV2(obj interface{}, hasValue bool) (map[string]interface{}, error)

参数
```shell
hasValue=true表示字段值不管是否存在都转换成map
hasValue=false表示字段为不为空或者不为0则转换成map
```

---

#### 类型转换 StructToMapV3 struct -> map
> func StructToMapV3(obj interface{}) map[string]interface{}

---

#### 类型转换 PanicToError panic -> error
> func PanicToError(fn func()) (err error)

---

#### 类型转换 P2E panic -> error
> func P2E()

---

#### 类型转换 ByteToBinaryString  字节 -> 二进制字符串
> func ByteToBinaryString(data byte) (str string)

---

#### 类型转换 MapStrToAny map[string]string -> map[string]interface{}
> func MapStrToAny(m map[string]string) map[string]interface{}

---

#### 类型转换 ByteToGBK   byte -> gbk byte
> func ByteToGBK(strBuf []byte) []byte

---

#### 类型转换 Int64ToStr int64 -> string
> func Int64ToStr(i int64) string

---

#### 
>

---

####
>

---

####
>

---

####
>

---

####
>

---

####
>

---

####
>

---

####
>

---

####
>

---

####
>

---

####
>

---

####
>

---

####
>

---

####
>

---

####
>

---

####
>

---

####
>

---

####
>

---

####
>

---

####
>

---

####
>

---

####
>

---

####
>

---

####
>

---

####
>

---

####
>

---

####
>

---

####
>

---

####
>

---

####
>

---

####
>

---

####
>

---

