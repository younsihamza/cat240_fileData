package sender

import (
	"os"
	"cat240/global"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func addClient(conn *websocket.Conn) {
	global.MuClient.Lock()
	defer global.MuClient.Unlock()
	global.Clients[conn] = true
}

func removeClient(conn *websocket.Conn) {
	global.MuClient.Lock()
	defer global.MuClient.Unlock()
	conn.Close()
	delete(global.Clients, conn)
}

func broadcastMessage(message []byte) {
	global.MuClient.Lock()
	defer global.MuClient.Unlock()
	for client := range global.Clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			global.Logger.Error("ERROR : ", zap.Error(err))
			client.Close()
			delete(global.Clients, client)
		}
	}
}

func addClientOneGeo(conn *websocket.Conn) {
	global.MuClientOneGeo.Lock()
	defer global.MuClientOneGeo.Unlock()
	global.ClientsOneGeo[conn] = true
}
func removeClientOneGeo(conn *websocket.Conn) {
	global.MuClientOneGeo.Lock()
	defer global.MuClientOneGeo.Unlock()
	conn.Close()
	delete(global.ClientsOneGeo, conn)
}



func broadcastMessageOneGeo(message []byte) {
	global.MuClientOneGeo.Lock()
	defer global.MuClientOneGeo.Unlock()
	for client := range global.ClientsOneGeo {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			global.Logger.Error("ERROR : ", zap.Error(err))
			client.Close()
			delete(global.ClientsOneGeo, client)
		}
	}
}

func Sender() {
	route := gin.Default()
	route.GET("/radar240", func(c *gin.Context) {
		conn, err := global.Upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			global.Logger.Error("ERROR : ", zap.Error(err))
			return
		}
		addClient(conn)
		defer removeClient(conn)
		for {
			message := <-global.ParsedData
			global.NUMBER_OF_SENT_MESSAGES.Inc()
			broadcastMessage(message)
		}
	})
	// route.GET("/radar240/geo", func(c *gin.Context) {
	// 	conn, err := global.Upgrader.Upgrade(c.Writer, c.Request, nil)
	// 	if err != nil {
	// 		global.Logger.Error("ERROR : ", zap.Error(err))
	// 		return
	// 	}
	// 	addClientOneGeo(conn)
	// 	defer removeClientOneGeo(conn)
	// 	for {
	// 		message := <-global.ParsedDataOneGeo
	// 		global.NUMBER_OF_SENT_MESSAGES.Inc()
	// 		broadcastMessageOneGeo(message)
	// 	}
	// })
	route.Run(":" + os.Getenv("WEBSOCKET_SERVER_PORT"))
}
