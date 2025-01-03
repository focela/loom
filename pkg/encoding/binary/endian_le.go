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

func LeEncodeByLength(length int, values ...interface{}) []byte {
	b := LeEncode(values...)
	if len(b) < length {
		b = append(b, make([]byte, length-len(b))...)
	} else if len(b) > length {
		b = b[:length]
	}
	return b
}

func LeEncodeString(s string) []byte {
	return []byte(s)
}

func LeEncodeBool(b bool) []byte {
	if b {
		return []byte{1}
	}
	return []byte{0}
}

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

func LeEncodeInt8(i int8) []byte { return []byte{byte(i)} }

func LeEncodeUint8(i uint8) []byte { return []byte{i} }

func LeEncodeInt16(i int16) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(i))
	return b
}

func LeEncodeUint16(i uint16) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, i)
	return b
}

func LeEncodeInt32(i int32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(i))
	return b
}

func LeEncodeUint32(i uint32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, i)
	return b
}

func LeEncodeInt64(i int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	return b
}

func LeEncodeUint64(i uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, i)
	return b
}

func LeEncodeFloat32(f float32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, math.Float32bits(f))
	return b
}

func LeEncodeFloat64(f float64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, math.Float64bits(f))
	return b
}

func LeDecode(b []byte, values ...interface{}) error {
	buf := bytes.NewBuffer(b)
	for _, value := range values {
		if err := binary.Read(buf, binary.LittleEndian, value); err != nil {
			return errors.Wrap(err, `binary.Read failed`)
		}
	}
	return nil
}

func LeDecodeToString(b []byte) string {
	return string(b)
}

func LeDecodeToInt(b []byte) int {
	return int(LeDecodeToInt64(b))
}

func LeDecodeToUint(b []byte) uint {
	return uint(LeDecodeToUint64(b))
}

func LeDecodeToBool(b []byte) bool { return len(b) > 0 && b[0] != 0 }

func LeDecodeToInt8(b []byte) int8 { return int8(b[0]) }

func LeDecodeToUint8(b []byte) uint8 { return b[0] }

func LeDecodeToInt16(b []byte) int16 { return int16(binary.LittleEndian.Uint16(LeFillUpSize(b, 2))) }

func LeDecodeToUint16(b []byte) uint16 { return binary.LittleEndian.Uint16(LeFillUpSize(b, 2)) }

func LeDecodeToInt32(b []byte) int32 { return int32(binary.LittleEndian.Uint32(LeFillUpSize(b, 4))) }

func LeDecodeToUint32(b []byte) uint32 { return binary.LittleEndian.Uint32(LeFillUpSize(b, 4)) }

func LeDecodeToInt64(b []byte) int64 { return int64(binary.LittleEndian.Uint64(LeFillUpSize(b, 8))) }

func LeDecodeToUint64(b []byte) uint64 { return binary.LittleEndian.Uint64(LeFillUpSize(b, 8)) }

func LeDecodeToFloat32(b []byte) float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(LeFillUpSize(b, 4)))
}

func LeDecodeToFloat64(b []byte) float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(LeFillUpSize(b, 8)))
}

func LeFillUpSize(b []byte, l int) []byte {
	if len(b) >= l {
		return b[:l]
	}
	c := make([]byte, l)
	copy(c, b)
	return c
}
