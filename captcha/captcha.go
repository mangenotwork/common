package captcha

import (
	"github.com/gin-gonic/gin"
	"github.com/mangenotwork/common/ginHelper"
	"github.com/mangenotwork/common/utils"
	"github.com/mojocn/base64Captcha"
)

var store = base64Captcha.DefaultMemStore

type OutData struct {
	ID     string `json:"id"`
	Base64 string `json:"base64"`
}

// NewDriver 验证码
func NewDriver() *base64Captcha.DriverString {
	driver := new(base64Captcha.DriverString)
	driver.Height = 44
	driver.Width = 120
	driver.NoiseCount = 5
	driver.ShowLineOptions = base64Captcha.OptionShowSineLine | base64Captcha.OptionShowSlimeLine | base64Captcha.OptionShowHollowLine
	driver.Length = 4
	driver.Source = "1234567890qwertyuipkjhgfdsazxcvbnm"
	driver.Fonts = []string{"wqy-microhei.ttc"}
	return driver
}

func Captcha(c *gin.Context) {
	var driver = NewDriver().ConvertFonts()
	capt := base64Captcha.NewCaptcha(driver, store)
	_, content, answer := capt.Driver.GenerateIdQuestionAnswer()
	id := utils.IDStr()
	item, _ := capt.Driver.DrawCaptcha(content)
	_ = capt.Store.Set(id, answer)
	img := item.EncodeB64string()
	ginHelper.APIOutPut(c, 0, "", &OutData{
		ID:     id,
		Base64: img,
	})
}

func Verify(id, code string) bool {
	if !store.Verify(id, code, true) {
		return false
	}
	return true
}
