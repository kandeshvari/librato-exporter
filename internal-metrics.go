package main

import (
	"sync"
	"sync/atomic"
)

type countersMap struct {
	sync.Map
}

var internalMetrics countersMap

func New() *countersMap {
	return &countersMap{}
}

func (mc *countersMap) Put(key string, value int64) error {
	v, _ := mc.LoadOrStore(key, new(int64))
	cntr := v.(*int64)
	atomic.StoreInt64(cntr, value)
	return nil
}
