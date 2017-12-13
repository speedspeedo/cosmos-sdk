package store

import (
	"bytes"
	"sort"
	"sync"
)

// If value is nil but deleted is false, it means the parent doesn't have the
// key.  (No need to delete upon Write())
type cValue struct {
	value   []byte
	deleted bool
	dirty   bool
}

// cacheKVStore wraps an in-memory cache around an underlying KVStore.
type cacheKVStore struct {
	mtx    sync.Mutex
	cache  map[string]cValue
	parent KVStore
}

var _ CacheKVStore = (*cacheKVStore)(nil)

func NewCacheKVStore(parent KVStore) *cacheKVStore {

	ci := &cacheKVStore{
		cache:  make(map[string]cValue),
		parent: parent,
	}

	return ci
}

func (ci *cacheKVStore) Get(key []byte) (value []byte) {
	ci.mtx.Lock()
	defer ci.mtx.Unlock()
	ci.assertValidKey(key)

	cacheValue, ok := ci.cache[string(key)]
	if !ok {
		value = ci.parent.Get(key)
		ci.setCacheValue(key, value, false, false)
	} else {
		value = cacheValue.value
	}

	return value
}

func (ci *cacheKVStore) Set(key []byte, value []byte) {
	ci.mtx.Lock()
	defer ci.mtx.Unlock()
	ci.assertValidKey(key)

	ci.setCacheValue(key, value, false, true)
}

func (ci *cacheKVStore) Has(key []byte) bool {
	value := ci.Get(key)
	return value != nil
}

func (ci *cacheKVStore) Delete(key []byte) {
	ci.mtx.Lock()
	defer ci.mtx.Unlock()
	ci.assertValidKey(key)

	ci.setCacheValue(key, nil, true, true)
}

// Write writes pending updates to the parent database and clears the cache.
func (ci *cacheKVStore) Write() {
	ci.mtx.Lock()
	defer ci.mtx.Unlock()

	// We need a copy of all of the keys.
	// Not the best, but probably not a bottleneck depending.
	keys := make([]string, 0, len(ci.cache))
	for key, dbValue := range ci.cache {
		if dbValue.dirty {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	// TODO in tmlibs/db we use Batch to write atomically.
	// Consider allowing usage of Batch.
	for _, key := range keys {
		cacheValue := ci.cache[key]
		if cacheValue.deleted {
			ci.parent.Delete([]byte(key))
		} else if cacheValue.value == nil {
			// Skip, it already doesn't exist in parent.
		} else {
			ci.parent.Set([]byte(key), cacheValue.value)
		}
	}

	// Clear the cache
	ci.cache = make(map[string]cValue)
}

//----------------------------------------
// To cache-wrap this cacheKVStore further.

func (ci *cacheKVStore) CacheWrap() CacheWrap {
	return NewCacheKVStore(ci)
}

//----------------------------------------
// Iteration

func (ci *cacheKVStore) Iterator(start, end []byte) Iterator {
	return ci.iterator(start, end, true)
}

func (ci *cacheKVStore) ReverseIterator(start, end []byte) Iterator {
	return ci.iterator(start, end, false)
}

func (ci *cacheKVStore) iterator(start, end []byte, ascending bool) Iterator {
	var parent, cache Iterator
	if ascending {
		parent = ci.parent.Iterator(start, end)
	} else {
		parent = ci.parent.ReverseIterator(start, end)
	}
	items := ci.dirtyItems(ascending)
	cache = newMemIterator(start, end, items)
	return newCacheMergeIterator(parent, cache, ascending)
}

// Constructs a slice of dirty items, to use w/ memIterator.
func (ci *cacheKVStore) dirtyItems(ascending bool) []KVPair {
	items := make([]KVPair, 0, len(ci.cache))
	for key, cacheValue := range ci.cache {
		if !cacheValue.dirty {
			continue
		}
		items = append(items,
			KVPair{[]byte(key), cacheValue.value})
	}
	sort.Slice(items, func(i, j int) bool {
		if ascending {
			return bytes.Compare(items[i].Key, items[j].Key) < 0
		} else {
			return bytes.Compare(items[i].Key, items[j].Key) > 0
		}
	})
	return items
}

//----------------------------------------
// etc

func (ci *cacheKVStore) assertValidKey(key []byte) {
	if key == nil {
		panic("key is nil")
	}
}

// Only entrypoint to mutate ci.cache.
func (ci *cacheKVStore) setCacheValue(key, value []byte, deleted bool, dirty bool) {
	cacheValue := cValue{
		value:   value,
		deleted: deleted,
		dirty:   dirty,
	}
	ci.cache[string(key)] = cacheValue
}
