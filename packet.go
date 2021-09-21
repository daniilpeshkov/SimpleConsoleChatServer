package main

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
