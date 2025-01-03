// Copyright (c) 2024 Focela Technologies. All rights reserved.
// Internal use only. Unauthorized use is prohibited.
// Contact: legal@focela.com

// Package binary provides APIs for handling binary/bytes data.
package binary

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/focela/loom/internal/core"
	"github.com/focela/loom/pkg/errors"
)

// LeEncode encodes multiple values into bytes using LittleEndian.
func LeEncode(values ...interface{}) []byte {
	buf := new(bytes.Buffer)
	for _, value := range values {
		if value == nil {
			return buf.Bytes()
		}
		switch v := value.(type) {
		case int:
			buf.Write(LeEncodeInt(v))
		case int8:
			buf.Write(LeEncodeInt8(v))
		case int16:
			buf.Write(LeEncodeInt16(v))
		case int32:
			buf.Write(LeEncodeInt32(v))
		case int64:
			buf.Write(LeEncodeInt64(v))
		case uint:
			buf.Write(LeEncodeUint(v))
		case uint8:
			buf.Write(LeEncodeUint8(v))
		case uint16:
			buf.Write(LeEncodeUint16(v))
		case uint32:
			buf.Write(LeEncodeUint32(v))
		case uint64:
			buf.Write(LeEncodeUint64(v))
		case bool:
			buf.Write(LeEncodeBool(v))
		case string:
			buf.Write(LeEncodeString(v))
		case []byte:
			buf.Write(v)
		case float32:
			buf.Write(LeEncodeFloat32(v))
		case float64:
			buf.Write(LeEncodeFloat64(v))
		default:
			if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
				core.Errorf(context.TODO(), `%+v`, err)
				buf.Write(LeEncodeString(fmt.Sprintf("%v", v)))
			}
		}
	}
	return buf.Bytes()
}

// LeEncodeByLength encodes values to bytes and ensures the result matches the specified length.
func LeEncodeByLength(length int, values ...interface{}) []byte {
	b := LeEncode(values...)
	if len(b) < length {
		b = append(b, make([]byte, length-len(b))...)
	} else if len(b) > length {
		b = b[:length]
	}
	return b
}

// LeEncodeString converts a string into a byte slice.
func LeEncodeString(s string) []byte {
	return []byte(s)
}

// LeEncodeBool converts a boolean value into a byte.
func LeEncodeBool(b bool) []byte {
	if b {
		return []byte{1}
	}
	return []byte{0}
}

// LeEncodeInt encodes an integer into a byte slice using LittleEndian.
func LeEncodeInt(i int) []byte {
	if i <= math.MaxInt8 {
		return LeEncodeInt8(int8(i))
	} else if i <= math.MaxInt16 {
		return LeEncodeInt16(int16(i))
	} else if i <= math.MaxInt32 {
		return LeEncodeInt32(int32(i))
	}
	return LeEncodeInt64(int64(i))
}

// LeEncodeUint encodes an unsigned integer into a byte slice using LittleEndian.
func LeEncodeUint(i uint) []byte {
	if i <= math.MaxUint8 {
		return LeEncodeUint8(uint8(i))
	} else if i <= math.MaxUint16 {
		return LeEncodeUint16(uint16(i))
	} else if i <= math.MaxUint32 {
		return LeEncodeUint32(uint32(i))
	}
	return LeEncodeUint64(uint64(i))
}

// LeEncodeInt8 encodes an int8 value into a single byte.
func LeEncodeInt8(i int8) []byte { return []byte{byte(i)} }

// LeEncodeUint8 encodes a uint8 value into a single byte.
func LeEncodeUint8(i uint8) []byte { return []byte{i} }

// LeEncodeInt16 encodes an int16 value into a 2-byte slice.
func LeEncodeInt16(i int16) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(i))
	return b
}

// LeEncodeUint16 encodes a uint16 value into a 2-byte slice.
func LeEncodeUint16(i uint16) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, i)
	return b
}

// LeEncodeInt32 encodes an int32 value into a 4-byte slice.
func LeEncodeInt32(i int32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(i))
	return b
}

// LeEncodeUint32 encodes a uint32 value into a 4-byte slice.
func LeEncodeUint32(i uint32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, i)
	return b
}

// LeEncodeInt64 encodes an int64 value into an 8-byte slice.
func LeEncodeInt64(i int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	return b
}

// LeEncodeUint64 encodes a uint64 value into an 8-byte slice.
func LeEncodeUint64(i uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, i)
	return b
}

// LeEncodeFloat32 encodes a float32 value into a 4-byte slice.
func LeEncodeFloat32(f float32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, math.Float32bits(f))
	return b
}

// LeEncodeFloat64 encodes a float64 value into an 8-byte slice.
func LeEncodeFloat64(f float64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, math.Float64bits(f))
	return b
}

// LeDecode decodes bytes into specified values using LittleEndian.
func LeDecode(b []byte, values ...interface{}) error {
	buf := bytes.NewBuffer(b)
	for _, value := range values {
		if err := binary.Read(buf, binary.LittleEndian, value); err != nil {
			return errors.Wrap(err, `binary.Read failed`)
		}
	}
	return nil
}

// LeDecodeToBool decodes a byte slice into a boolean value.
func LeDecodeToBool(b []byte) bool {
	return len(b) > 0 && b[0] != 0
}

// LeDecodeToInt8 decodes a byte slice into an int8 value.
func LeDecodeToInt8(b []byte) int8 { return int8(b[0]) }

// LeDecodeToUint8 decodes a byte slice into a uint8 value.
func LeDecodeToUint8(b []byte) uint8 { return b[0] }

// LeDecodeToInt16 decodes a byte slice into an int16 value.
func LeDecodeToInt16(b []byte) int16 {
	return int16(binary.LittleEndian.Uint16(LeFillUpSize(b, 2)))
}

// LeDecodeToUint16 decodes a byte slice into a uint16 value.
func LeDecodeToUint16(b []byte) uint16 {
	return binary.LittleEndian.Uint16(LeFillUpSize(b, 2))
}

// LeFillUpSize ensures a byte slice has the required length by padding zeros.
func LeFillUpSize(b []byte, l int) []byte {
	if len(b) >= l {
		return b[:l]
	}
	c := make([]byte, l)
	copy(c, b)
	return c
}
