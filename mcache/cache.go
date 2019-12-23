package mache

import (
	"container/list"
	"sync"
	"time"
)

// expValue stores timestamp and id of captchas. It is used in the list inside
// memoryStore for indexing generated captchas by timestamp to enable garbage
// collection of expired captchas.
type keyTime struct {
	timestamp time.Time
	id        string
}

// memoryStore is an internal store for captcha ids and their values.
type memoryStore struct {
	sync.RWMutex
	kv       map[string]interface{}
	idByTime *list.List
	// Number of items stored since last collection.
	numStored int
	// Number of saved items that triggers collection.
	collectNum int
	// Expiration time of captchas.
	expiration time.Duration
}

// NewMemoryStore returns a new standard memory store for captchas with the
// given collection threshold and expiration time (duration). The returned
// store must be registered with SetCustomStore to replace the default one.
func newMemoryStore(collectNum int, expiration time.Duration) *memoryStore {
	s := new(memoryStore)
	s.kv = make(map[string]interface{})
	s.idByTime = list.New()
	s.collectNum = collectNum
	s.expiration = expiration
	return s
}

func (s *memoryStore) set(id string, value interface{}) {
	s.Lock()
	if value == nil {
		delete(s.kv, id)
	} else {
		s.kv[id] = value
		s.idByTime.PushBack(keyTime{time.Now(), id})
		s.numStored++
	}
	s.Unlock()

	if s.numStored > s.collectNum {
		go s.collect()
	}
}

func (s *memoryStore) get(id string, clear bool) (value interface{}) {
	if !clear {
		s.RLock()
		defer s.RUnlock()
	} else {
		s.Lock()
		defer s.Unlock()
	}
	value, ok := s.kv[id]
	if !ok {
		return
	}
	if clear {
		delete(s.kv, id)
	}
	return
}

func (s *memoryStore) collect() {
	now := time.Now()
	s.Lock()
	defer s.Unlock()
	for e := s.idByTime.Front(); e != nil; {
		e = s.collectOne(e, now)
	}
	s.numStored = len(s.kv)
}

func (s *memoryStore) collectOne(e *list.Element, specifyTime time.Time) *list.Element {

	ev, ok := e.Value.(keyTime)
	if !ok {
		return nil
	}

	if ev.timestamp.Add(s.expiration).Before(specifyTime) {
		delete(s.kv, ev.id)
		next := e.Next()
		s.idByTime.Remove(e)
		return next
	}
	return nil
}
