package huobi

import (
	"net/http"

	"github.com/GeekChomolungma/Everest-Base-Camp/chomows"
	"github.com/GeekChomolungma/Everest-Base-Camp/logging/applogger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func WebsocketHandler(c *gin.Context) {
	// create the chomo ws conn
	var upGrader = &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	wsConn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		applogger.Error("WebSocket UpGrade failed:", err.Error())
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
				applogger.Error("HuoBiWs send BinaryMessage resp to Chomolungma failed:", err.Error())
			}
		})

	HuoBiWs.Connect(true)
}

func readLoop(wsConn *websocket.Conn, rch chan []byte, stopChannel chan int) {
	applogger.Info("Chomo readLoop started.")
	for {
		applogger.Debug("Chomo readLoop goting to read...")
		_, message, err := wsConn.ReadMessage()
		if err != nil {
			applogger.Error("Chomo WebSocket disconnected: %s", err.Error())
			wsConn.Close()
			stopChannel <- 1
			return
		}
		applogger.Debug("Chomo readLoop recieved msg:%s", string(message))
		rch <- message
	}
}

func sendLoop(WebSocketClientBase *chomows.WebSocketClientBase, sch chan []byte, stopChannel chan int) {
	for {
		select {
		case message := <-sch:
			WebSocketClientBase.Send(string(message))
			// TODO: if send err, should close the chomolungma client conn

		case <-stopChannel:
			applogger.Info("Chomo disconnected, close HuoBi conn too.")
			WebSocketClientBase.Close()
			return
		}
	}
}
