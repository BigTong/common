package filter

import (
	"sync"

	"github.com/willf/bloom"
)

type BloomFilter struct {
	bloomFilter *bloom.BloomFilter
	mutex       *sync.RWMutex
}

func NewBloomFilter(maxItemNum int, kFp float64) *BloomFilter {
	return &BloomFilter{
		mutex:       &sync.RWMutex{},
		bloomFilter: bloom.NewWithEstimates(uint(maxItemNum), kFp),
	}

}

func (self *BloomFilter) NeedFilter(src string) bool {
	self.mutex.RLock()
	defer self.mutex.RUnlock()

	if self.bloomFilter.TestString(src) {
		return true
	}
	return false
}

func (self *BloomFilter) SetFilter(src string) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.bloomFilter.AddString(src)
}

func (self *BloomFilter) Reset() {

}

func (self *BloomFilter) Dump() {

}
