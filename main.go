package main

import (
	"time"

	"github.com/yldoge/edu-go-blockchain/network"
)

func main() {
	trl := network.NewLocalTransport("LOCAL")
	trr := network.NewLocalTransport("REMOTE")

	trl.Connect(trr)
	trr.Connect(trl)

	go func() {
		for {
			trr.SendMessage(trl.Addr(), []byte("Hello, world!!"))
			time.Sleep(1 * time.Second)
		}
	}()

	ops := network.ServerOpts{
		Transports: []network.Transport{trl},
	}

	s := network.NewServer(ops)
	s.Start()
}
