package connection

func builtinDecodeLen(frame []byte) int {
	bufferLen := len(frame)

	if bufferLen < 3 {
		return 0
	}

	packetLength := (int(frame[0]) << 16) + (int(frame[1]) << 8) + int(frame[2])

	return packetLength
}
