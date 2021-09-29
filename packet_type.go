package main

import simpleTcpMessage "github.com/daniilpeshkov/go-simple-tcp-message"

func isLoginRequest(msg *simpleTcpMessage.Message) bool {
	sysBytes, sysOk := msg.GetField(TagSys)
	nameBytes, nameOk := msg.GetField(TagName)
	return sysOk && nameOk && len(sysBytes) == 1 && sysBytes[0] == SysLoginRequest && len(nameBytes) > 0
}

func isMessageRequest(msg *simpleTcpMessage.Message) bool {
	sysBytes, sysOk := msg.GetField(TagSys)
	msgBytes, textOk := msg.GetField(TagMessage)

	return sysOk && textOk && len(sysBytes) == 1 && sysBytes[0] == SysMessage && len(msgBytes) > 0

}
