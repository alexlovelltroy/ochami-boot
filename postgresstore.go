package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/lib/pq"

	log "github.com/sirupsen/logrus"
)

type PostgresStore struct {
	db *sql.DB
}

func (ps *PostgresStore) GetNode(mac string) (*Node, bool) {
	row := ps.db.QueryRow("SELECT * FROM nodes WHERE mac = $1", mac)

	node := &Node{}
	var optionsV4, optionsV6 []byte

	err := row.Scan(&node.MAC, &node.IPv4, &node.IPv6, &optionsV4, &optionsV6)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false
		}
		// handle other errors
	}

	json.Unmarshal(optionsV4, &node.OptionsV4)
	json.Unmarshal(optionsV6, &node.OptionsV6)

	return node, true
}

func (ps *PostgresStore) SetNode(mac string, node *Node) {
	optionsV4, _ := json.Marshal(node.OptionsV4)
	optionsV6, _ := json.Marshal(node.OptionsV6)

	_, err := ps.db.Exec("INSERT INTO nodes (mac, ipv4, ipv6, dhcpv4_options, dhcpv6_options) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (mac) DO UPDATE SET ipv4 = $2, ipv6 = $3, dhcpv4_options = $4, dhcpv6_options = $5",
		mac, node.IPv4, node.IPv6, optionsV4, optionsV6)
	if err != nil {
		log.Errorf("failed to insert node: %v", err)
	}
}

func (ps *PostgresStore) DeleteNode(mac string) {
	_, err := ps.db.Exec("DELETE FROM nodes WHERE mac = $1", mac)
	if err != nil {
		log.Errorf("failed to delete node: %v", err)
	}
}

func (ps *PostgresStore) AddNode(node *Node) {
	ps.SetNode(string(node.MAC), node)
}

func (ps *PostgresStore) createTablesIfNotExist() error {
	query := `
	CREATE TABLE IF NOT EXISTS nodes (
		mac MACADDR PRIMARY KEY,
		ipv4 INET,
		ipv6 INET,
		dhcpv4_options BYTEA,
		dhcpv6_options BYTEA
	)
	`

	_, err := ps.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	return nil
}

func NewPostgresStore(connStr string) (NodeStorage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	ps := &PostgresStore{db: db}

	err = ps.createTablesIfNotExist()
	if err != nil {
		return nil, err
	}

	return ps, nil
}
