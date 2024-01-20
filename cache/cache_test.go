package cache

import (
	"container/list"
	"reflect"
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		capacity int
	}
	tests := []struct {
		name    string
		args    args
		want    *Cache
		wantErr bool
	}{
		{
			name: "zero capacity",
			args: args{
				capacity: 0,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "negative capacity",
			args: args{
				capacity: -5,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid capacity",
			args: args{
				capacity: 10,
			},
			want: &Cache{
				capacity: 10,
				cache:    list.New(),
				elements: make(map[string]*list.Element),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.capacity)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_Set(t *testing.T) {

	type args struct {
		key   []byte
		value []byte
	}

	tests := []struct {
		name         string
		initialCache *Cache
		args         args
		want         *Cache
	}{
		{
			name: "add new item to empty cache",
			initialCache: &Cache{
				capacity: 3,
				cache:    list.New(),
				elements: make(map[string]*list.Element),
				mutex:    sync.RWMutex{},
			},
			args: args{
				key:   []byte("key1"),
				value: []byte("value1"),
			},
			want: &Cache{
				capacity: 3,
				cache:    list.New(),
				elements: map[string]*list.Element{string([]byte("key1")): {Value: Item{
					Key:   []byte("key1"),
					Value: []byte("value1"),
				}}},
				mutex: sync.RWMutex{},
			},
		},
		{
			name: "update existing item",
			initialCache: &Cache{
				capacity: 3,
				cache:    list.New(),
				elements: map[string]*list.Element{string([]byte("existing_key")): {Value: Item{Key: []byte("existing_key"), Value: []byte("old_value")}}},
				mutex:    sync.RWMutex{},
			},
			args: args{
				key:   []byte("existing_key"),
				value: []byte("new_value"),
			},
			want: &Cache{
				capacity: 3,
				cache:    list.New(),
				elements: map[string]*list.Element{string([]byte("existing_key")): {Value: Item{Key: []byte("existing_key"), Value: []byte("new_value")}}},
				mutex:    sync.RWMutex{},
			},
		},

		{
			name: "evict least recently used item",
			initialCache: &Cache{
				capacity: 2,
				cache:    list.New(),
				elements: map[string]*list.Element{
					string([]byte("key1")): {Value: Item{Key: []byte("key1"), Value: []byte("value1")}},
					string([]byte("key2")): {Value: Item{Key: []byte("key2"), Value: []byte("value2")}},
				},
				mutex: sync.RWMutex{},
			},
			args: args{
				key:   []byte("new_key"),
				value: []byte("new_value"),
			},
			want: &Cache{
				capacity: 2,
				cache:    list.New(),
				elements: map[string]*list.Element{
					string([]byte("key2")):    {Value: Item{Key: []byte("key2"), Value: []byte("value2")}},
					string([]byte("new_key")): {Value: Item{Key: []byte("new_key"), Value: []byte("new_value")}},
				},
				mutex: sync.RWMutex{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.initialCache.Set(tt.args.key, tt.args.value)
			for k, v := range tt.want.elements {
				if !reflect.DeepEqual(tt.initialCache.elements[k].Value, v) {
					t.Errorf("Set() = %v, want %v", tt.initialCache.elements[k].Value, v)
				}
			}

			for e := tt.initialCache.cache.Front(); e != nil; e = e.Next() {
				tt.want.cache.PushFront(e.Value.(*list.Element))
			}
			if !reflect.DeepEqual(tt.initialCache.cache, tt.want.cache) {
				t.Errorf("Set() cache = %v, want %v", tt.initialCache.cache, tt.want.cache)
			}

		})
	}
}
