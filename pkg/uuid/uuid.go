package uuid

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"time"
)

// UUID represents a UUID type.
type UUID [16]byte

// NewV7 generates a UUID v7 (time-ordered).
// The first 48 bits are the timestamp (milliseconds since Unix Epoch),
// the next 12 bits are random (for monotonicity within the same ms),
// and the last 64 bits are random.
func NewV7() (UUID, error) {
	var u UUID
	_, err := io.ReadFull(rand.Reader, u[6:]) // Fill random part first (bytes 6 to 15)
	if err != nil {
		return UUID{}, err
	}

	// Set timestamp (first 48 bits)
	// Unix time in milliseconds
	t := uint64(time.Now().UnixMilli())

	// Bytes 0-5: Big Endian timestamp
	u[0] = byte(t >> 40)
	u[1] = byte(t >> 32)
	u[2] = byte(t >> 24)
	u[3] = byte(t >> 16)
	u[4] = byte(t >> 8)
	u[5] = byte(t)

	// Set version bits (bits 48-51 of the UUID, which is byte 6 top 4 bits)
	// Version is 7 (0111)
	u[6] = (u[6] & 0x0F) | 0x70

	// Set variant bits (bits 64-65 of the UUID, which is byte 8 top 2 bits)
	// Variant is RFC 4122 (10xx xxxx)
	u[8] = (u[8] & 0x3F) | 0x80

	return u, nil
}

// Must returns uuid if err is nil else panics.
func Must(u UUID, err error) UUID {
	if err != nil {
		panic(err)
	}
	return u
}

// String returns the string form of uuid, xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func (u UUID) String() string {
	buf := make([]byte, 36)
	hex.Encode(buf[:8], u[:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], u[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], u[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], u[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], u[10:])
	return string(buf)
}
