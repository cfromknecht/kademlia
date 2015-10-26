package kademlia

type RoutingTable struct {
	self     Contact
	kbuckets [IDBytesLength]*KBucket
}

func (rt RoutingTable) Self() Contact {
	return rt.self
}

func NewRoutingTable(self Contact) *RoutingTable {
	rt := &RoutingTable{
		self:     self,
		kbuckets: [IDBytesLength]*KBucket{},
	}

	for i := 0; i < IDBytesLength; i++ {
		rt.kbuckets[i] = NewKBucket()
	}

	return rt
}

func (rt *RoutingTable) Update(contact Contact) {
	prefixLength := contact.ID.PrefixLen(rt.self.ID)
	if prefixLength == -1 {
		return
	}

	bucket := rt.kbuckets[prefixLength]
	bucket.Update(contact)
}

func (rt *RoutingTable) FindClosest(target NodeID, delta int) Contacts {
	prefixLength := target.PrefixLen(rt.self.ID)

	contacts := Contacts{}
	if prefixLength == -1 {
		contacts = append(contacts, rt.self)
		prefixLength = IDBytesLength - 1
	}

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

	if contacts.Len() > delta {
		return contacts[:delta]
	}

	return contacts
}
