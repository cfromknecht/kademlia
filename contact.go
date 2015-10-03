package kademlia

/*
 * Contact
 */

type Contact struct {
	ID      NodeID
	Address string
}

func NewContact(node NodeID, address string) Contact {
	return Contact{
		ID:      node,
		Address: address,
	}
}

/*
 * Contacts
 */

type Contacts []Contact

func (h Contacts) Len() int           { return len(h) }
func (h Contacts) Less(i, j int) bool { return h[i].ID.Less(h[j].ID) }
func (h Contacts) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *Contacts) Push(x interface{}) {
	*h = append(*h, x.(Contact))
}

func (h *Contacts) Pop() interface{} {
	oldHeap := *h
	oldLength := len(oldHeap)
	element := oldHeap[oldLength-1]
	*h = oldHeap[0 : oldLength-1]
	return element
}
