package types

import (
	"sync"
)

type FunctionCache struct {
	f          Function
	m          map[interface{}]FunctionResult
	lock       *sync.RWMutex
	InProgress map[Job]bool
	Waiting    map[Job][]chan FunctionResult
}

func NewFunctionCache(f Function) *FunctionCache {
	return &FunctionCache{
		f:          f,
		m:          make(map[interface{}]FunctionResult),
		lock:       &sync.RWMutex{},
		InProgress: make(map[Job]bool),
		Waiting:    make(map[Job][]chan FunctionResult),
	}
}

func (fc FunctionCache) Get(key interface{}) (FunctionResult, bool) {
	fc.lock.RLock()
	result, exists := fc.m[key]
	fc.lock.RUnlock()
	return result, exists
}

func (fc FunctionCache) Set(key interface{}, value FunctionResult) {
	fc.lock.Lock()
	fc.m[key] = value
	fc.lock.Unlock()
}
