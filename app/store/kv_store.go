package store

type KVStore struct {
	Entries map[string]string
}

func NewKVStore() *KVStore {
	return &KVStore{
		Entries: make(map[string]string),
	}
}

func (kv *KVStore) Get(key string) (string, bool) {
	entry, ok := kv.Entries[key]
	return entry, ok
}

func (kv *KVStore) Set(key string, value string) {
	kv.Entries[key] = value
}
