package utils

import "encoding/hex"

func RestoreToHexString(str string) string {
	s := hex.EncodeToString([]byte(str))
	return s
}

func RestoreToByteString(str string) string {
	byteSlice, err := hex.DecodeString(str)
	if err != nil {
		return str
	}

	return string(byteSlice)
}
