package safemap

import (
	"encoding/json"
	"sync"

	"github.com/BigTong/common/file"
)

type Map struct {
	data map[string]interface{}
	lock *sync.RWMutex
}

func NewMap() *Map {
	return &Map{
		data: make(map[string]interface{}),
		lock: &sync.RWMutex{},
	}
}

func (self *Map) CopyOf(m map[string]interface{}) {
	self.lock.Lock()
	defer self.lock.Unlock()
	if len(self.data) != 0 {
		self.data = make(map[string]interface{})
	}
	for k, v := range m {
		self.data[k] = v
	}
}

func NewMapFromConfig(config string) *Map {
	ret := &Map{
		data: make(map[string]interface{}),
		lock: &sync.RWMutex{},
	}
	conf := file.ReadFileToString(config)
	data := make(map[string]int64)
	json.Unmarshal([]byte(conf), &data)
	for k, v := range data {
		ret.data[k] = v
	}
	return ret
}

func (self *Map) Set(key string, val interface{}) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.data[key] = val
}

func (self *Map) Get(key string) interface{} {
	self.lock.Lock()
	defer self.lock.Unlock()
	val, ok := self.data[key]
	if ok {
		return val
	}
	return nil
}

func (self *Map) GetInt64(key string) (int64, bool) {
	val := self.Get(key)
	if val == nil {
		return 0, false
	}

	switch val.(type) {
	case int64:
		return val.(int64), true
	}
	return 0, false
}

func (self *Map) GetString(key string) (string, bool) {
	val := self.Get(key)
	if val == nil {
		return "", false
	}

	switch val.(type) {
	case string:
		return val.(string), true
	}
	return "", false
}

func (self *Map) GetInt(key string) (int, bool) {
	val := self.Get(key)
	if val == nil {
		return 0, false
	}

	switch val.(type) {
	case int:
		return val.(int), true
	}
	return 0, false
}
