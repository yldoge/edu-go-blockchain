package main

import (
	"bytes"
	"math/rand"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/yldoge/edu-go-blockchain/core"
	"github.com/yldoge/edu-go-blockchain/crypto"
	"github.com/yldoge/edu-go-blockchain/network"
)

func main() {
	trl := network.NewLocalTransport("LOCAL")
	trr := network.NewLocalTransport("REMOTE")

	trl.Connect(trr)
	trr.Connect(trl)

	go func() {
		for {
			if err := sendTransaction(trr, trl.Addr()); err != nil {
				log.Error(err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	ops := network.ServerOpts{
		Transports: []network.Transport{trl},
	}

	s := network.NewServer(ops)
	s.Start()
}

// tmp function to test functionaliy
func sendTransaction(tr network.Transport, to network.NetAdddr) error {
	pvk := crypto.GeneratePrivateKey()
	data := []byte(strconv.FormatInt(rand.Int63n(1000), 10))
	tx := core.NewTransaction(data)
	tx.Sign(pvk)
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())

	return tr.SendMessage(to, msg.Bytes())
}
