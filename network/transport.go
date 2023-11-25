package network

type NetAdddr string

type Transport interface {
	Consume() <-chan RPC
	Connect(Transport) error
	SendMessage(NetAdddr, []byte) error
	Broadcast([]byte) error
	Addr() NetAdddr
}
