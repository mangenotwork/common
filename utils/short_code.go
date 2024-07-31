package utils

import (
	randc "crypto/rand"
	"errors"
	"fmt"
	"math"
	randm "math/rand"
	"sync"
	"time"
)

// NewShortCode 获取一个不重复的短code;常用与短链code生成;
func NewShortCode() string {
	id, err := generate()
	if err == nil {
		return id
	}
	return ""
}

const DefaultABC = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type abcTable struct {
	alphabet []rune
}

type sCode struct {
	abc    abcTable
	worker uint
	epoch  time.Time
	ms     uint
	count  uint
	mx     sync.Mutex
}

var shortCode *sCode

func init() {
	shortCode = mustNew(0, DefaultABC, 1)
}

func generate() (string, error) {
	return shortCode.generate()
}

func newShortCode(worker uint8, alphabet string, seed uint64) (*sCode, error) {
	if worker > 31 {
		return nil, errors.New("expected worker in the range [0,31]")
	}
	abc, err := newAbc(alphabet, seed)
	if err == nil {
		sid := &sCode{
			abc:    abc,
			worker: uint(worker),
			epoch:  time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC),
			ms:     0,
			count:  0,
		}
		return sid, nil
	}
	return nil, err
}

func mustNew(worker uint8, alphabet string, seed uint64) *sCode {
	sid, err := newShortCode(worker, alphabet, seed)
	if err == nil {
		return sid
	}
	panic(err)
}

func (sid *sCode) generate() (string, error) {
	return sid.generateInternal(nil, sid.epoch)
}

func (sid *sCode) mustGenerate() string {
	id, err := sid.generate()
	if err == nil {
		return id
	}
	panic(err)
}

func (sid *sCode) generateInternal(tm *time.Time, epoch time.Time) (string, error) {
	ms, count := sid.getMsAndCounter(tm, epoch)
	idRunes := make([]rune, 7)
	if tmp, err := sid.abc.enCode(ms, 8, 5); err == nil {
		copy(idRunes, tmp)
	} else {
		return "", err
	}
	if tmp, err := sid.abc.enCode(sid.worker, 1, 5); err == nil {
		idRunes[6] = tmp[0]
	} else {
		return "", err
	}
	if count > 0 {
		if countRunes, err := sid.abc.enCode(count, 0, 6); err == nil {
			idRunes = append(idRunes, countRunes...)
		} else {
			return "", err
		}
	}
	return string(idRunes), nil
}

func (sid *sCode) getMsAndCounter(tm *time.Time, epoch time.Time) (uint, uint) {
	sid.mx.Lock()
	defer sid.mx.Unlock()
	var ms uint
	if tm != nil {
		ms = uint(tm.Sub(epoch).Nanoseconds() / 1000000)
	} else {
		ms = uint(time.Now().Sub(epoch).Nanoseconds() / 1000000)
	}
	if ms == sid.ms {
		sid.count++
	} else {
		sid.count = 0
		sid.ms = ms
	}
	return sid.ms, sid.count
}

func newAbc(alphabet string, seed uint64) (abcTable, error) {
	runes := []rune(alphabet)
	if len(runes) != len(DefaultABC) {
		return abcTable{}, fmt.Errorf("alphabet must contain %v unique characters", len(DefaultABC))
	}
	if nonUnique(runes) {
		return abcTable{}, errors.New("alphabet must contain unique characters only")
	}
	abc := abcTable{alphabet: nil}
	abc.shuffle(alphabet, seed)
	return abc, nil
}

func nonUnique(runes []rune) bool {
	found := make(map[rune]struct{})
	for _, r := range runes {
		if _, seen := found[r]; !seen {
			found[r] = struct{}{}
		}
	}
	return len(found) < len(runes)
}

func (abc *abcTable) shuffle(alphabet string, seed uint64) {
	source := []rune(alphabet)
	for len(source) > 1 {
		seed = (seed*9301 + 49297) % 233280
		i := int(seed * uint64(len(source)) / 233280)
		abc.alphabet = append(abc.alphabet, source[i])
		source = append(source[:i], source[i+1:]...)
	}
	abc.alphabet = append(abc.alphabet, source[0])
}

func (abc *abcTable) enCode(val, nSymbols, digits uint) ([]rune, error) {
	if digits < 4 || 6 < digits {
		return nil, fmt.Errorf("allowed digits range [4,6], found %v", digits)
	}

	var computedSize uint = 1
	if val >= 1 {
		computedSize = uint(math.Log2(float64(val)))/digits + 1
	}
	if nSymbols == 0 {
		nSymbols = computedSize
	} else if nSymbols < computedSize {
		return nil, fmt.Errorf("cannot accommodate data, need %v digits, got %v", computedSize, nSymbols)
	}

	mask := 1<<digits - 1

	random := make([]int, int(nSymbols))
	if digits < 6 {
		copy(random, maskedRandomInt(len(random), 0x3e-mask))
	}

	res := make([]rune, int(nSymbols))
	for i := range res {
		shift := digits * uint(i)
		index := (int(val>>shift) & mask) | random[i]
		if index >= 62 {
			continue
		}
		res[i] = abc.alphabet[index]
	}
	return res, nil
}

func (abc *abcTable) mustEncode(val, size, digits uint) []rune {
	res, err := abc.enCode(val, size, digits)
	if err == nil {
		return res
	}
	panic(err)
}

func maskedRandomInt(size, mask int) []int {
	intList := make([]int, size)
	bytes := make([]byte, size)
	if _, err := randc.Read(bytes); err == nil {
		for i, b := range bytes {
			intList[i] = int(b) & mask
		}
	} else {
		for i := range intList {
			intList[i] = randm.Intn(0xff) & mask
		}
	}
	return intList
}
