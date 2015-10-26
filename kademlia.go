package kademlia

import (
	"encoding/hex"
	"errors"
	"fmt"
	db "github.com/syndtr/goleveldb/leveldb"
	"log"
	"net"
	"net/rpc"
	"time"
)

const (
	VALUES_DB_PATH = "db/values-"
)

type Kademlia struct {
	routes    *RoutingTable
	valuesDB  *db.DB
	NetworkID string
}

func NewKademlia(self Contact, networkID string) *Kademlia {
	ret := &Kademlia{
		routes:    NewRoutingTable(self),
		valuesDB:  nil,
		NetworkID: networkID,
	}

	hexID := hex.EncodeToString(self.ID[:])
	conn, err := db.OpenFile(VALUES_DB_PATH+hexID, nil)
	if err != nil {
		log.Println(err)
		panic("Unable to open values database")
	}

	defer conn.Close()
	ret.valuesDB = conn

	return ret
}

// Generic RPC base
type RPCHeader struct {
	Sender    Contact
	NetworkID string
}

// Every RPC updates routing tables in Kademlia
func (k *Kademlia) HandleRPC(request RPCHeader, response *RPCHeader) error {
	if request.NetworkID != k.NetworkID {
		return errors.New(fmt.Sprintf("Expected Network ID %s, go %s", k.NetworkID, request.NetworkID))
	}

	// Update routing table for all incoming RPCs
	k.routes.Update(request.Sender)
	// Pong with sender
	response.Sender = k.routes.self

	return nil
}

func dialContact(contact Contact) (*rpc.Client, error) {
	connection, err := net.DialTimeout("tcp", contact.Address, 5*time.Second)
	if err != nil {
		return nil, err
	}

	return rpc.NewClient(connection), nil
}

func (k *Kademlia) Serve() error {
	rpc.Register(&KademliaCore{k})

	l, err := net.Listen("tcp", k.routes.self.Address)
	if err != nil {
		return err
	}

	go rpc.Accept(l)

	return nil
}

/*
 * KademliaCore
 * Handles RPC interactions between client/server
 */

type KademliaCore struct {
	kad *Kademlia
}
