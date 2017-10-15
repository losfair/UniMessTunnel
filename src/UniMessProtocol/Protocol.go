package UniMessProtocol

import (
	"net"
	"log"
	"encoding/binary"
	"UniMessCore"
)

type StateMachine struct {
	chain *UniMessCore.ProtocolChain
	encodedConn net.Conn
	decodedConn net.Conn
}

func NewStateMachine(
	chain *UniMessCore.ProtocolChain,
	encodedConn net.Conn,
	decodedConn net.Conn,
) *StateMachine {
	return &StateMachine {
		chain: chain,
		encodedConn: encodedConn,
		decodedConn: decodedConn,
	}
}

func (sm *StateMachine) Start() {
	// decoded -> encoded
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()

		defer sm.encodedConn.Close()

		for {
			buf := make([]byte, 4094)
			n, err := sm.decodedConn.Read(buf)
			if err != nil {
				break
			}

			data := buf[:n]

			outData := make([]byte, 2)
			binary.LittleEndian.PutUint16(outData[0:], uint16(len(data)))

			outData = append(outData, sm.chain.EncodePacket(data)...)
			_, err = sm.encodedConn.Write(outData)
			if err != nil {
				break
			}
		}
	}()

	// encoded -> decoded
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()

		defer sm.decodedConn.Close()

		var recvBuf []byte = make([]byte, 0)

		feed := func (b byte) bool {
			recvBuf = append(recvBuf, b)
			if len(recvBuf) <= 2 {
				return true
			}

			var expLen uint16 = 0
			expLen = binary.LittleEndian.Uint16(recvBuf[:2])

			if len(recvBuf) == int(expLen) + 2 {
				_, err := sm.decodedConn.Write(sm.chain.DecodePacket(recvBuf[2:]))
				if err != nil {
					return false
				}
				recvBuf = make([]byte, 0)
			}

			return true
		}

		buf := make([]byte, 4096)

		for {
			n, err := sm.encodedConn.Read(buf)
			if err != nil {
				break
			}

			data := buf[:n]
			for _, b := range data {
				ret := feed(b)
				if !ret {
					return;
				}
			}
		}
	}()
}
