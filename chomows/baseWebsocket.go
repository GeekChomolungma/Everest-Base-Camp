package chomows

import (
	"fmt"
	"net/http"

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
	for {
		// read
		mt, message, err := wsConn.ReadMessage()
		fmt.Println(string(message))
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		// write
		responseMsg := "hello, websocket!"
		err = wsConn.WriteMessage(mt, []byte(responseMsg))
		if err != nil {
			fmt.Println(err.Error())
			break
		}
	}
}
