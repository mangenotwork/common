package qrcode

import (
	"encoding/base64"
	"github.com/skip2/go-qrcode"
)

func GetQRCodeIO(content string) string {
	var png []byte
	png, err := qrcode.Encode(content, qrcode.High, 260)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(png)
}
