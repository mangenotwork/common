package utils

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/mangenotwork/common/log"
)

// 字符串相关

// CleaningStr 清理字符串前后空白 和回车 换行符号
func CleaningStr(str string) string {
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\r", "", -1)
	str = strings.Replace(str, "\\n", "", -1)
	//str = strings.Replace(str, "\"", "", -1)
	str = strings.TrimSpace(str)
	str = StrDeleteSpace(str)
	return str
}

// StrDeleteSpace 删除字符串前后的空格
func StrDeleteSpace(str string) string {
	strList := []byte(str)
	spaceCount, count := 0, len(strList)
	for i := 0; i <= len(strList)-1; i++ {
		if strList[i] == 32 {
			spaceCount++
		} else {
			break
		}
	}
	strList = strList[spaceCount:]
	spaceCount, count = 0, len(strList)
	for i := count - 1; i >= 0; i-- {
		if strList[i] == 32 {
			spaceCount++
		} else {
			break
		}
	}
	return string(strList[:count-spaceCount])
}

// SizeFormat 字节的单位转换 保留两位小数
func SizeFormat(size int64) string {
	if size < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2fB", float64(size)/float64(1))
	} else if size < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(size)/float64(1024))
	} else if size < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(size)/float64(1024*1024))
	} else if size < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(size)/float64(1024*1024*1024))
	} else if size < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(size)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(size)/float64(1024*1024*1024*1024*1024))
	}
}

// DeepCopy 深copy
func DeepCopy(dst, src interface{}) error {
	return deepCopy(dst, src)
}

func deepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

// IsContainStr  字符串是否等于items中的某个元素
func IsContainStr(items []string, item string) bool {
	for i := 0; i < len(items); i++ {
		if items[i] == item {
			return true
		}
	}
	return false
}

// FileMd5  file md5   文件md5
func FileMd5(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", md5hash.Sum(nil)), nil
}

// PathExists 目录不存在则创建
func PathExists(path string) {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		os.MkdirAll(path, 0777)
	}
}

// StrDuplicates  数组，切片去重和去空串
func StrDuplicates(a []string) []string {
	m := make(map[string]struct{})
	ret := make([]string, 0, len(a))
	for i := 0; i < len(a); i++ {
		if a[i] == "" {
			continue
		}
		if _, ok := m[a[i]]; !ok {
			m[a[i]] = struct{}{}
			ret = append(ret, a[i])
		}
	}
	return ret
}

// IsElementStr 判断字符串是否与数组里的某个字符串相同
func IsElementStr(listData []string, element string) bool {
	for _, k := range listData {
		if k == element {
			return true
		}
	}
	return false
}

// windowsPath windows平台需要转一下
func windowsPath(path string) string {
	if runtime.GOOS == "windows" {
		path = strings.Replace(path, "\\", "/", -1)
	}
	return path
}

// GetNowPath 获取当前运行路径
func GetNowPath() string {
	path, err := os.Getwd()
	if err != nil {
		log.Error(err)
		return ""
	}
	return windowsPath(path)
}

// FileMd5sum 文件 Md5
func FileMd5sum(fileName string) string {
	fin, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	if err != nil {
		log.Error(fileName, err)
		return ""
	}
	defer func() {
		_ = fin.Close()
	}()
	buf, bufErr := ioutil.ReadFile(fileName)
	if bufErr != nil {
		log.Error(fileName, bufErr)
		return ""
	}
	m := md5.Sum(buf)
	return hex.EncodeToString(m[:16])
}

// SearchBytesIndex []byte 字节切片 循环查找
func SearchBytesIndex(bSrc []byte, b byte) int {
	for i := 0; i < len(bSrc); i++ {
		if bSrc[i] == b {
			return i
		}
	}
	return -1
}

// IF 三元表达式
// use: IF(a>b, a, b).(int)
func IF(condition bool, a, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}

// ReplaceAllToOne 批量统一替换字符串
func ReplaceAllToOne(str string, from []string, to string) string {
	arr := make([]string, len(from)*2)
	for i, s := range from {
		arr[i*2] = s
		arr[i*2+1] = to
	}
	r := strings.NewReplacer(arr...)
	return r.Replace(str)
}

// Exists 路径是否存在
func Exists(path string) bool {
	if stat, err := os.Stat(path); stat != nil && !os.IsNotExist(err) {
		return true
	}
	return false
}

// IsDir 是否是目录
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile 是否是文件
func IsFile(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !s.IsDir()
}

// GzipCompress gzip压缩
func GzipCompress(src []byte) []byte {
	var in bytes.Buffer
	w := gzip.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

// GzipDecompress gzip解压
func GzipDecompress(src []byte) []byte {
	dst := make([]byte, 0)
	br := bytes.NewReader(src)
	gr, err := gzip.NewReader(br)
	if err != nil {
		return dst
	}
	defer gr.Close()
	tmp, err := ioutil.ReadAll(gr)
	if err != nil {
		return dst
	}
	dst = tmp
	return dst
}

// AbPathByCaller 获取当前执行文件绝对路径（go run）
func AbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return path.Join(abPath, "../../../")
}

// GetWD 获取当前工作目录
func GetWD() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return wd
}

// IsUtf8 是否是utf8编码
func IsUtf8(buf []byte) bool {
	return utf8.Valid(buf)
}

// Get16MD5Encode 返回一个16位md5加密后的字符串
func Get16MD5Encode(data string) string {
	return GetMD5Encode(data)[8:24]
}

// GetMD5Encode 获取Md5编码
func GetMD5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// GetAllFile 获取目录下的所有文件
func GetAllFile(pathname string) ([]string, error) {
	s := make([]string, 0)
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		log.Error("read dir fail:", err)
		return s, err
	}
	for _, fi := range rd {
		if !fi.IsDir() {
			fullName := pathname + "/" + fi.Name()
			s = append(s, fullName)
		}
	}
	return s, nil
}

func subString(str string, start, end int) string {
	rs := []rune(str)
	length := len(rs)
	if start < 0 || start > length {
		log.Error("start is wrong")
		return ""
	}
	if end < start || end > length {
		log.Error("end is wrong")
		return ""
	}
	return string(rs[start:end])
}

// DeCompressZIP zip解压文件
func DeCompressZIP(zipFile, dest string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer func() {
		_ = reader.Close()
	}()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		filename := dest + file.Name
		err = os.MkdirAll(subString(filename, 0, strings.LastIndex(filename, "/")), 0755)
		if err != nil {
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			return err
		}
		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
		_ = w.Close()
		_ = rc.Close()
	}
	return nil
}

// DeCompressTAR tar 解压文件
func DeCompressTAR(tarFile, dest string) error {
	file, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()
	tr := tar.NewReader(file)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		filename := dest + hdr.Name
		err = os.MkdirAll(subString(filename, 0, strings.LastIndex(filename, "/")), 0755)
		if err != nil {
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			return err
		}
		_, err = io.Copy(w, tr)
		if err != nil {
			return err
		}
		_ = w.Close()
	}
	return nil
}

// DecompressionZipFile zip压缩文件
func DecompressionZipFile(src, dest string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		_ = reader.Close()
	}()
	for _, file := range reader.File {
		filePath := path.Join(dest, file.Name)
		if file.FileInfo().IsDir() {
			_ = os.MkdirAll(filePath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
				return err
			}
			inFile, err := file.Open()
			if err != nil {
				return err
			}
			outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return err
			}
			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
			_ = inFile.Close()
			_ = outFile.Close()
		}
	}
	return nil
}

// CompressFiles 压缩很多文件
// files 文件数组，可以是不同dir下的文件或者文件夹
// dest 压缩文件存放地址
func CompressFiles(files []string, dest string) error {
	d, _ := os.Create(dest)
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()
	for _, file := range files {
		err := compressFiles(file, "", w)
		if err != nil {
			return err
		}
	}
	return nil
}

func compressFiles(filePath string, prefix string, zw *zip.Writer) error {
	file, err := os.Open(filePath)
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f := file.Name() + "/" + fi.Name()
			if err != nil {
				return err
			}
			err = compressFiles(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// CompressDirZip 压缩目录
func CompressDirZip(src, outFile string) error {
	// 预防：旧文件无法覆盖
	_ = os.RemoveAll(outFile)
	// 创建：zip文件
	zipFile, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer zipFile.Close()
	// 打开：zip文件
	archive := zip.NewWriter(zipFile)
	defer archive.Close()
	// 遍历路径信息
	filepath.Walk(src, func(path string, info os.FileInfo, _ error) error {
		// 如果是源路径，提前进行下一个遍历
		if path == src {
			return nil
		}
		// 获取：文件头信息
		header, _ := zip.FileInfoHeader(info)
		header.Name = strings.TrimPrefix(path, src+`/`)
		// 判断：文件是不是文件夹
		if info.IsDir() {
			header.Name += `/`
		} else {
			// 设置：zip的文件压缩算法
			header.Method = zip.Deflate
		}
		// 创建：压缩包头部信息
		writer, _ := archive.CreateHeader(header)
		if !info.IsDir() {
			file, _ := os.Open(path)
			defer file.Close()
			io.Copy(writer, file)
		}
		return nil
	})
	return nil
}

// OutJsonFile 将data输出到json文件
func OutJsonFile(data interface{}, fileName string) error {
	var (
		f   *os.File
		err error
	)
	if Exists(fileName) { //如果文件存在
		f, err = os.OpenFile(fileName, os.O_APPEND, 0666) //打开文件
	} else {
		f, err = os.Create(fileName) //创建文件
	}
	if err != nil {
		return err
	}
	str, err := AnyToJson(data)
	if err != nil {
		return err
	}
	_, err = io.WriteString(f, str)
	if err != nil {
		return err
	}
	return nil
}

// FileExists 文件是否存在
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func Md5Uppercase(str string) string {
	s := GetMD5Encode(str)
	return strings.ToUpper(s)
}

// picSuffix 图片文件后缀
var picSuffix = []string{".png", ".jpg", ".gif", ".jpeg", "svg", "bmp"}

// IsPic 判断是否是图片
func IsPic(suffix string) bool {
	for _, su := range picSuffix {
		suffix = strings.ToLower(suffix)
		if su == suffix {
			return true
		}
	}
	return false
}

// RandomIntCaptcha 生成 captchaLen 位随机数，理论上会重复
func RandomIntCaptcha(captchaLen int) string {
	var arr string
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < captchaLen; i++ {
		arr = arr + fmt.Sprintf("%d", r.Intn(10))
	}
	return arr
}

// DeepEqual 深度比较任意类型的两个变量的是否相等,类型一样值一样反回true
// 如果元素都是nil，且类型相同，则它们是相等的; 如果它们是不同的类型，它们是不相等的
func DeepEqual(a, b interface{}) bool {
	ra := reflect.Indirect(reflect.ValueOf(a))
	rb := reflect.Indirect(reflect.ValueOf(b))
	if raValid, rbValid := ra.IsValid(), rb.IsValid(); !raValid && !rbValid {
		return reflect.TypeOf(a) == reflect.TypeOf(b)
	} else if raValid != rbValid {
		return false
	}
	return reflect.DeepEqual(ra.Interface(), rb.Interface())
}

// StrLen 获取字符长度
func StrLen(str string) int {
	return utf8.RuneCountInString(str)
}

func RandomString(list []string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return list[r.Intn(len(list))]
}
