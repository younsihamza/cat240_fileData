package utils

import (
	"cat240/global"
	cat240 "cat240/parser"
	"encoding/json"
	"time"

	"go.uber.org/zap"
)

// ParseData continuously processes radar data from the FilteredData channel.
// It parses the data, decodes it, and sends the JSON-encoded result to the ParsedData channel.
// The function includes a timeout check to monitor data processing activity.
func ParseData() {
	const timeoutDuration = 5 * time.Second

	for {
		select {
		case data := <-global.FilteredData:
			if err := processData(data); err != nil {
				global.Logger.Error("Failed to process data", zap.Error(err))
				global.NUMBER_OF_FAILED_MESSAGES.Inc()
				continue
			}
		case <-time.After(timeoutDuration):
			global.Logger.Info("No data to parse in the last 5 seconds")
		}
	}
}

// processData handles the parsing and encoding of a single data block.
// It returns an error if any step of the process fails.
func processData(data []byte) error {
	// Parse the data
	validData, err := cat240.Parser(data)
	if err != nil {
		return err
	}

	// Decode the parsed data
	multipleGeometry , _ := cat240.Decode(validData)
	if multipleGeometry == nil {
		return nil
	}

	// Convert to JSON
	jsonData, err := json.Marshal(multipleGeometry)
	if err != nil {
		return err
	}
	// onegeo, err := json.Marshal(OneGeometry)
	// if err != nil {
	// 	return err
	// }

	// Send to output channel
	global.ParsedData <- jsonData
	// global.ParsedDataOneGeo <- onegeo
	global.NUMBER_OF_PARSED_MESSAGES.Inc()
	return nil
}
