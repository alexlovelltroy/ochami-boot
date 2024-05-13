package main

import (
	"sync"
)

type MemoryStore struct {
	sync.RWMutex
	Nodes map[string]*Node // key is MAC address as string
	data  map[string]interface{}
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		Nodes: make(map[string]*Node),
	}
}

func (store *MemoryStore) AddNode(node *Node) {
	store.Lock()
	defer store.Unlock()
	store.Nodes[node.MAC.String()] = node
}

func (store *MemoryStore) GetNode(mac string) (*Node, bool) {
	store.RLock()
	defer store.RUnlock()
	node, exists := store.Nodes[mac]
	return node, exists
}

func (store *MemoryStore) DeleteNode(mac string) {
	store.Lock()
	defer store.Unlock()
	delete(store.Nodes, mac)
}

func (store *MemoryStore) Set(key string, value interface{}) {
	store.Lock()
	defer store.Unlock()
	store.data[key] = value
}

func (store *MemoryStore) Get(key string) (interface{}, bool) {
	store.RLock()
	defer store.RUnlock()
	val, exists := store.data[key]
	return val, exists
}

func (store *MemoryStore) Delete(key string) {
	store.Lock()
	defer store.Unlock()
	delete(store.data, key)
}
