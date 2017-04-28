package rand

import (
	"math/rand"
	"sync"
	"time"
)

type SafeRand struct {
	random *rand.Rand
	mutex  *sync.Mutex
}

func NewSafeRand() *SafeRand {
	return &SafeRand{
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
		mutex:  &sync.Mutex{},
	}
}

func (self *SafeRand) Intn(n int) int {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	return self.random.Intn(n)
}

func (self *SafeRand) Float64() float64 {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	return self.random.Float64()
}

func (self *SafeRand) Perm(n int) []int {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	return self.random.Perm(n)
}
