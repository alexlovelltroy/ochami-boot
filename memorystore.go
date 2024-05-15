package main

import (
	"sync"
	"time"
)

type BootAttempt struct {
	Timestamp time.Time
	Known     bool
}

type MacMemoryStore struct {
	sync.RWMutex
	Macs map[string][]BootAttempt // key is MAC address as string
}

func NewMacMemoryStore() *MacMemoryStore {
	return &MacMemoryStore{
		Macs: make(map[string][]BootAttempt),
	}
}

func (store *MacMemoryStore) AddBootAttempt(mac string, timestamp time.Time, known bool) {
	store.Lock()
	defer store.Unlock()
	store.Macs[mac] = append(store.Macs[mac], BootAttempt{Timestamp: timestamp, Known: known})
}

func (store *MacMemoryStore) GetBootAttempts(mac string) int {
	store.RLock()
	defer store.RUnlock()
	bootAttempts, exists := store.Macs[mac]
	if exists {
		return len(bootAttempts)
	}
	return 0
}

func (store *MacMemoryStore) GetLastBootAttempt(mac string) (BootAttempt, bool) {
	store.RLock()
	defer store.RUnlock()
	bootAttempts, exists := store.Macs[mac]
	if exists {
		return bootAttempts[len(bootAttempts)-1], true
	}
	return BootAttempt{}, false
}

func (store *MacMemoryStore) GetUnknownMacs() []string {
	store.RLock()
	defer store.RUnlock()
	var unknownMacs []string
	for mac, bootAttempts := range store.Macs {
		if len(bootAttempts) == 0 || bootAttempts[len(bootAttempts)-1].Known {
			unknownMacs = append(unknownMacs, mac)
		}
	}
	return unknownMacs
}

type NodeMemoryStore struct {
	sync.RWMutex
	Nodes map[string]*Node // key is MAC address as string
}

func NewNodeMemoryStore() *NodeMemoryStore {
	return &NodeMemoryStore{
		Nodes: make(map[string]*Node),
	}
}

func (store *NodeMemoryStore) AddNode(node *Node) {
	store.Lock()
	defer store.Unlock()
	store.Nodes[node.MAC.String()] = node
}

func (store *NodeMemoryStore) GetNode(mac string) (*Node, bool) {
	store.RLock()
	defer store.RUnlock()
	node, exists := store.Nodes[mac]
	return node, exists
}

func (store *NodeMemoryStore) DeleteNode(mac string) {
	store.Lock()
	defer store.Unlock()
	delete(store.Nodes, mac)
}

func (store *NodeMemoryStore) SetNode(mac string, node *Node) {
	store.Lock()
	defer store.Unlock()
	store.Nodes[mac] = node
}
