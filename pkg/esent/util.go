package esent

import (
	"encoding/binary"
	"math"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

func Float32frombytes(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}

func Float64frombytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

// Convert a slice of bytes to string, strip whitespace, and split by newline.
func SplitLinesBytes(input []byte) []string {
	return strings.Split(strings.TrimSpace(string(input[:])), "\n")
}

// Remove duplicates from a slice of strings.
func RemoveDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}

	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}

	return list
}

// Filter a slice using a test function.
//
// NOTE: This requires Go 1.18+
func Filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

// Return a string representing the date and time in ISO 8601 format.
func Isoformat(t time.Time) string {
	return t.Format("2006-01-02T15:04:05+00:00")
}

// Return a string representing the current date and time (UTC) in ISO 8601 format.
func IsoformatUtcNow() string {
	return Isoformat(time.Now().UTC())
}

func GetLogger() *log.Logger {
	logger := log.New(os.Stderr)
	return logger
}
