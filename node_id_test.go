package kademlia

import (
	"fmt"
	"testing"
)

func TestNewRandomNodeIDLength(t *testing.T) {
	id := NewRandomNodeID()

	if len(id) != IDLength {
		t.Error(fmt.Sprintf("Expected %d", IDLength))
	}
}

func TestNewRandomNodeDifferent(t *testing.T) {
	id1 := NewRandomNodeID()
	id2 := NewRandomNodeID()

	if id1 == id2 {
		t.Error("NewRandomNodeID should generate different Node IDs")
	}
}

func TestNewNodeID(t *testing.T) {
	hexID := "c06349c2f47c837f96d782f2753b2266d548bfa3"
	decodedID := NewNodeID(hexID)
	encodedID := decodedID.String()

	if hexID != encodedID {
		t.Error("IDs should be equal after decoding and re-enconding")
	}
}

func TestString(t *testing.T) {
	id := NewRandomNodeID()
	idString := id.String()
	recoveredID := NewNodeID(idString)

	if id != recoveredID {
		t.Error("IDs should be equal after encoding as string and decoding")
	}
}

func TestEqual(t *testing.T) {
	id := NewRandomNodeID()

	if !id.Equals(id) {
		t.Error("IDs should be equal to themselves")
	}
}

func TestLess(t *testing.T) {
	id := NewRandomNodeID()
	lesserID := id

	// decrement lesserID by 1 in the least significant, non-zero byte
	for i := IDLength - 1; i >= 0; i-- {
		if lesserID[i] != 0 {
			lesserID[i] -= 1
			break
		}
	}

	if !lesserID.Less(id) {
		t.Error("Smaller ID should be less than original ID")
	}
}

func TestXor(t *testing.T) {
	expectedXor := "15d9a75528691e5ebc0a415b99ed3b98f88110bd"
	id1String := "66472dba5cf4e1cbad155ad05beb14cb19d7c65a"
	id2String := "739e8aef749dff95111f1b8bc2062f53e156d6e7"

	id1 := NewNodeID(id1String)
	id2 := NewNodeID(id2String)

	if id1.Xor(id2).String() != expectedXor {
		t.Error(fmt.Sprintf("XOR of %s and %s should be %s", id1, id2, expectedXor))
	}
}

var prefixTests = []struct {
	id1       string
	id2       string
	prefixLen int
}{
	{
		"8000000000000000000000000000000000000000",
		"0000000000000000000000000000000000000000",
		0,
	},
	{
		"4000000000000000000000000000000000000000",
		"0000000000000000000000000000000000000000",
		1,
	},
	{
		"0100000000000000000000000000000000000000",
		"0000000000000000000000000000000000000000",
		7,
	},
	{
		"0080000000000000000000000000000000000000",
		"0000000000000000000000000000000000000000",
		8,
	},
	{
		"0000000000000000000000000000000000000001",
		"0000000000000000000000000000000000000000",
		IDBytesLength - 1,
	},
	{
		"0000000000000000000000000000000000000000",
		"0000000000000000000000000000000000000000",
		-1,
	},
}

func TestPrefixLen(t *testing.T) {
	for _, tt := range prefixTests {
		id1 := NewNodeID(tt.id1)
		id2 := NewNodeID(tt.id2)
		expectedLen := tt.prefixLen

		actualLen := id1.PrefixLen(id2)
		if actualLen != expectedLen {
			t.Error(fmt.Sprintf("Expected prefix length of %s and %s is %d, got %d", id1, id2, expectedLen, actualLen))
		}
		if actualLen != id2.PrefixLen(id1) {
			t.Error("XOR distance metric should be symmetric")
		}
	}
}
