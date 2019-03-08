package helper

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/libp2p/go-conn-security-multistream"
	ic "github.com/libp2p/go-libp2p-crypto"
	"github.com/libp2p/go-libp2p-peer"
	"github.com/libp2p/go-libp2p-peerstore"
	"github.com/libp2p/go-libp2p-peerstore/pstoremem"
	"github.com/libp2p/go-libp2p-secio"
	"github.com/libp2p/go-libp2p-swarm"
	tptu "github.com/libp2p/go-libp2p-transport-upgrader"
	"github.com/libp2p/go-tcp-transport"
	"github.com/whyrusleeping/go-smux-multistream"
	"github.com/whyrusleeping/go-smux-yamux"
)

func GenSwarmByKey(key ic.PrivKey) (*swarm.Swarm, *tptu.Upgrader) {
	ctx := context.Background()
	priv, pub := key, key.GetPublic()
	pid, err := peer.IDFromPublicKey(pub)
	if err != nil {
		panic(err)
	}
	ps := pstoremem.NewPeerstore()
	ps.AddPubKey(pid, pub)
	ps.AddPrivKey(pid, priv)
	s := swarm.NewSwarm(ctx, pid, ps, nil)

	//NewTCPTransport
	u := GenUpgrader(s)
	tcpTransport := tcp.NewTCPTransport(GenUpgrader(s))

	if err := s.AddTransport(tcpTransport); err != nil {
		fmt.Println(err)
	}
	// TODO 多地址
	//maddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port))
	//s.Listen(maddr)
	//s.AddListenAddr(maddr)
	s.Peerstore().AddAddrs(pid, s.ListenAddresses(), peerstore.PermanentAddrTTL)
	return s, u
}
func GenSwarm() *swarm.Swarm {
	s, _ := GenSwarm2()
	return s
}

func GenSwarm2() (*swarm.Swarm, *tptu.Upgrader) {
	priv, _, err := ic.GenerateKeyPairWithReader(ic.RSA, 2048, rand.Reader)
	if err != nil {
		panic(err) // oh no!
	}
	return GenSwarmByKey(priv)
}

// GenUpgrader creates a new connection upgrader for use with this swarm.
func GenUpgrader(n *swarm.Swarm) *tptu.Upgrader {
	id := n.LocalPeer()
	pk := n.Peerstore().PrivKey(id)
	secMuxer := new(csms.SSMuxer)
	secMuxer.AddTransport(secio.ID, &secio.Transport{
		LocalID:    id,
		PrivateKey: pk,
	})
	multistream.NewBlankTransport()
	stMuxer := multistream.NewBlankTransport()
	stMuxer.AddTransport("/yamux/1.0.0", sm_yamux.DefaultTransport)
	return &tptu.Upgrader{
		Secure:  secMuxer,
		Muxer:   stMuxer,
		Filters: n.Filters,
	}
}
