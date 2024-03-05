package p2p

import (
	"reflect"
	"time"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cometbft/cometbft/p2p/conn"
	"github.com/cometbft/cometbft/proto/tendermint/p2p"
	protop2p "github.com/cometbft/cometbft/proto/tendermint/p2p"
	"github.com/cosmos/gogoproto/proto"
)

var _ Reactor = &PPIReactor{}
var _ proto.Message = &p2p.PPMessage{}

const (
	PPAlive int = iota
	PPinged
	PPonged
)

type peerState struct {
	pingAT int64
	pongAT int64
	id     ID
}

type PPIReactor struct {
	peerList  []ID
	peerState map[ID]peerState
	logger    log.Logger
	quit      chan struct{}
	BaseReactor
}

func NewPPeactor() (PPIReactor, error) {
	r := PPIReactor{
		peerList:  make([]ID, 0),
		peerState: map[ID]peerState{},
		quit:      make(chan struct{}),
	}
	r.BaseReactor = *NewBaseReactor("PP", &r)
	return r, nil
}

// AddPeer implements Reactor.
func (r *PPIReactor) AddPeer(peer Peer) {
	go func() {
		ticker := time.NewTicker(time.Second * 10)
		for {
			select {
			case <-r.quit:
				return
			case <-ticker.C:
				state, ok := r.peerState[peer.ID()]
				if !ok {
					r.logger.Error("peer not exist in peerState", "id", peer.ID())
					return
				}
				r.logger.Error("pp send envelope ping", "id", peer.ID())
				success := peer.SendEnvelope(Envelope{
					Message: &protop2p.PP{
						FromId: string(r.Switch.NetAddress().ID),
						PingAt: uint64(time.Now().Unix()),
						PongAt: 0,
					},
					ChannelID: r.GetChannels()[0].ID,
				})
				if !success {
					r.logger.Error("pp send envelope failed", "id", peer.ID())
					ticker.Reset(time.Second * 5)
				} else {
					state.pingAT = time.Now().Unix()
					r.peerState[peer.ID()] = state
					ticker.Reset(time.Second * 10)
				}
			}
		}
	}()
}

// GetChannels implements Reactor.
func (PPIReactor) GetChannels() []*conn.ChannelDescriptor {
	return []*conn.ChannelDescriptor{
		&conn.ChannelDescriptor{
			ID:                byte(0x43),
			Priority:          6,
			SendQueueCapacity: 100,
			MessageType:       &protop2p.PPMessage{},
		},
	}
}

// InitPeer implements Reactor.
func (r *PPIReactor) InitPeer(peer Peer) Peer {
	r.peerList = append(r.peerList, peer.ID())
	r.peerState[peer.ID()] = peerState{
		pingAT: 0,
		id:     peer.ID(),
	}
	r.logger.Error("init peer id", "id", peer.ID())
	println("init peer id", peer.ID())
	return peer
}

// OnReset implements Reactor.
func (r *PPIReactor) OnReset() error {
	return nil
}

// OnStart implements Reactor.
func (r *PPIReactor) OnStart() error {
	r.logger.Error("pp custom start")
	println("pp custome start")
	return nil
}

// OnStop implements Reactor.
func (PPIReactor) OnStop() {}

// ReceiveEnvelope implements Reactor.
func (r *PPIReactor) ReceiveEnvelope(e Envelope) {
	r.logger.Error("receive envelope")
	ppMsg, ok := e.Message.(*protop2p.PP)
	if !ok {
		r.logger.Error("ppreactor receive message type missmatch actual", "type", reflect.TypeOf(e.Message).String())
	}
	if ppMsg.PongAt == 0 {
		// receive ping
		r.logger.Error("send envelope pong")
		e.Src.SendEnvelope(
			Envelope{
				Message: &protop2p.PP{
					FromId: string(r.Switch.NetAddress().ID),
					PingAt: ppMsg.PingAt,
					PongAt: uint64(time.Now().Unix()),
				},
				ChannelID: r.GetChannels()[0].ID,
			},
		)
	} else {
		// receive pong
		state, ok := r.peerState[e.Src.ID()]
		if !ok {
			return
		}
		state.pongAT = time.Now().Unix()
		r.peerState[e.Src.ID()] = state
	}
}

// RemovePeer implements Reactor.
func (r *PPIReactor) RemovePeer(peer Peer, reason interface{}) {
	r.logger.Error("remove peer", "id", peer.ID())
	delete(r.peerState, peer.ID())
}

// Reset implements Reactor.
func (r *PPIReactor) Reset() error {
	r.peerList = make([]ID, 0)
	r.peerState = map[ID]peerState{}
	r.quit = make(chan struct{})
	return nil
}

// SetLogger implements Reactor.
func (r *PPIReactor) SetLogger(l log.Logger) {
	r.logger = l
}

// Start implements Reactor.
func (PPIReactor) Start() error {
	return nil
}

// Stop implements Reactor.
func (r *PPIReactor) Stop() error {
	close(r.quit)
	return nil
}

// String implements Reactor.
func (r *PPIReactor) String() string {
	return "p2p custom reactor"
}
