package utils

import (
	"bufio"
	"cat240/global"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	// "text/scanner"
	"time"

	"go.uber.org/zap"
)

const (
	// Buffer size for reading radar data
	BUFFER_SIZE = 1000000
	// Delay between reconnection attempts
	RECONNECT_DELAY = 5 * time.Second
)

// ReadData continuously reads data from a radar connection.
// It establishes a TCP connection to the radar, authenticates, and streams data.
// The function handles reconnection attempts and data buffering.
func ReadData() {
	file, err := os.Open("data.txt")
	
	if err != nil {
		fmt.Println("Error opening file:", err)
		return 
	}
	defer file.Close()
	for {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			hexData, err := hex.DecodeString(line)
			if err != nil {
				global.Logger.Error("Failed to decode hex string", zap.Error(err))
				continue
			}
			global.FilteredData <- hexData
			time.Sleep(100 * time.Microsecond) // Simulate processing delay
			
		// conn, err := establishConnection()
		// if err != nil {
		// 	global.Logger.Error("Failed to establish connection", zap.Error(err))
		// 	time.Sleep(RECONNECT_DELAY)
		// 	continue
		// }

		// if err := authenticateConnection(conn); err != nil {
		// 	global.Logger.Error("Failed to authenticate", zap.Error(err))
		// 	conn.Close()
		// 	continue
		// }

		// if err := readAndProcessData(conn); err != nil {
		// 	global.Logger.Error("Error during data reading", zap.Error(err))
		}

		// conn.Close()
	}
}

// establishConnection creates a TCP connection to the radar.
func establishConnection() (net.Conn, error) {
	radarAddr := fmt.Sprintf("%s:%s", os.Getenv("RADAR_IP"), os.Getenv("RADAR_PORT"))
	return net.Dial("tcp", radarAddr)
}

// authenticateConnection sends the authentication password to the radar.
func authenticateConnection(conn net.Conn) error {
	password := os.Getenv("RADAR_PASSWORD")
	_, err := conn.Write([]byte(password + "\n"))
	return err
}

// readAndProcessData reads data from the connection and sends it to the FilteredData channel.
// It handles connection errors and channel buffering.
func readAndProcessData(conn net.Conn) error {
	buffer := make([]byte, BUFFER_SIZE)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			return err
		}

		select {
		case global.FilteredData <- buffer[:n]:
			global.NUMBER_OF_RECEIVED_MESSAGES.Inc()
		default:
			global.Logger.Info("FilteredData channel is full, dropping data")
		}
	}
}
