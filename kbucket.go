package kademlia

import (
	"container/list"
	"fmt"
)

type ContactList *list.List

type KBucket struct {
	*list.List
	UpdateChan chan Contact
}

func NewKBucket() *KBucket {
	kb := &KBucket{
		list.New(),
		make(chan Contact),
	}
	return kb
}

func (kb *KBucket) run() {
	for {
		select {
		case contact := <-kb.UpdateChan:
			kb.update(contact)
		}
	}
}

func (kb *KBucket) update(contact Contact) {
	foundPtr := kb.findContact(contact)
	fmt.Println("Found?", foundPtr)
	if foundPtr != nil {
		// If entry is already in KBucket, move it to back of list
		kb.MoveToBack(foundPtr)
	} else if kb.isFull() {
		// Ping node, and remove if unresponsive
		// TODO(@cfromknecht) Build internal ping

		//kb.Remove(bucket.Front())
		//kb.PushBack(foundPtr)
	} else {
		// KBucket is not full, simply add contact
		fmt.Println("Pushing contact", contact)
		kb.PushBack(contact)
	}
}

func (kb KBucket) findContact(contact Contact) *list.Element {
	return kb.findById(contact.ID)
}

func (kb KBucket) findById(nodeID NodeID) *list.Element {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()

	for el := kb.Front(); el != nil; el = el.Next() {
		fmt.Println("el:", el, "value:", el.Value, "type:", fmt.Sprintf("%T",
			el.Value))
		if nodeID == el.Value.(Contact).ID {
			fmt.Println("Returning el")
			return el
		}
	}

	return nil
}

func (kb KBucket) isFull() bool {
	return kb.Len() >= BucketSize
}
