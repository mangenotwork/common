package utils

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// Charset 字符集类型
type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
	GBK     = Charset("GBK")
	GB2312  = Charset("GB2312")
)

// ConvertByte2String 编码转换
func ConvertByte2String(byte []byte, charset Charset) string {
	var str string
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)

	case GBK:
		var decodeBytes, _ = simplifiedchinese.GBK.NewDecoder().Bytes(byte)
		str = string(decodeBytes)

	case GB2312:
		var decodeBytes, _ = simplifiedchinese.HZGB2312.NewDecoder().Bytes(byte)
		str = string(decodeBytes)

	case UTF8:
		fallthrough

	default:
		str = string(byte)
	}

	return str
}

func UnicodeDec(raw string) string {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(raw), `\\u`, `\u`, -1))
	if err != nil {
		return ""
	}
	return str
}

func UnicodeDecByte(raw []byte) []byte {
	rawStr := string(raw)
	return []byte(UnicodeDec(rawStr))
}

// UnescapeUnicode Unicode 转码
func UnescapeUnicode(raw []byte) ([]byte, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

// Base64Encode base64 编码
func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// Base64Decode base64 解码
func Base64Decode(str string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(str)
	return string(b), err
}

// Base64UrlEncode base64 url 编码
func Base64UrlEncode(str string) string {
	return base64.URLEncoding.EncodeToString([]byte(str))
}

// Base64UrlDecode base64 url 解码
func Base64UrlDecode(str string) (string, error) {
	b, err := base64.URLEncoding.DecodeString(str)
	return string(b), err
}

// convert
func convert(dstCharset string, srcCharset string, src string) (dst string, err error) {
	if dstCharset == srcCharset {
		return src, nil
	}
	dst = src
	// Converting `src` to UTF-8.
	if srcCharset != "UTF-8" {
		if e := getEncoding(srcCharset); e != nil {
			tmp, err := ioutil.ReadAll(
				transform.NewReader(bytes.NewReader([]byte(src)), e.NewDecoder()),
			)
			if err != nil {
				return "", err
			}
			src = string(tmp)
		} else {
			return dst, err
		}
	}
	// Do the converting from UTF-8 to `dstCharset`.
	if dstCharset != "UTF-8" {
		if e := getEncoding(dstCharset); e != nil {
			tmp, err := ioutil.ReadAll(
				transform.NewReader(bytes.NewReader([]byte(src)), e.NewEncoder()),
			)
			if err != nil {
				return "", err
			}
			dst = string(tmp)
		} else {
			return dst, err
		}
	} else {
		dst = src
	}
	return dst, nil
}

var (
	// Alias for charsets.
	charsetAlias = map[string]string{
		"HZGB2312": "HZ-GB-2312",
		"hzgb2312": "HZ-GB-2312",
		"GB2312":   "HZ-GB-2312",
		"gb2312":   "HZ-GB-2312",
	}
)

func getEncoding(charset string) encoding.Encoding {
	if c, ok := charsetAlias[charset]; ok {
		charset = c
	}
	enc, err := ianaindex.MIB.Encoding(charset)
	if err != nil {
		log.Println(err)
	}
	return enc
}

// ToUTF8  to utf8
func ToUTF8(srcCharset string, src string) (dst string, err error) {
	return convert("UTF-8", srcCharset, src)
}

// UTF8To utf8 to
func UTF8To(dstCharset string, src string) (dst string, err error) {
	return convert(dstCharset, "UTF-8", src)
}

// ToUTF16 to utf16
func ToUTF16(srcCharset string, src string) (dst string, err error) {
	return convert("UTF-16", srcCharset, src)
}

// UTF16To utf16 to
func UTF16To(dstCharset string, src string) (dst string, err error) {
	return convert(dstCharset, "UTF-16", src)
}

// ToBIG5  to big5
func ToBIG5(srcCharset string, src string) (dst string, err error) {
	return convert("big5", srcCharset, src)
}

// BIG5To  big to
func BIG5To(dstCharset string, src string) (dst string, err error) {
	return convert(dstCharset, "big5", src)
}

// ToGDK to gdk
func ToGDK(srcCharset string, src string) (dst string, err error) {
	return convert("gbk", srcCharset, src)
}

// GDKTo gdk to
func GDKTo(dstCharset string, src string) (dst string, err error) {
	return convert(dstCharset, "gbk", src)
}

// ToGB18030  to gb18030
func ToGB18030(srcCharset string, src string) (dst string, err error) {
	return convert("gb18030", srcCharset, src)
}

// GB18030To gb18030 to
func GB18030To(dstCharset string, src string) (dst string, err error) {
	return convert(dstCharset, "gb18030", src)
}

// ToGB2312 to gb2312
func ToGB2312(srcCharset string, src string) (dst string, err error) {
	return convert("GB2312", srcCharset, src)
}

// GB2312To gb2312 to
func GB2312To(dstCharset string, src string) (dst string, err error) {
	return convert(dstCharset, "GB2312", src)
}

// ToHZGB2312 to hzgb2312
func ToHZGB2312(srcCharset string, src string) (dst string, err error) {
	return convert("HZGB2312", srcCharset, src)
}

// HZGB2312To hzgb2312 to
func HZGB2312To(dstCharset string, src string) (dst string, err error) {
	return convert(dstCharset, "HZGB2312", src)
}

// ConvertStrToGBK 将utf-8编码的字符串转换为GBK编码
func ConvertStrToGBK(str string) string {
	ret, err := simplifiedchinese.GBK.NewEncoder().String(str)
	if err != nil {
		ret = str
	}
	return ret
}

// ConvertGBKToStr 将GBK编码的字符串转换为utf-8编码
func ConvertGBKToStr(gbkStr string) string {
	ret, err := simplifiedchinese.GBK.NewDecoder().String(gbkStr)
	if err != nil {
		ret = gbkStr
	}
	return ret
}
