package utils

import (
	"github.com/mangenotwork/common/log"
	"testing"
)

// go test utils_test.go charset.go -test.run Test_charset -v
// go test utils_test.go charset.go -cover -v -count=1 -timeout 15s
// go test -benchmem -bench=Test_charset charset.go utils_test.go
func Test_charset(t *testing.T) {
	a := Base64Encode("aaaaaa")
	log.Print(a)
	log.Print(Base64Decode(a))
	b, _ := UTF8To("GB2312", "阿萨声卡十分罕见卡所发海军")
	log.Print(b)
	log.Print(ToUTF8("GB2312", b))
	c := ConvertStrToGBK("阿萨声卡十分罕见卡所发海军")
	log.Print(c)
	log.Print(ConvertGBKToStr(c))

}

// go test -benchmem -bench=charset
func Benchmark_charset(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Base64Encode("aaaaaa")
	}
}

// go test -test.run Test_crypto -v
func Test_crypto(t *testing.T) {
	aes := NewAES(CBC)
	str := "aaaaaaaaaaaa"
	key := "0123456789123456" // key 必须16位
	s, _ := aes.Encrypt([]byte(str), []byte(key))
	log.Print(string(s))
	s1, _ := aes.Decrypt(s, []byte(key))
	log.Print(string(s1))
}

// go test -test.run Test_html_extract -v
func Test_html_extract(t *testing.T) {
	htmlStr := `
<div class="p-20 userbg ub"><img src="https://img003.33633.cn/avatar/20221128/f2c7d78362333c14.jpg" class="avatar50 mr-10"> <div class="ub-f1"><div class="text-white mb-10 text-16">eeeexxxx</div> <div class="text-white mb-10 text-16"><span class="opacity-5">粉丝 </span><span>0</span><span class="opacity-5"> 人</span></div> <div class="text-white text-16"><span class="opacity-5">中奖</span> <span>0</span> <span class="opacity-5">次，共计</span> <span>0</span><span class="opacity-5">元</span></div></div></div>
`
	log.Print(GetPointHTML(htmlStr, "img", "", ""))
}

// go test -test.run Test_id_worker -v
func Test_id_worker(t *testing.T) {
	log.Print(ID64())
	log.Print(ID())
	log.Print(IDStr())
	log.Print(IDMd5())
}

// TODO json_find test
// TODO map test
// TODO reg test
// TODO set test
// TODO slice test
// TODO stack test
// TODO str test
// TODO time test
// TODO type_conversion test
