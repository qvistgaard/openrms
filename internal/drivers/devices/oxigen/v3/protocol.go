package v3

import (
	"encoding/binary"
	"time"
)

func unpackLapTime(high byte, low byte) time.Duration {
	be := binary.BigEndian.Uint64([]byte{0, 0, 0, 0, 0, 0, high, low})
	lt := float64(be) / 99.25
	ltd := time.Duration(lt * float64(time.Second))
	return ltd
}

func unpackRaceTime(b [4]byte, lag byte) time.Duration {
	be := make([]byte, 8)
	copy(be[4:], b[:])

	rt := binary.BigEndian.Uint64(be) - uint64(lag)
	// rt := (uint(b[0]) * 16777216) + (uint(b[1]) * 65536) + (uint(b[2]) * 256) + uint(b[3]) - uint(lag)
	return time.Duration(rt*10) * time.Millisecond
}

func unpackPitStatus(b dongleRxMessage) bool {
	return 0x40&b[8] == 0x40
}

func packRaceCounter(start time.Time) []byte {
	centiSeconds := time.Now().Sub(start).Milliseconds() / 10
	be := make([]byte, 8)
	binary.BigEndian.PutUint64(be, uint64(centiSeconds))
	return be[len(be)-3:]
}

func unpackRaceCounter(b [3]byte) time.Duration {
	be := make([]byte, 8)
	copy(be[5:], b[:])

	u := binary.BigEndian.Uint64(be)
	return time.Duration(u*10) * time.Millisecond
}

func percentageToByte(percent uint8) uint8 {
	if percent > 100 {
		percent = 100
	}
	return uint8(255.0 * (float64(percent) / 100.0))
}
