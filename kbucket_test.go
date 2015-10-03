package kademlia

import (
	"testing"
)

func TestNewKBucket(t *testing.T) {
	kb := NewKBucket()

	if kb.Len() != 0 {
		t.Error("NewKBucket should have 0 length, got", kb.Len())
	}
}

func TestFindContact(t *testing.T) {
	kb := NewKBucket()

	firstContact := randomContact()

	if foundPtr := kb.findContact(*firstContact); foundPtr != nil {
		t.Error("KBucket is empty, findContact should return nil")
	}

	kb.PushBack(firstContact)

	foundPtr := kb.findContact(*firstContact)
	if foundPtr == nil || foundPtr.Value.(*Contact) != firstContact {
		t.Error("findContact should return non-nil Contact equal to firstContact")
	}

	// add contact to beginning of slice
	kb.PushFront(randomContact())

	foundPtr = kb.findContact(*firstContact)
	if foundPtr == nil || foundPtr.Value.(*Contact) != firstContact {
		t.Error("findContact should return non-nil Contact equal to firstContact")
	}
}

func TestFindByID(t *testing.T) {
	kb := NewKBucket()

	firstContact := randomContact()

	if foundPtr := kb.findById(firstContact.ID); foundPtr != nil {
		t.Error("KBucket is empty, findById should return nil")
	}

	kb.PushBack(firstContact)

	foundPtr := kb.findById(firstContact.ID)
	if foundPtr == nil || foundPtr.Value.(*Contact) != firstContact {
		t.Error("findById should return non-nil Contact equal to firstContact")
	}

	// add contact to beginning of slice
	kb.PushFront(randomContact())

	foundPtr = kb.findById(firstContact.ID)
	if foundPtr == nil || foundPtr.Value.(*Contact) != firstContact {
		t.Error("findById should return non-nil Contact equal to firstContact")
	}
}

func TestIsFull(t *testing.T) {
	kb := NewKBucket()

	for i := 0; i < BucketSize; i++ {
		if kb.isFull() {
			t.Error("KBucket should not be full before adding another contact")
		}
		kb.PushBack(randomContact())
	}

	if !kb.isFull() {
		t.Error("KBucket should be full after adding", BucketSize, "contacts")
	}
}

func randomContact() (contact *Contact) {
	contact = new(Contact)
	contactID := NewRandomNodeID()
	*contact = NewContact(contactID, "")
	return
}
