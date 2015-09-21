package main

import (
	"flag"
	"fmt"
	"github.com/cfromknecht/kademlia"
)

func parseFlags() (port *int, firstContact *kademlia.Contact) {
	port = flag.Int("port", 6000, "a int")
	firstID := flag.String("first-id", "", "a hexideicimal node ID")
	firstIP := flag.String("first-ip", "", "the TCP address of an existing node")

	flag.Parse()

	if *firstID == "" || *firstIP == "" {
		firstID = nil
		firstIP = nil
	} else {
		firstContact = new(kademlia.Contact)
		*firstContact = kademlia.NewContact(kademlia.NewNodeID(*firstID), *firstIP)
	}

	return
}

func main() {
	port, firstContact := parseFlags()

	if port == nil {
		panic("Must supply desired port number")
	}

	fmt.Println("Initializing Kademlia DHT ...")

	selfID := kademlia.NewRandomNodeID()

	selfAddress := fmt.Sprintf("127.0.0.1:%d", *port)
	self := kademlia.NewContact(selfID, selfAddress)
	fmt.Println("Self:", selfID, selfAddress)

	selfNetwork := kademlia.NewKademlia(self, "Certcoin-DHT")

	selfNetwork.Serve()

	if firstContact != nil {
		contactRecords := selfNetwork.IterativeFindNode(*firstContact, 5)
		fmt.Println("Contacts: %s", contactRecords)
	}

	done := make(chan bool)
	_ = <-done
}
