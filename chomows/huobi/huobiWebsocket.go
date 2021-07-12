package huobi

import (
	"fmt"
	"net/http"

	"github.com/GeekChomolungma/Everest-Base-Camp/chomows"
	"github.com/GeekChomolungma/Everest-Base-Camp/logging/applogger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func WebsocketHandler(c *gin.Context) {
	go loop(c)
}

func loop(c *gin.Context) {
	// create the chomo ws conn
	var upGrader = &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	wsConn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// channel to read Chomolungma msg
	chomoReadChannel := make(chan []byte, 10)
	stopChannel := make(chan int, 10)
	go readLoop(wsConn, chomoReadChannel, stopChannel)

	//create ws to huobi server
	HuoBiWs := new(chomows.WebSocketClientBase).Init("api-aws.huobi.pro", "/ws")
	HuoBiWs.SetHandler(
		func() {
			go sendLoop(HuoBiWs, chomoReadChannel, stopChannel)
		},
		func(response []byte) {
			// send BinaryMessage resp to Chomolungma
			err = wsConn.WriteMessage(websocket.BinaryMessage, response)
			if err != nil {
				fmt.Println(err.Error())
			}
		})

	HuoBiWs.Connect(true)
}

func readLoop(wsConn *websocket.Conn, rch chan []byte, stopChannel chan int) {
	for {
		_, message, err := wsConn.ReadMessage()
		fmt.Println(message)
		if err != nil {
			applogger.Info("Chomo WebSocket disconnected:", err.Error())
			wsConn.Close()
			stopChannel <- 1
			return
		}
		rch <- message
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
