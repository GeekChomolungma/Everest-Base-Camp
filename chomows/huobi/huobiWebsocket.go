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
	defer wsConn.Close()

	// channel to read Chomolungma msg
	chomoReadChannel := make(chan []byte, 10)
	go readLoop(wsConn, chomoReadChannel)

	//create ws to huobi server
	HuoBiWs := new(chomows.WebSocketClientBase).Init("api-aws.huobi.pro", "/ws")
	HuoBiWs.SetHandler(
		func() {
			go sendLoop(HuoBiWs, chomoReadChannel)
		},
		func(response []byte) {
			// send BinaryMessage resp to Chomolungma
			err = wsConn.WriteMessage(websocket.BinaryMessage, response)
			if err != nil {
				fmt.Println(err.Error())
			}
		})

	HuoBiWs.Connect(true)
	HuoBiWs.Close()
}

func readLoop(wsConn *websocket.Conn, rch chan []byte) {
	for {
		_, message, err := wsConn.ReadMessage()
		fmt.Println(message)
		if err != nil {
			applogger.Info("WebSocket connected", err.Error())
		}
		rch <- message
	}
}

func sendLoop(WebSocketClientBase *chomows.WebSocketClientBase, sch chan []byte) {
	for {
		select {
		case message := <-sch:
			WebSocketClientBase.Send(string(message))
		}
	}
}
