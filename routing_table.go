package kademlia

import ()

type RoutingTable struct {
	self              Contact
	kbuckets          [IDBytesLength]*KBucket
	UpdateChan        chan Contact
	LookupRequestChan chan chan Contacts
}

func (rt RoutingTable) Self() Contact {
	return rt.self
}

func NewRoutingTable(self Contact) *RoutingTable {
	rt := &RoutingTable{
		self:              self,
		kbuckets:          [IDBytesLength]*KBucket{},
		UpdateChan:        make(chan Contact),
		LookupRequestChan: make(chan chan Contacts),
	}

	for i := 0; i < IDBytesLength; i++ {
		rt.kbuckets[i] = NewKBucket()
	}

	go rt.run()

	return rt
}

func (rt *RoutingTable) run() {
	for {
		select {
		case contact := <-rt.UpdateChan:
			prefixLength := contact.ID.PrefixLen(rt.self.ID)
			// Check if same ID as self
			if prefixLength == -1 {
				break
			}

			bucket := rt.kbuckets[prefixLength]
			bucket.update(contact)
		case contactsReturnChan := <-rt.LookupRequestChan:
			targetSlice := <-contactsReturnChan
			contacts := rt.findClosest(targetSlice[0].ID)
			contactsReturnChan <- contacts
		}
	}
}

func (rt *RoutingTable) findClosest(nodeID NodeID) Contacts {
	prefixLength := nodeID.PrefixLen(rt.self.ID)
	if prefixLength == -1 {
		prefixLength = IDBytesLength - 1
	}

	contacts := Contacts{}
	for i := prefixLength; i >= 0; i-- {
		elt := rt.kbuckets[i].Front()
		for elt != nil {
			if len(contacts) < BucketSize {
				contacts = append(contacts, elt.Value.(Contact))
			} else {
				return contacts
			}
			elt = elt.Next()
		}
	}

	for i := prefixLength + 1; i < IDBytesLength; i++ {
		elt := rt.kbuckets[i].Front()
		for elt != nil {
			if len(contacts) < BucketSize {
				contacts = append(contacts, elt.Value.(Contact))
				elt = elt.Next()
			} else {
				return contacts
			}
		}
	}

	return contacts
}
