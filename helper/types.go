package helper

import (
	"context"
	"fmt"
	"reflect"

	"github.com/alecthomas/log4go"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-circuit"
	"github.com/libp2p/go-libp2p-host"
	"github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p-peer"
	"github.com/libp2p/go-libp2p-peerstore"
	"github.com/libp2p/go-libp2p/p2p/host/basic"

	ic "github.com/libp2p/go-libp2p-crypto"
	opts "github.com/libp2p/go-libp2p-kad-dht/opts"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
)

type blankValidator struct{}

func (blankValidator) Validate(_ string, _ []byte) error        { return nil }
func (blankValidator) Select(_ string, _ [][]byte) (int, error) { return 0, nil }

type Node struct {
	Host    host.Host
	Routing *dht.IpfsDHT
}

func NewLocalNode() *Node {
	h, _ := basichost.NewHost(context.Background(), GenSwarm(), &basichost.HostOpts{})
	d, _ := dht.New(context.Background(), h, opts.NamespacedValidator("cc14514", blankValidator{}))
	return &Node{h, d}
}

func libp2popts(key ic.PrivKey, port int) []libp2p.Option {
	var libp2pOpts []libp2p.Option

	// set pk
	libp2pOpts = append(libp2pOpts, libp2p.Identity(key))

	// TODO 需要探测整个网络能监听的所有 ip
	// listen address
	libp2pOpts = append(libp2pOpts, libp2p.ListenAddrStrings(
		fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port),
		fmt.Sprintf("/ip6/::/tcp/%d", port),
	))

	// 支持 relay
	libp2pOpts = append(libp2pOpts, libp2p.EnableRelay(relay.OptHop))

	return libp2pOpts
}

func NewNode(key ic.PrivKey, port int) *Node {
	ctx := context.Background()
	host, err := libp2p.New(ctx, libp2popts(key, port)...)
	if err != nil {
		panic(err)
	}
	d, err := dht.New(ctx, host, opts.NamespacedValidator("cc14514", blankValidator{}))
	if err != nil {
		panic(err)
	}
	return &Node{host, d}
}

func (self *Node) Close() {
	self.Routing.Close()
	self.Host.Close()
}

func (self *Node) Bootstrap(ctx context.Context) error {
	return self.Routing.Bootstrap(ctx)
}

func (self *Node) Connect(ctx context.Context, targetID interface{}, targetAddrs []ma.Multiaddr) error {
	var tid peer.ID
	if "string" == reflect.TypeOf(targetID).Name() {
		stid := targetID.(string)
		tid, _ = peer.IDFromString(stid)
	} else {
		tid = targetID.(peer.ID)
	}
	a := self.Host
	a.Peerstore().AddAddrs(tid, targetAddrs, peerstore.TempAddrTTL)
	pi := peerstore.PeerInfo{ID: tid}
	if err := a.Connect(ctx, pi); err != nil {
		return err
	}
	return nil
}

func (self *Node) PutValue(ctx context.Context, key string, value []byte) error {
	return self.Routing.PutValue(ctx, key, value)
}

func (self *Node) GetValue(ctx context.Context, key string) ([]byte, error) {
	return self.Routing.GetValue(ctx, key)
}

func (self *Node) FindPeer(ctx context.Context, targetID interface{}) (pstore.PeerInfo, error) {
	var (
		tid peer.ID
		err error
	)
	if "string" == reflect.TypeOf(targetID).Name() {
		stid := targetID.(string)
		tid, err = peer.IDB58Decode(stid)
		if err != nil {
			log4go.Error(err)
			return pstore.PeerInfo{}, err
		}
	} else {
		tid = targetID.(peer.ID)
	}
	return self.Routing.FindPeer(ctx, tid)
}
