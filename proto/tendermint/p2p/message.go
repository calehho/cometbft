package p2p

import (
	"fmt"

	"github.com/cosmos/gogoproto/proto"
)

func (m *PP) Wrap() proto.Message {
	mm := &PPMessage{}
	mm.Sum = &PPMessage_Pp{
		Pp: m,
	}
	return mm
}

func (m *PPMessage) Unwrap() (proto.Message, error) {
	switch msg := m.Sum.(type) {
	case *PPMessage_Pp:
		return m.GetPp(), nil
	default:
		return nil, fmt.Errorf("unknown message: %T", msg)
	}
}

// // Wrap implements the p2p Wrapper interface and wraps a mempool message.
// func (m *Txs) Wrap() proto.Message {
// 	mm := &Message{}
// 	mm.Sum = &Message_Txs{Txs: m}
// 	return mm
// }

// // Unwrap implements the p2p Wrapper interface and unwraps a wrapped mempool
// // message.
// func (m *Message) Unwrap() (proto.Message, error) {
// 	switch msg := m.Sum.(type) {
// 	case *Message_Txs:
// 		return m.GetTxs(), nil

// 	default:
// 		return nil, fmt.Errorf("unknown message: %T", msg)
// 	}
// }
