package huobi

import (
	"net/http"
	"time"

	"github.com/GeekChomolungma/Everest-Base-Camp/chomows"
	"github.com/GeekChomolungma/Everest-Base-Camp/logging/applogger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var clientID int = 0

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
	chomoReadChannel := make(chan []byte, 100)
	stopChannel := make(chan int, 10)
	clientID = clientID + 1
	go readLoop(wsConn, chomoReadChannel, stopChannel, clientID)

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

func readLoop(wsConn *websocket.Conn, rch chan []byte, stopChannel chan int, clientID int) {
	applogger.Info("Chomo client-%d readLoop started.", clientID)
	for {
		applogger.Debug("Chomo client-%d readLoop goting to read...", clientID)
		_, message, err := wsConn.ReadMessage()
		if err != nil {
			applogger.Error("Chomo client-%d WebSocket disconnected: %s", clientID, err.Error())
			wsConn.Close()
			stopChannel <- 1
			return
		}
		applogger.Debug("Chomo client-%d readLoop recieved msg:%s", clientID, string(message))
		rch <- message
	}
}

func sendLoop(WebSocketClientBase *chomows.WebSocketClientBase, sch chan []byte, stopChannel chan int) {
	for {
		select {
		case message := <-sch:
			// frequency limit
			time.Sleep(time.Duration(100) * time.Millisecond)
			WebSocketClientBase.Send(string(message))
			// TODO: if send err, should close the chomolungma client conn

		case <-stopChannel:
			applogger.Info("Chomo disconnected, close HuoBi conn too.")
			WebSocketClientBase.Close()
			return
		}
	}
}
