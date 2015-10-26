package kademlia

func (k *Kademlia) Bootstrap(target, self Contact) ([]Contact, error) {
	client, err := dialContact(target)
	if err != nil {
		return nil, err
	}

	req := k.NewFindNodeRequest(self.ID)
	res := FindNodeResponse{}

	err = client.Call("KademliaCore.FindNodeRPC", &req, &res)
	if err != nil {
		return nil, err
	}

	k.routes.Update(res.Sender)

	return res.Contacts, nil
}
