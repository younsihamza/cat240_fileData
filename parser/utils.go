package cat240

import (
	"encoding/hex"
	"strings"
)
// HexStringToBytes convert hex string to bytes and remove all spaces, tabs and new lines
func HexStringToBytes(s string) ([]byte, error) {
	s =  strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\t", "")
	return hex.DecodeString(s)
} 


