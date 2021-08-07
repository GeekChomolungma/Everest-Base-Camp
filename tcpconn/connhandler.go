package tcpconn

import (
	"net"

	"github.com/GeekChomolungma/Everest-Base-Camp/chomows"
	"github.com/GeekChomolungma/Everest-Base-Camp/config"
	"github.com/GeekChomolungma/Everest-Base-Camp/logging/applogger"
)

func TcpServerStart() {
	ln, err := net.Listen("tcp", config.TcpServerSetting.Host)
	if err != nil {
		// handle error
		applogger.Error("TcpServerStart Listen failed:", err.Error())
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			applogger.Error("TcpServerStart Accept failed:", err.Error())
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	// channel to read Chomolungma msg
	chomoReadChannel := make(chan []byte, 10)
	stopChannel := make(chan int, 10)
	go readLoop(conn, chomoReadChannel, stopChannel)

	//create ws to huobi server
	HuoBiWs := new(chomows.WebSocketClientBase).Init("api-aws.huobi.pro", "/ws")
	HuoBiWs.SetHandler(
		func() {
			go sendLoop(HuoBiWs, chomoReadChannel, stopChannel)
		},
		func(response []byte, msgType int) {
			// send BinaryMessage resp to Chomolungma
			_, err := conn.Write(response)
			if err != nil {
				applogger.Error("HuoBiWs send BinaryMessage resp to Chomolungma failed:", err.Error())
			}
		})

	HuoBiWs.Connect(true)
}

func readLoop(conn net.Conn, rch chan []byte, stopChannel chan int) {
	for {
		buf := make([]byte, 1024*8)
		n, err := conn.Read(buf)
		if err != nil {
			applogger.Info("Chomo TCP handleConn error:", err.Error())
			conn.Close()
			stopChannel <- 1
			return
		} else {
			applogger.Info("Chomo TCP handleConn Received bytes length %d, Payload: %s", n, string(buf))
		}

		rch <- buf
	}
}

func sendLoop(WebSocketClientBase *chomows.WebSocketClientBase, sch chan []byte, stopChannel chan int) {
	for {
		select {
		case message := <-sch:
			WebSocketClientBase.Send(string(message))
		case <-stopChannel:
			applogger.Info("Chomo disconnected, close HuoBi conn too.")
			WebSocketClientBase.Close()
			return
		}
	}
}
