package main

import (
	"sync"

	_ "github.com/lib/pq"
)

type MemoryStore struct {
	sync.RWMutex
	Nodes map[string]*Node // key is MAC address as string
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

func (store *MemoryStore) SetNode(mac string, node *Node) {
	store.Lock()
	defer store.Unlock()
	store.Nodes[mac] = node
}
