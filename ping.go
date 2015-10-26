package kademlia

type PingRequest struct {
	RPCHeader
}

func (k *Kademlia) NewPingRequest() PingRequest {
	return PingRequest{
		RPCHeader{
			Sender:    k.routes.self,
			NetworkID: k.NetworkID,
		},
	}
}

type PingResponse struct {
	RPCHeader
}

func (k *Kademlia) Ping(target Contact) error {
	client, err := dialContact(target)
	if err != nil {
		return err
	}

	req := k.NewPingRequest()
	res := PingResponse{}

	return client.Call("KademliaCore.PingRPC", &req, &res)
}

func (kc *KademliaCore) PingRPC(req PingRequest, res *PingResponse) error {
	return kc.kad.HandleRPC(req.RPCHeader, &res.RPCHeader)
}
