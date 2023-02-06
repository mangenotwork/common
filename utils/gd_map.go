package utils

import (
	"sync"
)

// GDMaper 固定顺序 Map 接口
type GDMaper interface {
	Add(key string, value interface{}) *GDMap
	Get(key string) interface{}
	Del(key string) *GDMap
	Len() int
	KeyList() []string
	AddMap(data map[string]interface{}) *GDMap
	Range(f func(k string, v interface{})) *GDMap
	RangeAt(f func(id int, k string, v interface{})) *GDMap
	CheckValue(value interface{}) bool // 检查是否存在某个值
	Reverse()                          //反序
}

// GDMap 固定顺序map
type GDMap struct {
	mux     sync.Mutex
	data    map[string]interface{}
	keyList []string
	size    int
}

// NewGDMap ues: NewGDMap().Add(k,v)
func NewGDMap() *GDMap {
	return &GDMap{
		data:    make(map[string]interface{}),
		keyList: make([]string, 0),
		size:    0,
	}
}

// Add  添加kv
func (m *GDMap) Add(key string, value interface{}) *GDMap {
	m.mux.Lock()
	defer m.mux.Unlock()
	if _, ok := m.data[key]; ok {
		m.data[key] = value
		return m
	}
	m.keyList = append(m.keyList, key)
	m.size++
	m.data[key] = value
	return m
}

// Get 通过key获取值
func (m *GDMap) Get(key string) interface{} {
	m.mux.Lock()
	defer m.mux.Unlock()
	return m.data[key]
}

// Del 删除指定key的值
func (m *GDMap) Del(key string) *GDMap {
	m.mux.Lock()
	defer m.mux.Unlock()
	if _, ok := m.data[key]; ok {
		delete(m.data, key)
		for i := 0; i < m.size; i++ {
			if m.keyList[i] == key {
				m.keyList = append(m.keyList[:i], m.keyList[i+1:]...)
				m.size--
				return m
			}
		}
	}
	return m
}

// Len map的长度
func (m *GDMap) Len() int {
	return m.size
}

// KeyList 打印map所有的key
func (m *GDMap) KeyList() []string {
	return m.keyList
}

// AddMap 写入map
func (m *GDMap) AddMap(data map[string]interface{}) *GDMap {
	for k, v := range data {
		m.Add(k, v)
	}
	return m
}

// Range 遍历map
func (m *GDMap) Range(f func(k string, v interface{})) *GDMap {
	for i := 0; i < m.size; i++ {
		f(m.keyList[i], m.data[m.keyList[i]])
	}
	return m
}

// RangeAt Range 遍历map含顺序id
func (m *GDMap) RangeAt(f func(id int, k string, v interface{})) *GDMap {
	for i := 0; i < m.size; i++ {
		f(i, m.keyList[i], m.data[m.keyList[i]])
	}
	return m
}

// CheckValue 查看map是否存在指定的值
func (m *GDMap) CheckValue(value interface{}) bool {
	m.mux.Lock()
	defer m.mux.Unlock()
	for i := 0; i < m.size; i++ {
		if m.data[m.keyList[i]] == value {
			return true
		}
	}
	return false
}

// Reverse map反序
func (m *GDMap) Reverse() {
	for i, j := 0, len(m.keyList)-1; i < j; i, j = i+1, j-1 {
		m.keyList[i], m.keyList[j] = m.keyList[j], m.keyList[i]
	}
}
