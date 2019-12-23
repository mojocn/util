package mache

import (
	"time"
)

var mStore = newMemoryStore(1024*1024*10, time.Hour*24*30)

func Set(key string, v interface{}) {
	mStore.set(key, v)
}

func Get(key string, clear bool) interface{} {
	return mStore.get(key, clear)
}

//GetOrSet getOrSet
func GetOrSet(key string, setValueFunc func() interface{}) interface{} {
	vv := Get(key, false)
	if vv == nil {
		vv = setValueFunc()
		if vv != nil {
			Set(key, vv)
		}
	}
	return vv
}
