package kademlia

import (
	"crypto/rand"
	"encoding/hex"
)

type NodeID [IDLength]byte

func NewNodeID(data string) (ret NodeID) {
	decoded, _ := hex.DecodeString(data)
	for i := 0; i < IDLength; i++ {
		ret[i] = decoded[i]
	}
	return
}

func NewRandomNodeID() (ret NodeID) {
	buffer := make([]byte, IDLength)
	_, err := rand.Read(buffer)
	check(err)

	for i, b := range buffer {
		ret[i] = b
	}

	return
}

func (node NodeID) String() string {
	return hex.EncodeToString(node[0:IDLength])
}

func (node NodeID) Equals(other NodeID) bool {
	for i := 0; i < IDLength; i++ {
		if node[i] != other[i] {
			return false
		}
	}
	return true
}

func (node NodeID) Less(other interface{}) bool {
	for i := 0; i < IDLength; i++ {
		if node[i] != other.(NodeID)[i] {
			return node[i] < other.(NodeID)[i]
		}
	}

	return false
}

func (node NodeID) Xor(other NodeID) (ret NodeID) {
	for i := 0; i < IDLength; i++ {
		ret[i] = node[i] ^ other[i]
	}
	return
}

func (node NodeID) PrefixLen(other NodeID) (ret int) {
	distance := node.Xor(other)
	for i := 0; i < IDLength; i++ {
		for j := 0; j < 8; j++ {
			if (distance[i]>>uint8(7-j))&0x1 != 0 {
				return 8*i + j
			}
		}
	}
	return -1
}
