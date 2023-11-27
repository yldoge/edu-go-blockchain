package main

import (
	"bytes"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/yldoge/edu-go-blockchain/core"
	"github.com/yldoge/edu-go-blockchain/crypto"
	"github.com/yldoge/edu-go-blockchain/network"
)

func main() {
	trLocal := network.NewLocalTransport("LOCAL")
	trRemote0 := network.NewLocalTransport("REMOTE_0")
	trRemote1 := network.NewLocalTransport("REMOTE_1")
	trRemote2 := network.NewLocalTransport("REMOTE_2")

	trLocal.Connect(trRemote0)
	trRemote0.Connect(trLocal)

	trRemote0.Connect(trRemote1)
	trRemote1.Connect(trRemote2)

	initRemoteServers([]network.Transport{trRemote0, trRemote1, trRemote2})

	go func() {
		for {
			if err := sendTransaction(trRemote0, trLocal.Addr()); err != nil {
				log.Error(err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	pvk := crypto.GeneratePrivateKey()
	localServer := makeServer("LOCAL", trLocal, &pvk)
	localServer.Start()
}

func initRemoteServers(trs []network.Transport) {
	for i, tr := range trs {
		s := makeServer(
			fmt.Sprintf("REMOTE_%d", i),
			tr, nil,
		)
		go s.Start()
	}
}

func makeServer(id string, tr network.Transport, pk *crypto.PrivateKey) *network.Server {
	opts := network.ServerOpts{
		PrivateKey: pk,
		ID:         id,
		Transports: []network.Transport{tr},
	}
	s, err := network.NewServer(opts)
	if err != nil {
		panic(err)
	}
	return s
}

// tmp function to test functionaliy
func sendTransaction(tr network.Transport, to network.NetAdddr) error {
	pvk := crypto.GeneratePrivateKey()
	data := []byte{0x03, 0x0a, 0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0c, 0x0d, 0x05, 0x0a, 0x0f}
	tx := core.NewTransaction(data)
	tx.Sign(pvk)
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())

	return tr.SendMessage(to, msg.Bytes())
}
