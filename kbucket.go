package kademlia

import (
	"container/list"
)

type ContactList *list.List

type KBucket struct {
	*list.List
	UpdateChan         chan Contact
	LookupRequestChan  chan Contact
	LookupResponseChan chan Contact
}

func NewKBucket() (kb *KBucket) {
	kb = new(KBucket)
	*kb = KBucket{
		list.New(),
		make(chan Contact),
		make(chan Contact),
		make(chan Contact),
	}
	go kb.Run()
	return
}

func (kb *KBucket) Run() {
	for {
		select {
		case contact := <-kb.UpdateChan:
			kb.Update(contact)
		case contact := <-kb.LookupRequestChan:
			kb.LookupResponseChan <- kb.Lookup(contact)
		}
	}
}

func (kb *KBucket) Update(contact Contact) {
	foundPtr := kb.findContact(contact)
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
		kb.PushBack(foundPtr)
	}
}

func (kb *KBucket) Lookup(contact Contact) Contact {
	return contact
}

func (kb KBucket) findContact(contact Contact) *list.Element {
	return kb.findById(contact.id)
}

func (kb KBucket) findById(nodeID NodeID) *list.Element {
	for el := kb.Front(); el != nil; el = el.Next() {
		if nodeID == el.Value.(*Contact).id {
			return el
		}
	}

	return nil
}

func (kb KBucket) isFull() bool {
	return kb.Len() >= BucketSize
}
