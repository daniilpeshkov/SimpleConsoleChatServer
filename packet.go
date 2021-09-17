package main

import (
	"bufio"
	"io"
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

type NetReader struct {
	reader *bufio.Reader
}

func (r NetReader) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	var i int
	for i = 0; i < len(p); i++ {
		b, err := r.ReadByte()
		if err != nil {
			return i, err
		}
		p[i] = b
	}
	return i + 1, nil
}

func (r NetReader) ReadByte() (byte, error) {
	b, err := r.reader.ReadByte()
	if err != nil {
		return 0, err
	}
	switch b {
	case StopByte:
		return 0, EOP{}
	case MaskByte:
		b, err = r.reader.ReadByte()
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

type NetConverter struct {
	reader *io.Reader
	writer *io.Writer
}

// func (nw NetConverter) ConvertAll() {
// 	for {
// 		b, err := nw.reader.Read()
// 	}
// }

// func (w NetConverter) Write(b []byte) (int, error) {

// }

func ConvertToNet(bytes []byte) []byte {

	net_bytes := make([]byte, 0, len(bytes)*2)

	//startByte is unused for now
	//net_bytes = append(net_bytes, StartByte)

	for _, v := range bytes {
		switch v {
		case 0x7E:
			net_bytes = append(net_bytes, 0x7D, 0x5E)
		case 0x7D:
			net_bytes = append(net_bytes, 0x7D, 0x5D)
		case 0x15:
			net_bytes = append(net_bytes, 0x7D, 0x35)
		default:
			net_bytes = append(net_bytes, v)
		}

	}
	net_bytes = append(net_bytes, 0x15)
	return net_bytes
}

func ConvertFromNet(net_bytes []byte) []byte {
	bytes := make([]byte, 0, len(net_bytes))

	for i := 0; i < len(net_bytes); i++ {
		if net_bytes[i] == 0x7D {
			i += 1
			switch net_bytes[i] {
			case 0x5E:
				bytes = append(bytes, 0x7E)
			case 0x5D:
				bytes = append(bytes, 0x7D)
			case 0x35:
				bytes = append(bytes, 0x15)
			}
		} else if net_bytes[i] != 0x15 && net_bytes[i] != 0x7E {
			bytes = append(bytes, net_bytes[i])
		}
	}
	return bytes
}
