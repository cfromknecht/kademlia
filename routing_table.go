package kademlia

type RoutingTable struct {
	self       Contact
	kbuckets   [IDBytesLength]*KBucket
	UpdateChan chan Contact
}

func (rt RoutingTable) Self() Contact {
	return rt.self
}

func NewRoutingTable(self Contact) (rt *RoutingTable) {
	rt = new(RoutingTable)

	rt.self = self

	for i := 0; i < IDBytesLength; i++ {
		rt.kbuckets[i] = NewKBucket()
	}

	rt.UpdateChan = make(chan Contact)
	go rt.Run()

	return
}

func (rt *RoutingTable) Run() {
	for {
		select {
		case contact := <-rt.UpdateChan:
			prefixLength := contact.id.PrefixLen(rt.self.id)
			// Check if same ID as self
			if prefixLength == -1 {
				break
			}

			bucket := rt.kbuckets[prefixLength]
			bucket.UpdateChan <- contact
		}
	}
}
