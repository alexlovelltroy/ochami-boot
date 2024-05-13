package main

// NodeStorage is the interface that wraps the basic methods needed for a store.
type NodeStorage interface {
	AddNode(node *Node)
	GetNode(mac string) (*Node, bool)
	SetNode(mac string, node *Node)
	DeleteNode(mac string)
}
