package main

// NodeStorage is the interface that wraps the basic methods needed for a store.
type NodeStorage interface {
	GetNode(mac string) (*Node, bool)
	SetNode(mac string, node *Node)
	DeleteNode(mac string)
}
