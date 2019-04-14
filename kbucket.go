package kademlia

import (
	"container/list"
	"fmt"
)

type ContactList *list.List

type KBucket struct {
	*list.List
}

func NewKBucket() *KBucket {
	return &KBucket{
		list.New(),
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
		if nodeID == el.Value.(*Contact).ID {
			return el
		}
	}

	return nil
}

func (kb KBucket) isFull() bool {
	return kb.Len() >= BucketSize
}
