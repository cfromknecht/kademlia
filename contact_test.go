package kademlia_test

import (
	"fmt"
	"github.com/cfromknecht/kademlia"
	"testing"
)

func TestNewContact(t *testing.T) {
	expectedNodeID := kademlia.NewRandomNodeID()
	expectedAddress := "127.0.0.1:1234"

	contact := kademlia.NewContact(expectedNodeID, expectedAddress)

	if expectedNodeID != contact.ID() {
		t.Error(fmt.Sprintf("Expected ID %s, got %s", expectedNodeID, contact.ID()))
	}
	if expectedAddress != contact.Address() {
		t.Error(fmt.Sprintf("Expected address %s, got %s", expectedAddress, contact.Address()))
	}
}

func TestContactsLen(t *testing.T) {
	expectedLen := 20
	contacts := kademlia.Contacts{}

	for i := 0; i < expectedLen; i++ {
		if contacts.Len() != i {
			t.Error(fmt.Sprintf("Contacts length should be %d, got %d", i, contacts.Len()))
		}

		id := kademlia.NewRandomNodeID()
		address := fmt.Sprintf("127.0.0.1:60%d", i)
		contacts = append(contacts, kademlia.NewContact(id, address))
	}

	if contacts.Len() != expectedLen {
		t.Error(fmt.Sprintf("Contacts length should be %d, got %d", expectedLen, contacts.Len()))
	}
}

var contactsLessTests = []struct {
	contact1          kademlia.Contact
	contact2          kademlia.Contact
	lessThanForwards  bool
	lessThanBackwards bool
}{
	{
		kademlia.NewContact(kademlia.NewNodeID("0000000000000000000000000000000000000000"), ""),
		kademlia.NewContact(kademlia.NewNodeID("0000000000000000000000000000000000000001"), ""),
		true,
		false,
	},
	{
		kademlia.NewContact(kademlia.NewNodeID("000000000000000000000000000000000FFFFFFF"), ""),
		kademlia.NewContact(kademlia.NewNodeID("0000000000000000000000000000000000000000"), ""),
		false,
		true,
	},
	{
		kademlia.NewContact(kademlia.NewNodeID("0000000000000000000000000000000000000000"), ""),
		kademlia.NewContact(kademlia.NewNodeID("0000000000000000000000000000000000000000"), ""),
		false,
		false,
	},
}

func TestContactsLess(t *testing.T) {
	for _, tt := range contactsLessTests {
		contacts := kademlia.Contacts{tt.contact1, tt.contact2}

		if tt.lessThanForwards != contacts.Less(0, 1) {
			t.Error("Contacts Less not ordering based on NodeID")
		}

		if tt.lessThanBackwards != contacts.Less(1, 0) {
			t.Error("Contacts Less not ordering based on NodeID")
		}
	}
}

func TestContactsSwap(t *testing.T) {
	contact0 := kademlia.NewContact(kademlia.NewRandomNodeID(), "")
	contact1 := kademlia.NewContact(kademlia.NewRandomNodeID(), "")
	contacts := kademlia.Contacts{contact0, contact1}

	// swap once
	contacts.Swap(0, 1)
	if contact1 != contacts[0] || contact0 != contacts[1] {
		t.Error("Contacts were not swapped properly")
	}

	// return to original ordering
	contacts.Swap(0, 1)
	if contact0 != contacts[0] || contact1 != contacts[1] {
		t.Error("Contacts were not swapped properly")
	}
}

func TestContactsPush(t *testing.T) {
	contactToAdd := kademlia.NewContact(kademlia.NewRandomNodeID(), "")
	contacts := kademlia.Contacts{}

	contacts.Push(contactToAdd)

	if contacts.Len() != 1 {
		t.Error("Contact was not pushed to Contacts array")
	}
	if contacts[0].ID() != contactToAdd.ID() || contacts[0].Address() != contactToAdd.Address() {
		t.Error("Contact not copied correctly during Push")
	}
}

func TestContactsPop(t *testing.T) {
	contactToRemove := kademlia.NewContact(kademlia.NewRandomNodeID(), "")
	contacts := kademlia.Contacts{contactToRemove}

	removedContact := contacts.Pop().(kademlia.Contact)

	if contacts.Len() != 0 {
		t.Error("Contact was not popped from Contacts array")
	}
	if contactToRemove.ID() != removedContact.ID() || contactToRemove.Address() != removedContact.Address() {
		t.Error("Contact was not copied correctly during Pop")
	}
}
