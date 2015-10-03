package kademlia

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"time"
)

type Kademlia struct {
	routes    *RoutingTable
	NetworkID string
}

func NewKademlia(self Contact, networkID string) *Kademlia {
	ret := &Kademlia{
		routes:    NewRoutingTable(self),
		NetworkID: networkID,
	}
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
	k.routes.UpdateChan <- request.Sender
	// Pong with sender
	response.Sender = k.routes.self

	return nil
}

func (k *Kademlia) Ping(target Contact) error {
	client, err := dialContact(target)
	if err != nil {
		return err
	}

	req := k.NewPingRequest()
	res := PingResponse{}

	return client.Call("KademliaCore.PingRPC", &req, &res)
}

func (k *Kademlia) Bootstrap(target, self Contact) ([]Contact, error) {
	client, err := dialContact(target)
	if err != nil {
		return nil, err
	}

	req := k.NewFindNodeRequest(self)
	res := FindNodeResponse{}

	err = client.Call("KademliaCore.FindNodeRPC", &req, &res)
	if err != nil {
		return nil, err
	}

	return res.Contacts, nil
}

func (k *Kademlia) FindNode(target Contact) ([]Contact, error) {
	client, err := dialContact(target)
	if err != nil {
		return nil, err
	}

	req := k.NewFindNodeRequest(target)
	res := FindNodeResponse{}

	err = client.Call("KademliaCore.FindNodeRPC", &req, &res)
	if err != nil {
		return nil, err
	}

	return res.Contacts, nil
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

/*
 * PING
 */

type PingRequest struct {
	RPCHeader
}

func (k *Kademlia) NewPingRequest() PingRequest {
	return PingRequest{
		RPCHeader{
			Sender:    k.routes.self,
			NetworkID: k.NetworkID,
		},
	}
}

type PingResponse struct {
	RPCHeader
}

func (kc *KademliaCore) PingRPC(req PingRequest, res *PingResponse) error {
	return kc.kad.HandleRPC(req.RPCHeader, &res.RPCHeader)
}

/*
 * FIND NODE
 */

type FindNodeRequest struct {
	RPCHeader
	Target Contact
}

func (k *Kademlia) NewFindNodeRequest(target Contact) FindNodeRequest {
	return FindNodeRequest{
		RPCHeader: RPCHeader{
			Sender:    k.routes.self,
			NetworkID: k.NetworkID,
		},
		Target: target,
	}
}

type FindNodeResponse struct {
	RPCHeader
	Contacts Contacts
}

func (kc *KademliaCore) FindNodeRPC(req FindNodeRequest, res *FindNodeResponse) error {
	err := kc.kad.HandleRPC(req.RPCHeader, &res.RPCHeader)
	if err != nil {
		return err
	}

	contactsReturnChan := make(chan Contacts)
	kc.kad.routes.LookupRequestChan <- contactsReturnChan
	contactsReturnChan <- Contacts{req.Target}

	contacts := <-contactsReturnChan
	res.Contacts = contacts

	return nil
}
