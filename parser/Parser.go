package cat240

import (
	"fmt"
	"math"
)

const (
	cat240Identifier = 240
	bitsPerByte      = 8
	azimuthScale     = 360.0 / 65535.0

	// Field sizes in bytes
	dataSourceIdentifierSize = 2
	messageTypeSize          = 1
	recordHeaderSize         = 4
	videoHeaderSize          = 12
	videoCellsResolutionSize = 2
	videoOctetsCountersSize  = 5
	timeOfDaySize            = 3
	reservedFieldSize        = 1

	// Block multipliers
	lowBlockMultiplier    = 4
	mediumBlockMultiplier = 64
	highBlockMultiplier   = 256

	// Time scales
	nanoScale  = -9
	femtoScale = -15

	// Error messages
	errDataTooShort        = "data too short: expected %d bytes"
	errDataTooShortAtLeast = "data too short: expected at least %d bytes"
)

// Parser parses CAT240 radar data and returns a ValidData struct or an error if parsing fails.
// The input data should be a byte slice containing the CAT240 message.
func Parser(data []byte) (*ValidData, error) {
	if len(data) < 3 {
		return nil, fmt.Errorf("data too short: expected at least 3 bytes, got %d", len(data))
	}
	if data[0] != cat240Identifier {
		return nil, fmt.Errorf("invalid CAT240 identifier: expected %d, got %d", cat240Identifier, data[0])
	}

	messageLength := (int(data[1]) << 8) + int(data[2])
	if messageLength != len(data) {
		return nil, fmt.Errorf("invalid message length: expected %d, got %d", messageLength, len(data))
	}

	currentByte := 3
	fspec := fmt.Sprintf("%08b", int64(data[currentByte]))
	if fspec[7] == '1' {
		currentByte++
		fspec = fmt.Sprintf("%016b", int64(int(data[3])<<8+int(data[4])))
	}
	currentByte++

	return convertToVideoDataItem(fspec, data[currentByte:])
}

// convertToVideoDataItem converts the data to a VideoDataItem struct based on the FSPEC field.
func convertToVideoDataItem(fspec string, data []byte) (*ValidData, error) {
	var item VideoDataItem
	remainingData := data

	for i, val := range fspec {
		fieldIndex := i + 1
		if fieldIndex == 8 && val != '1' {
			return nil, fmt.Errorf("invalid FSPEC: last bit must be 1")
		}
		if val != '1' {
			continue
		}

		var err error
		remainingData, err = processField(fieldIndex, remainingData, &item)
		if err != nil {
			return nil, fmt.Errorf("field %d: %w", fieldIndex, err)
		}
	}

	return validator(item)
}

// processField handles the processing of a single field based on its index.
func processField(fieldIndex int, data []byte, item *VideoDataItem) ([]byte, error) {
	switch fieldIndex {
	case 1:
		return processDataSourceIdentifier(data, item)
	case 2:
		return processMessageType(data, item)
	case 3:
		return processRecordHeader(data, item)
	case 4:
		return processVideoSummary(data, item)
	case 5:
		return processVideoHeaderNano(data, item)
	case 6:
		return processVideoHeaderFemto(data, item)
	case 7:
		return processVideoCellsResolution(data, item)
	case 9:
		return processVideoOctetsVideoCellCounters(data, item)
	case 10:
		return processVideoBlock(data, item, lowBlockMultiplier, &item.VideoBlockLowDataVolume)
	case 11:
		return processVideoBlock(data, item, mediumBlockMultiplier, &item.VideoBlockMediumDataVolume)
	case 12:
		return processVideoBlock(data, item, highBlockMultiplier, &item.VideoBlockHighDataVolume)
	case 13:
		return processTimeOfDay(data, item)
	case 14:
		return processReservedExpansionField(data, item)
	case 15:
		return processSpecialPurposeField(data, item)
	default:
		return data, nil
	}
}

// Helper functions for processing each field
func processDataSourceIdentifier(data []byte, item *VideoDataItem) ([]byte, error) {
	if len(data) < dataSourceIdentifierSize {
		return nil, fmt.Errorf(errDataTooShort, dataSourceIdentifierSize)
	}
	item.DataSourceIndentifier = &DataSourceIndentifier{
		SAC: int(data[0]),
		SIC: int(data[1]),
	}
	return data[dataSourceIdentifierSize:], nil
}

func processMessageType(data []byte, item *VideoDataItem) ([]byte, error) {
	if len(data) < messageTypeSize {
		return nil, fmt.Errorf(errDataTooShort, messageTypeSize)
	}
	item.MessageType = int(data[0])
	return data[messageTypeSize:], nil
}

func processRecordHeader(data []byte, item *VideoDataItem) ([]byte, error) {
	if len(data) < recordHeaderSize {
		return nil, fmt.Errorf(errDataTooShort, recordHeaderSize)
	}
	item.RecodeHeader = int(data[0])<<24 + int(data[1])<<16 + int(data[2])<<8 + int(data[3])
	return data[recordHeaderSize:], nil
}

func processVideoHeaderNano(data []byte, item *VideoDataItem) ([]byte, error) {
	if len(data) < videoHeaderSize {
		return nil, fmt.Errorf(errDataTooShort, videoHeaderSize)
	}
	item.VideoHeaderNano = createVideoHeader(data, nanoScale)
	return data[videoHeaderSize:], nil
}

func processVideoHeaderFemto(data []byte, item *VideoDataItem) ([]byte, error) {
	if len(data) < videoHeaderSize {
		return nil, fmt.Errorf(errDataTooShort, videoHeaderSize)
	}
	item.VideoHeaderFemto = createVideoHeader(data, femtoScale)
	return data[videoHeaderSize:], nil
}

// createVideoHeader creates a VideosHeaders struct with the given data and time scale
func createVideoHeader(data []byte, timeScale int) *VideosHeaders {
	return &VideosHeaders{
		StartAzimuth: float64(int(data[0])<<8+int(data[1])) * azimuthScale,
		EndAzimuth:   float64(int(data[2])<<8+int(data[3])) * azimuthScale,
		StartRange:   int(data[4])<<24 + int(data[5])<<16 + int(data[6])<<8 + int(data[7]),
		CellDuration: float64(int(data[8])<<24+int(data[9])<<16+int(data[10])<<8+int(data[11])) * math.Pow(10, float64(timeScale)),
	}
}

func processVideoCellsResolution(data []byte, item *VideoDataItem) ([]byte, error) {
	if len(data) < videoCellsResolutionSize {
		return nil, fmt.Errorf(errDataTooShort, videoCellsResolutionSize)
	}
	item.VideoCellsResolution = &VideoCellsResolution{
		CompressionIndicator: data[0],
		BitResolution:        int(data[1]),
	}
	return data[videoCellsResolutionSize:], nil
}

func processVideoOctetsVideoCellCounters(data []byte, item *VideoDataItem) ([]byte, error) {
	if len(data) < videoOctetsCountersSize {
		return nil, fmt.Errorf(errDataTooShort, videoOctetsCountersSize)
	}
	item.VideoOctetsVideoCellCounters = &VideoOctetsVideoCellCounters{
		ValidOctetsInVideoBlock: int(data[0])<<8 + int(data[1]),
		ValidCellsInVideoBlock:  int(data[2])<<16 + int(data[3])<<8 + int(data[4]),
	}
	return data[videoOctetsCountersSize:], nil
}

// processVideoBlock handles processing of video blocks with different multipliers
func processVideoBlock(data []byte, item *VideoDataItem, multiplier int, target *[]byte) ([]byte, error) {
	if len(data) < 1 {
		return nil, fmt.Errorf(errDataTooShortAtLeast, 1)
	}
	blockSize := int(data[0]) * multiplier
	if blockSize+1 > len(data) {
		return nil, fmt.Errorf(errDataTooShort, blockSize+1)
	}
	*target = data[1 : blockSize+1]
	return data[blockSize+1:], nil
}

func processTimeOfDay(data []byte, item *VideoDataItem) ([]byte, error) {
	if len(data) < timeOfDaySize {
		return nil, fmt.Errorf(errDataTooShort, timeOfDaySize)
	}
	item.TimeOfDay = float64(int(data[0])<<16 + int(data[1])<<8 + int(data[2]))
	return data[timeOfDaySize:], nil
}

func processReservedExpansionField(data []byte, item *VideoDataItem) ([]byte, error) {
	if len(data) < reservedFieldSize {
		return nil, fmt.Errorf(errDataTooShort, reservedFieldSize)
	}
	item.ReservedExpansionField = string(data[0])
	return data[reservedFieldSize:], nil
}

func processSpecialPurposeField(data []byte, item *VideoDataItem) ([]byte, error) {
	if len(data) < 1 {
		return nil, fmt.Errorf(errDataTooShortAtLeast, 1)
	}
	fieldLength := int(data[0])
	if fieldLength > len(data) {
		return nil, fmt.Errorf(errDataTooShort, fieldLength)
	}
	item.SpecialPurposeField = data[1:fieldLength]
	return data[fieldLength:], nil
}

func processVideoSummary(data []byte, item *VideoDataItem) ([]byte, error) {
	if len(data) < 1 {
		return nil, fmt.Errorf(errDataTooShortAtLeast, 1)
	}
	summaryLength := int(data[0])
	if summaryLength+1 > len(data) {
		return nil, fmt.Errorf(errDataTooShort, summaryLength+1)
	}
	item.VideoSummary = data[1 : summaryLength+1]
	return data[summaryLength+1:], nil
}

// validator validates the VideoDataItem and returns a ValidData struct or an error if validation fails.
func validator(item VideoDataItem) (*ValidData, error) {
	if !isValidVideoDataItem(item) {
		return nil, fmt.Errorf("invalid video data item: missing mandatory fields")
	}

	validData := ValidData{
		DataSourceIndentifier:        item.DataSourceIndentifier,
		MessageType:                  item.MessageType,
		RecodeHeader:                 item.RecodeHeader,
		VideoCellsResolution:         item.VideoCellsResolution,
		VideoOctetsVideoCellCounters: item.VideoOctetsVideoCellCounters,
		TimeOfDay:                    item.TimeOfDay,
		ReservedExpansionField:       item.ReservedExpansionField,
		SpecialPurposeField:          item.SpecialPurposeField,
	}

	// Set video header based on available data
	if item.VideoHeaderNano != nil {
		validData.VideoHeader = item.VideoHeaderNano
	} else if item.VideoHeaderFemto != nil {
		validData.VideoHeader = item.VideoHeaderFemto
	}

	// Set video block based on available data
	if item.VideoBlockLowDataVolume != nil {
		validData.VideoBlock = item.VideoBlockLowDataVolume
	} else if item.VideoBlockMediumDataVolume != nil {
		validData.VideoBlock = item.VideoBlockMediumDataVolume
	} else if item.VideoBlockHighDataVolume != nil {
		validData.VideoBlock = item.VideoBlockHighDataVolume
	}

	return &validData, nil
}

// isValidVideoDataItem checks if all mandatory fields are present in the VideoDataItem.
func isValidVideoDataItem(item VideoDataItem) bool {
	return item.DataSourceIndentifier != nil &&
		item.MessageType != 0 &&
		item.RecodeHeader != 0 &&
		(item.VideoHeaderNano != nil || item.VideoHeaderFemto != nil) &&
		item.VideoCellsResolution != nil &&
		item.VideoOctetsVideoCellCounters != nil &&
		(item.VideoBlockLowDataVolume != nil ||
			item.VideoBlockMediumDataVolume != nil ||
			item.VideoBlockHighDataVolume != nil)
}
