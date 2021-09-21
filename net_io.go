package main

import (
	"bufio"
	"io"
	"net"
)

const (
	StartByte = byte(0x7E)
	StopByte  = byte(0x15)
	MaskByte  = byte(0x7D)
)

type EOP struct{} // reader reached packet delimiter (End Of Packet)

func (err EOP) Error() string {
	return "End of packet reached"
}

type NetIO struct {
	writer     io.Writer
	byteReader io.ByteReader
}

func NewNetIO(conn net.Conn) *NetIO {
	return &NetIO{conn, bufio.NewReader(conn)}
}

//writes slice of bytes into net
func (nw NetIO) WriteBytes(p []byte, SendEOP bool) (int, error) {
	for i, b := range p {
		var err error
		var w_slice []byte
		switch b {
		case 0x7E:
			w_slice = []byte{0x7D, 0x5E}
		case 0x7D:
			w_slice = []byte{0x7D, 0x5D}
		case 0x15:
			w_slice = []byte{0x7D, 0x35}
		default:
			w_slice = []byte{b}
		}
		_, err = nw.writer.Write(w_slice)
		if err != nil {
			return i, err
		}
	}
	var err error
	if SendEOP {
		_, err = nw.writer.Write([]byte{StopByte})
	}
	return len(p), err
}

//write bytes from reader (until EOF) into net
func (nio NetIO) ReadFrom(r io.Reader) (n int64, err error) {
	buf := make([]byte, 128)
	w_len := int64(0)
	sendEOP := false
	for sendEOP == false {
		r_len, readErr := r.Read(buf)
		var w_slice []byte
		if r_len > 0 {
			w_slice = buf[:r_len]
		}
		if readErr != nil {
			sendEOP = true
		}
		sent, err := nio.WriteBytes(w_slice, sendEOP)
		w_len += int64(sent)
		if err != nil {
			return w_len, err
		}
	}
	return w_len, nil
}

//read from net
func (nio NetIO) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	var i int
	for i = 0; i < len(p); i++ {
		b, err := nio.readConvertedByte()
		if err != nil {
			return i, err
		}
		p[i] = b
	}
	return i + 1, nil
}

func (nio NetIO) readConvertedByte() (byte, error) {
	b, err := nio.byteReader.ReadByte()
	if err != nil {
		return 0, err
	}
	switch b {
	case StopByte:
		return 0, EOP{}
	case MaskByte:
		b, err = nio.byteReader.ReadByte()
		if err != nil {
			return 0, err
		}
		switch b {
		case 0x5E:
			return 0x7E, nil
		case 0x5D:
			return 0x7D, nil
		case 0x35:
			return 0x15, nil
		default:
			return b, nil // should never be reached (maybe need special error)
		}
	default:
		return b, nil
	}
}
