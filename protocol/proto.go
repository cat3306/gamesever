package protocol

import (
	"encoding/binary"
	"errors"
	"github.com/panjf2000/gnet/v2"
)

//
// * 0           2                       6
// * +-----------+-----------------------+
// * |   magic   |       body len        |
// * +-----------+-----------+-----------+
// * |                                   |
// * +                                   +
// * |           body bytes              |
// * +                                   +
// * |            ... ...                |
// * +-----------------------------------+

const (
	payloadLen  = uint32(4)
	protocolLen = uint32(4)
	codeTypeLen = uint32(2)
)

var (
	packetEndian        = binary.LittleEndian
	ErrIncompletePacket = errors.New("incomplete packet")
)

type Context struct {
	Payload  []byte
	CodeType CodeType
	Proto    uint32
	Conn     gnet.Conn
}

func Decode(c gnet.Conn) (*Context, error) {
	bodyOffset := int(payloadLen + protocolLen + codeTypeLen)
	buf, err := c.Peek(bodyOffset)
	if err != nil {
		return nil, err
	}
	bodyLen := packetEndian.Uint32(buf[:payloadLen])
	protocol := packetEndian.Uint32(buf[payloadLen : payloadLen+protocolLen])
	codeType := packetEndian.Uint16(buf[payloadLen+protocolLen : payloadLen+protocolLen+codeTypeLen])
	msgLen := bodyOffset + int(bodyLen)
	if c.InboundBuffered() < msgLen {
		return nil, ErrIncompletePacket
	}
	buf, err = c.Peek(msgLen)
	if err != nil {
		return nil, err
	}
	_, err = c.Discard(msgLen)
	if err != nil {
		return nil, err
	}
	packet := &Context{
		Payload:  buf[bodyOffset:msgLen],
		CodeType: CodeType(codeType),
		Proto:    protocol,
		Conn:     c,
	}
	return packet, nil
}
func Encode(v interface{}, codeType CodeType, proto uint32) []byte {
	if v == nil {
		panic("v nil")
	}
	raw, err := GameCoder(codeType).Marshal(v)
	if err != nil {
		panic(err)
	}
	bodyOffset := int(payloadLen + protocolLen + codeTypeLen)
	msgLen := bodyOffset + len(raw)
	data := make([]byte, msgLen)
	packetEndian.PutUint32(data, uint32(len(raw)))
	packetEndian.PutUint32(data[payloadLen:], proto)
	packetEndian.PutUint16(data[payloadLen+protocolLen:], uint16(codeType))
	copy(data[bodyOffset:msgLen], raw)
	return data
}
