package test

import (
	"testing"
	wamSerialize "ws/messenger/client/im/wam/serialize"
)

func BenchmarkUsingByteBufferPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buffer := wamSerialize.AcquireEventBuffer()
		buffer.WriteByte(1)

		//encoder := waBinary.AcquireBinaryEncoder()
		//encoder.Marshal(waBinary.Node{
		//	Tag: "abc",
		//	Attrs: waBinary.Attrs{
		//		"xml": "w:stats",
		//	},
		//})

		if buffer.ByteBuffer().Len() != 1 {
			b.Error("not a 1 byte buffer")
		}

		wamSerialize.ReleaseEventBuffer(buffer)
		//waBinary.ReleaseBinaryEncoder(encoder)
	}
}
