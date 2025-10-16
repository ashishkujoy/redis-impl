package store

import "time"

type EntryValue struct {
	value     string
	createdAt time.Time
	px        int
}
type KVStore struct {
	Entries map[string]EntryValue
}

func NewKVStore() *KVStore {
	return &KVStore{
		Entries: make(map[string]EntryValue),
	}
}

func (kv *KVStore) Get(key string) (string, bool) {
	entry, ok := kv.Entries[key]
	if !ok {
		return "", false
	}
	if entry.px != -1 && time.Since(entry.createdAt) >= time.Duration(entry.px)*time.Millisecond {
		delete(kv.Entries, key)
		return "", false
	}
	return entry.value, ok
}

func (kv *KVStore) Set(key string, value string, px int) {
	kv.Entries[key] = EntryValue{value, time.Now(), px}
}
