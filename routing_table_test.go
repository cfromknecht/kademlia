package kademlia

import "testing"

func TestNewRoutingTable(t *testing.T) {
	selfID := NewRandomNodeID()
	self := NewContact(selfID, "127.0.0.1:6000")
	table := NewRoutingTable(self)

	if table.Self() != self {
		t.Error("Routing table self Contact not copied properly")
	}

	for i := 0; i < IDBytesLength; i++ {
		if table.kbuckets[i].Len() != 0 {
			t.Error("All buckets should have 0 length after initialization")
		}
	}
}

var selfID = NewRandomNodeID()
var selfContact = NewContact(selfID, "127.0.0.1:6000")
var updateTests = []struct {
	self    Contact
	updates Contacts
}{
	{
		selfContact,
		Contacts{selfContact},
	},
}

func TestUpdate(t *testing.T) {
	for _, tt := range updateTests {
		table := NewRoutingTable(tt.self)

		for _, contact := range tt.updates {
			table.UpdateChan <- contact
		}
	}
}
