package network

type NetAdddr string

type RPC struct {
	From    NetAdddr
	Payload []byte
}

type Transport interface {
	Consume() <-chan RPC
	Connect(Transport) error
	SendMessage(NetAdddr, []byte) error
	Addr() NetAdddr
}
