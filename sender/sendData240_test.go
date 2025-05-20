package sender_test

import (
	"net/url"
	"os"
	"cat240/global"
	"cat240/sender"
	"testing"
	"time"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestSendData240(t *testing.T) {
	os.Setenv("WEBSOCKET_SERVER_PORT", "8080")
	global.ParsedData = make(chan []byte, 30)
	go sender.Sender()
	time.Sleep(5 * time.Second)
	serverUrl := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/radar240"}
	conn , _, err := websocket.DefaultDialer.Dial(serverUrl.String(), nil)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()
	global.ParsedData <- []byte("test data")
	t.Run("Successful data received and  disconnection" , func(t *testing.T) {
		_, message, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("Failed to read message: %v", err)
		}
		assert.Equal(t, "test data", string(message))
		conn.Close()
		global.ParsedData <- []byte("test data")
		global.ParsedData <- []byte("test data")
		time.Sleep(3 * time.Second)
		global.MuClient.Lock()
		assert.Equal(t, 0, len(global.Clients))
		global.MuClient.Unlock()
	})
}