package huobi

import (
	"net/http"

	"github.com/GeekChomolungma/Everest-Base-Camp/chomows"
	"github.com/GeekChomolungma/Everest-Base-Camp/logging/applogger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func WebsocketHandlerV2(c *gin.Context) {
	// create the chomo ws conn
	var upGrader = &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	wsConn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		applogger.Error("WebSocket V2 UpGrade failed:", err.Error())
		return
	}

	// channel to read Chomolungma msg
	chomoReadChannel := make(chan []byte, 100)
	stopChannel := make(chan int, 10)
	clientID = clientID + 1
	go readLoop(wsConn, chomoReadChannel, stopChannel, clientID)

	//create ws V2 to huobi server
	HuoBiWs := new(chomows.WebSocketClientBase).Init("api-aws.huobi.pro", "/ws/v2")
	HuoBiWs.SetHandler(
		func() {
			go sendLoop(HuoBiWs, chomoReadChannel, stopChannel, clientID)
		},
		func(response []byte, msgType int) {
			// send BinaryMessage resp to Chomolungma
			err = wsConn.WriteMessage(msgType, response)
			if err != nil {
				applogger.Error("HuoBiWs V2 send TextMessage resp to Chomolungma failed:", err.Error())
			}
		})

	HuoBiWs.Connect(true)
}
