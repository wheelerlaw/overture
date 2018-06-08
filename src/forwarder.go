package main

import "encoding/json"

type Forwarder struct {
	Name             string
	Address          string
	//Port             int
	Protocols        []string
	Timeout          int
}


func (forwarder *Forwarder) UnmarshalJSON(b []byte) error {
	type ForwarderAlias Forwarder

	forwarderAlias := new(ForwarderAlias)

	if err := json.Unmarshal(b, &forwarderAlias); err != nil {
		return err
	}

	if len(forwarderAlias.Protocols) == 0 {
		forwarderAlias.Protocols = []string{"tcp", "udp"}
	}

	*forwarder = Forwarder(*forwarderAlias)

	return nil
}