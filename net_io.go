package main

import (
	"bufio"
	"io"
	"net"
)

const TypeMask = byte(0x1F)
const HasNext = 7
const HeaderLen = 2
const MaxDataLen = 255
const MaxMessageLen = MaxDataLen + HeaderLen

type EOP struct{} // reader reached packet delimiter (End Of Packet)

func (err EOP) Error() string {
	return "End of packet reached"
}

type NetIO struct {
	w io.Writer
	r *bufio.Reader
}

func NewNetIO(conn net.Conn) *NetIO {
	return &NetIO{conn, bufio.NewReader(conn)}
}

func (nio NetIO) RecieveMessage() (*message, error) {
	msg := NewMessage()
	header := [2]byte{}
	hasNext := true
	var dlc uint
	var _type byte

	for hasNext {
		_, err := io.ReadFull(nio.r, header[:])
		if err != nil {
			return nil, err
		}

		hasNext = header[0]&(1<<HasNext) != 0
		dlc = uint(header[1])
		_type = header[0] & TypeMask
		buf := make([]byte, dlc)
		_, err = io.ReadFull(nio.r, buf)
		if err != nil {
			return nil, err
		}
		msg.appendField(_type, buf)
	}
	return msg, nil
}

func (nio NetIO) SendMessage(msg *message) error {
	fieldsLeft := len(msg.fields)
	buf := [MaxMessageLen]byte{}
	for k, v := range msg.fields {
		bLeft := len(v)

		buf[0] = byte(k) & TypeMask

		for bLeft > 0 {
			if fieldsLeft > 1 || bLeft > MaxDataLen {
				buf[0] |= (1 << HasNext)
			} else {
				buf[0] &= (1 << HasNext) ^ 0xFF
			}
			var bytesToSend int
			if bLeft > MaxDataLen {
				bytesToSend = MaxDataLen
			} else {
				bytesToSend = bLeft
			}
			buf[1] = byte(bytesToSend)
			copy(buf[2:], v[:bytesToSend])
			v = v[bytesToSend:]
			bLeft -= bytesToSend
			fieldsLeft -= 1

			_, err := nio.w.Write(buf[:HeaderLen+bytesToSend])
			if err != nil {
				return err
			}
		}
	}

	return nil
}
