// Package websocket implements a websocket based transport for go-libp2p.
package websocket

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/transport"
	ma "github.com/multiformats/go-multiaddr"
	mafmt "github.com/multiformats/go-multiaddr-fmt"
	manet "github.com/multiformats/go-multiaddr/net"
)

// WsProtocol is the multiaddr protocol definition for this transport.
//
// Deprecated: use `ma.ProtocolWithCode(ma.P_WS)
var WsProtocol = ma.ProtocolWithCode(ma.P_WS)

// WsFmt is multiaddr formatter for WsProtocol
var WsFmt = mafmt.And(mafmt.TCP, mafmt.Or(
	mafmt.Base(ma.P_WS),
	mafmt.Base(ma.P_WSS),
))

// WsCodec is the multiaddr-net codec definition for the websocket transport
var WsCodec = &manet.NetCodec{
	NetAddrNetworks:  []string{"websocket"},
	ProtocolName:     "ws",
	ConvertMultiaddr: ConvertWebsocketMultiaddrToNetAddr,
	ParseNetAddr:     ParseWebsocketNetAddr,
}
var WssCodec = &manet.NetCodec{
	NetAddrNetworks:  []string{"websocket secure"},
	ProtocolName:     "wss",
	ConvertMultiaddr: ConvertWebsocketMultiaddrToNetAddr,
	ParseNetAddr:     ParseWebsocketNetAddr,
}

// This is _not_ WsFmt because we want the transport to stick to dialing fully
// resolved addresses.
var WsFmtDial = mafmt.And(mafmt.IP, mafmt.Base(ma.P_TCP), mafmt.Or(
	mafmt.Base(ma.P_WS),
	mafmt.Base(ma.P_WSS),
))

func init() {
	manet.RegisterNetCodec(WsCodec)
	manet.RegisterNetCodec(WssCodec)
}
// Option is the type implemented by functional options.
//
// Actual options that one can use vary based on build target environment, i.e.
// the options available on the browser differ from those available natively.
type Option func(cfg *WebsocketConfig) error

var _ transport.Transport = (*WebsocketTransport)(nil)

func (t *WebsocketTransport) CanDial(a ma.Multiaddr) bool {
	return WsFmtDial.Matches(a)
}

func (t *WebsocketTransport) Protocols() []int {
	return []int{ma.P_WS, ma.P_WSS}
}

func (t *WebsocketTransport) Proxy() bool {
	return false
}

func (t *WebsocketTransport) Dial(ctx context.Context, raddr ma.Multiaddr, p peer.ID) (transport.CapableConn, error) {
	macon, err := t.maDial(ctx, raddr)
	if err != nil {
		return nil, err
	}
	return t.Upgrader.UpgradeOutbound(ctx, t, macon, p)
}
