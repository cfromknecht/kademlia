package kademlia

const (
	Alpha         = 3
	IDLength      = 20
	IDBytesLength = 8 * IDLength
	BucketSize    = 20
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}
