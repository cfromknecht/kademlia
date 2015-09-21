package kademlia

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
)

type Kademlia struct {
	routes    *RoutingTable
	NetworkID string
}

func NewKademlia(self Contact, networkID string) (ret *Kademlia) {
	ret = new(Kademlia)
	ret.routes = NewRoutingTable(self)
	ret.NetworkID = networkID
	return
}

type RPCHeader struct {
	Sender    Contact
	NetworkID string
}

func (k *Kademlia) HandleRPC(request, response *RPCHeader) error {
	if request.NetworkID != k.NetworkID {
		return errors.New(fmt.Sprintf("Expected Network ID %s, go %s", k.NetworkID, request.NetworkID))
	}

	k.routes.UpdateChan <- request.Sender

	response.Sender = k.routes.self
	return nil
}

type KademliaCore struct {
	kad *Kademlia
}

type PingRequest struct {
	RPCHeader
}

type PingResponse struct {
	RPCHeader
}

func (kc *KademliaCore) Ping(args *PingRequest, response *PingResponse) (err error) {
	err = kc.kad.HandleRPC(&args.RPCHeader, &response.RPCHeader)
	return
}

func (k *Kademlia) Serve() (err error) {
	rpc.Register(&KademliaCore{k})

	rpc.HandleHTTP()
	l, err := net.Listen("tcp", k.routes.self.address)
	check(err)

	go http.Serve(l, nil)

	return
}

type FindNodeRequest struct {
	RPCHeader
	target NodeID
}

type FindNodeResponse struct {
	RPCHeader
	contacts Contacts
}

func (kc *KademliaCore) FindNode(args *FindNodeRequest, response *FindNodeResponse) (err error) {
	if err = kc.kad.HandleRPC(&args.RPCHeader, &response.RPCHeader); err == nil {
		//contactRecords := kc.kad.routes.FindClosest(args.target, BucketSize)
		response.contacts = Contacts{}
	}

	return
}

func (k *Kademlia) Call(contact Contact, method string, args, reply interface{}) (err error) {
	if client, err := rpc.DialHTTP("tcp", contact.address); err == nil {
		if err = client.Call(method, args, reply); err == nil {
			k.routes.UpdateChan <- contact
		}
	}
	return
}

func (k *Kademlia) sendQuery(node Contact, target NodeID, done chan Contacts) {
	args := FindNodeRequest{RPCHeader{k.routes.self, k.NetworkID}, target}
	reply := FindNodeResponse{}

	if err := k.Call(node, "KademliaCore.FindNode", &args, &reply); err == nil {
		fmt.Println("Returning contacts: ", reply.contacts)
		done <- reply.contacts
	} else {
		fmt.Println("Returning empty contacts")
		done <- []Contact{}
	}
}

/*
func (k *Kademlia) IterativeFindNode(target Contact, delta int) (ret ContactRecords) {
		done := make(chan Contacts)

		ret = make(ContactRecords, BucketSize)

		frontier := new(Contacts)

		seen := make(map[string]bool)

		for _, record := range k.routes.FindClosest(target.id, delta) {
			ret = append(ret, record)
			frontier.Push(record.Contact)
			seen[record.id.String()] = true
		}

		pending := 0
		for i := 0; i < delta && frontier.Len() > 0; i++ {
			pending++
			go k.sendQuery(frontier.Pop().(Contact), target.id, done)
		}

		for pending > 0 {
			contacts := <-done
			pending--
			for _, contact := range contacts {
				if _, ok := seen[contact.id.String()]; ok == false {
					ret = append(ret, ContactRecord{contact, contact.id.Xor(target.id)})
					frontier.Push(contact)
					seen[contact.id.String()] = true
				}
			}

			for pending < delta && frontier.Len() > 0 {
				go k.sendQuery(frontier.Pop().(Contact), target.id, done)
				pending++
			}
		}

		sort.Sort(ret)
		ret = ret[0:BucketSize]

		return
}
*/
