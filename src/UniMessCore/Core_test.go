package UniMessCore

import (
	"testing"
	"crypto/rand"
)

func BenchmarkEncodeDecode(b *testing.B) {
	cc := GenerateConfigChain(16)
	pc := cc.GetProtocolChain()
	data := make([]byte, 4096)

	_, err := rand.Read(data)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		encData := pc.EncodePacket(data)
		pc.DecodePacket(encData)
	}
}
