package cat240

import (
	"bytes"
	"cat240/global"
	"compress/zlib"
	"fmt"
	"io"
	"math"

	"go.uber.org/zap"
)

var (
	coordinates []BlockData = nil
	last_start_azimuth = 0.0
	last_end_azimuth = 0.0	
	Last_range = 1
	current_Range = 1
)

const (
	speedOfLight     = 299792458.0 // speed of light in meters per second
	earthRadius      = 6378137.0   // Earth radius in meters
	radiansToDegrees = 180.0 / math.Pi
	degreesToRadians = math.Pi / 180.0
)

// BlockData represents a single radar data point with its coordinates and intensity
type BlockData struct {
	Longitude    float64 `json:"longitude"`
	Latitude     float64 `json:"latitude"`
	Intensity    int     `json:"intensity"`
	StartAzimuth float64 `json:"start_azimuth"`
	EndAzimuth   float64 `json:"end_azimuth"`
	StartRange   float64 `json:"start_range"`
}



// Decode processes the CAT240 radar data and converts it to GeoJSON format
func Decode(data *ValidData) (*map[string]interface{}, *map[string]interface{}) {
	if checkHighOrderBit(data.VideoCellsResolution.CompressionIndicator) {
		data.VideoBlock = decompressData(data.VideoBlock)
	}
	bitResolution(data, int(data.VideoCellsResolution.BitResolution))
	fmt.Println("bitResolution", data.VideoCellsResolution.BitResolution)
	var readyCoordinatesToSend []BlockData
	if int(int(data.VideoHeader.StartAzimuth) / 2) + 1 != Last_range {
		readyCoordinatesToSend = coordinates
		Last_range = int(int(data.VideoHeader.StartAzimuth) / 2) + 1
		coordinates = nil
	}
	coordinates = append(coordinates, coordinateTransformation(data)...)
	current_Range = int(int(data.VideoHeader.StartAzimuth) / 2) + 1
	if len(readyCoordinatesToSend) > 0 {
		return toGeoJSON(readyCoordinatesToSend) , toOneGeometry(coordinates)
	}
	if len(coordinates) > 0 {
		return nil , toOneGeometry(coordinates)
	}
	return nil , nil
}


func toOneGeometry(data []BlockData) *map[string]interface{}{
	Geometrys := make([]map[string]interface{},0)
	for _, block := range data {
		if block.Intensity > 1 {
			Geometrys = append(Geometrys, map[string]interface{}{
				"type":        "Feature",
				"geometry":    map[string]interface{}{"type": "Point", "coordinates": []float64{block.Longitude, block.Latitude}},
				"properties":  map[string]interface{}{"intensity": block.Intensity},
			})
		}
	}
	return &map[string]interface{}{
		"type":        "FeatureCollection",
		"features":    Geometrys,
	}
}

// coordinateTransformation converts radar data to geographic coordinates
func coordinateTransformation(data *ValidData) []BlockData {
	var coordinates []BlockData
	rangeCell := data.VideoHeader.CellDuration * speedOfLight / 2.0
	azimuthIncrement := (data.VideoHeader.EndAzimuth - data.VideoHeader.StartAzimuth) / float64(data.VideoOctetsVideoCellCounters.ValidCellsInVideoBlock)
	currentRange := rangeCell * float64(len(data.VideoBlock)-1+data.VideoHeader.StartRange)

	// Add initial point
	coordinates = append(coordinates, BlockData{
		Longitude:    0,
		Latitude:     0,
		Intensity:    0,
		StartAzimuth: data.VideoHeader.StartAzimuth,
		EndAzimuth:   data.VideoHeader.EndAzimuth,
		StartRange:   currentRange,
	})

	// Process each data point
	for i := 0; i < len(data.VideoBlock)-1; i++ {
		currentRange := rangeCell * float64(i + data.VideoHeader.StartRange)
		currentAzimuth := data.VideoHeader.StartAzimuth + azimuthIncrement * float64(i)
		x, y := polarToCartesian(currentRange, currentAzimuth)
		lat, longitude := cartesianToGeo(33.126474, -8.641668, x, y)
		coordinates = append(coordinates, BlockData{
			Longitude:    longitude,
			Latitude:     lat,
			Intensity:    int(data.VideoBlock[i]),
			StartAzimuth: data.VideoHeader.StartAzimuth,
			EndAzimuth:   data.VideoHeader.EndAzimuth,
			StartRange:   currentRange,
		})
	}

	return coordinates
}

// cartesianToGeo converts Cartesian coordinates to geographic coordinates
func cartesianToGeo(originLat, originLon, x, y float64) (float64, float64) {
	originLatRad := originLat * degreesToRadians
	deltaLatRad := y / earthRadius
	deltaLonRad := x / (earthRadius * math.Cos(originLatRad))

	deltaLatDeg := deltaLatRad * radiansToDegrees
	deltaLonDeg := deltaLonRad * radiansToDegrees

	return originLat + deltaLatDeg, originLon + deltaLonDeg
}

// polarToCartesian converts polar coordinates to Cartesian coordinates
func polarToCartesian(rangeCell, azimuth float64) (float64, float64) {
	azimuthRad := azimuth * degreesToRadians
	x := rangeCell * math.Sin(azimuthRad)
	y := rangeCell * math.Cos(azimuthRad)
	return x, y
}

// decompressData decompresses zlib-compressed data
func decompressData(data []byte) []byte {
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		global.Logger.Error("Failed to create zlib reader", zap.Error(err))
		return nil
	}
	defer r.Close()

	decompressed, err := io.ReadAll(r)
	if err != nil {
		global.Logger.Error("Failed to read decompressed data", zap.Error(err))
		return nil
	}
	return decompressed
}

// checkHighOrderBit checks if the high-order bit of a byte is set
func checkHighOrderBit(b byte) bool {
	return (b >> 7) == 0x01
}

// bitResolution processes the video data based on bit resolution
func bitResolution(data *ValidData, bitPerCell int) {
	var videoData []byte
	switch bitPerCell {
	case 1, 2, 4, 8:
		for i := 0; i < len(data.VideoBlock); i++ {
			for j := 0; j < 8; j += bitPerCell {
				mask := (1 << bitPerCell) - 1
				shift := 8 - bitPerCell - j
				videoData = append(videoData, (data.VideoBlock[i]>>shift)&byte(mask))
			}
		}
	}
	data.VideoBlock = videoData
}

// toGeoJSON converts the radar data to GeoJSON format
func toGeoJSON(data []BlockData) *map[string]interface{} {
	coordinateWithBigerOpacity := make([]map[string]interface{}, 0)
	for _, block := range data {
		if block.Intensity > 1 {
			coordinateWithBigerOpacity = append(coordinateWithBigerOpacity, map[string]interface{}{"opacity":  block.Intensity, "coordinates": []float64{block.Longitude, block.Latitude}})
		}
}
	return &map[string]interface{}{
		"coordinates": coordinateWithBigerOpacity,
		"range_id":  current_Range,
	}
}
