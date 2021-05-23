package utilities

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"runtime/trace"
	"time"
)

/* ByteTo* functions used to simplify Modbus register value conversions. */

func BytesToInt16(bytes []byte) int16 {
	resultsUint16 := binary.BigEndian.Uint16(bytes)
	int16_value := int16(resultsUint16)
	return int16_value
}
func BytesToUInt16(bytes []byte) uint16 {
	resultsUint16 := binary.BigEndian.Uint16(bytes)
	return resultsUint16
}

func BytesToUint32(bytes []byte) uint32 {
	resultsUint32 := binary.BigEndian.Uint32(bytes)
	return resultsUint32
}

func TimeNowInUnixMs() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// get a [a-z0-9] random {strlen} long string.
func RandomString(strlen int) string {
	var r *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}

// Trace to file fd
var traceFile *os.File

// Start Trace to file (trace-{random}.out)
func StartTrace() {
	println("Tracing...")

	var err error
	traceFile, err = os.Create(fmt.Sprintf("trace-%s.out", RandomString(10)))
	if err != nil {
		panic(err)
	}

	err = trace.Start(traceFile)
	if err != nil {
		panic(err)
	}
}

// Stop Trace to file.
func StopTrace() {
	trace.Stop()
	traceFile.Close()
}
