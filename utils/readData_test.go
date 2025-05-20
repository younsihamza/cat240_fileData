package utils_test

import (
	"cat240/global"
	"cat240/utils"
	"net"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockServer struct {
	listener net.Listener
	data	[]byte
}

func CreateMockServer(t *testing.T) *mockServer {
	connection , err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create mock server")
	}
	return &mockServer{listener : connection}
}

func TestReadData(t *testing.T) {
	global.FilteredData = make(chan []byte, 2)
	global.ParsedData = make(chan []byte, 2)

	server := CreateMockServer(t)
	defer server.listener.Close()
	addr := server.listener.Addr().(*net.TCPAddr)
    os.Setenv("RADAR_IP", "127.0.0.1")
    os.Setenv("RADAR_PORT", strconv.Itoa(addr.Port))
    os.Setenv("RADAR_PASSWORD", "testpass")
	t.Run("TestReadData", func(t *testing.T) {
		go utils.ReadData()
		conn, err := server.listener.Accept()
		defer conn.Close()
		buf := make([]byte, 1024)
		n , err := conn.Read(buf)
		assert.NoError(t, err)
		assert.Equal(t, "testpass\n", string(buf[:n]))
		conn.Write([]byte("data"))
		select {
			case data := <-global.FilteredData:
				assert.Equal(t, "data", string(data))
			case <-time.After(time.Second):
				t.Fatalf("Failed to read data")
		}
	})

	
}
