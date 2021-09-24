package main

import (
	"fmt"
)

const (
	TypeText     = 0
	TypeFileName = 1
	TypeFile     = 2
	TypeDate     = 3
	TypeName     = 4
)

var supportedTypes = map[byte]struct{}{
	TypeText:     {},
	TypeFileName: {},
	TypeFile:     {},
	TypeDate:     {},
	TypeName:     {},
}

func typeSupported(t byte) bool {
	_, ok := supportedTypes[t]
	return ok
}

type message struct {
	fields map[byte][]byte
}

func (msg message) String() string {
	return fmt.Sprintf("Message:\nfrom: %s\ntext: %s\n", msg.fields[TypeName], msg.fields[TypeText])
}

func (msg *message) appendField(_type byte, p []byte) {
	if typeSupported(_type) && len(p) > 0 {
		if _, ok := msg.fields[_type]; !ok {
			msg.fields[_type] = make([]byte, 0, len(p))
		}
		msg.fields[_type] = append(msg.fields[_type], p...)
	}

}

func NewMessage() *message {
	return &message{
		fields: make(map[byte][]byte),
	}
}

// func (msg Message) HasNameFrom() bool {
// 	return len(msg.nameFrom) != 0
// }
// func (msg Message) HasText() bool {
// 	return len(msg.text) != 0
// }
// func (msg Message) HasFileName() bool {
// 	return len(msg.fileName) != 0
// }
// func (msg Message) HasFile() bool {
// 	return len(msg.file) != 0
// }
// func (msg Message) HasTime() bool {
// 	return len(msg.time) != 0
// }
