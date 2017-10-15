package main

import (
	"os"
	"log"
	"net"
	"io/ioutil"
	"UniMessCore"
	"UniMessProtocol"
)

var connectAddr string
var protocolChain *UniMessCore.ProtocolChain
var mode string

func main() {
	protocolChain = loadProtocolChain(os.Args[1])
	mode = os.Args[2]
	listenAddr := os.Args[3]
	connectAddr = os.Args[4]

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleConn(conn)
	}
}

func loadProtocolChain(path string) *UniMessCore.ProtocolChain {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	cc := UniMessCore.LoadConfigChain(data)
	return cc.GetProtocolChain()
}

func handleConn(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	clientConn, err := net.Dial("tcp", connectAddr)
	if err != nil {
		panic(err)
	}

	if mode != "encode" && mode != "decode" {
		log.Println("Warning: Unknown mode. Default to `encode`.")
		mode = "encode"
	}

	if mode == "encode" {
		sm := UniMessProtocol.NewStateMachine(protocolChain, conn, clientConn)
		sm.Start()
	} else {
		sm := UniMessProtocol.NewStateMachine(protocolChain, clientConn, conn)
		sm.Start()
	}
}
