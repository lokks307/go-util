package bytesbuilder

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"reflect"
	"strings"

	"github.com/btcsuite/btcutil/base58"
)

type ByteBuilder struct {
	BytesArr []byte
}

func NewBuilder() *ByteBuilder {
	return &ByteBuilder{}
}

func (m *ByteBuilder) Clear() {
	m.BytesArr = nil
}

func (m *ByteBuilder) Append(data ...interface{}) {
	var byteData []byte

	for idx := range data {
		byteData = nil
		switch t := data[idx].(type) {
		case bool:
			if t {
				byteData = append(byteData, 0x01)
			} else {
				byteData = append(byteData, 0x00)
			}
		case string:
			byteData = []byte(t)
		case uint8:
			byteData = append(byteData, t)
		case []uint8: // == []byte
			byteData = t
		case int16, uint16:
			var v uint16
			if reflect.TypeOf(t).Kind() == reflect.Int16 {
				v = uint16(t.(int16))
			} else {
				v = t.(uint16)
			}
			b := make([]byte, 2)
			binary.BigEndian.PutUint16(b, v)
			byteData = b
		case int32, uint32:
			var v uint32
			if reflect.TypeOf(t).Kind() == reflect.Int32 {
				v = uint32(t.(int32))
			} else {
				v = t.(uint32)
			}
			b := make([]byte, 4)
			binary.BigEndian.PutUint32(b, v)
			byteData = b
		case int64, uint64:
			var v uint64
			if reflect.TypeOf(t).Kind() == reflect.Int64 {
				v = uint64(t.(int64))
			} else {
				v = t.(uint64)
			}
			b := make([]byte, 8)
			binary.BigEndian.PutUint64(b, v)
			byteData = b
		}

		if byteData != nil {
			m.BytesArr = append(m.BytesArr, byteData...)
		}
	}
}

func (m *ByteBuilder) AppendHex(hexStrArr ...string) {
	for idx := range hexStrArr {
		if hexStrArr[idx] == "" {
			continue
		}

		if strings.HasPrefix(hexStrArr[idx], "0x") {
			if data, err := hex.DecodeString(hexStrArr[idx][2:]); err == nil {
				m.BytesArr = append(m.BytesArr, data...)
			}
		} else {
			if data, err := hex.DecodeString(hexStrArr[idx]); err == nil {
				m.BytesArr = append(m.BytesArr, data...)
			}
		}
	}
}

func (m *ByteBuilder) AppendBase64(b64StrArr ...string) {
	for idx := range b64StrArr {
		if b64StrArr[idx] == "" {
			continue
		}

		if decodedData, err := base64.StdEncoding.DecodeString(b64StrArr[idx]); err == nil {
			m.BytesArr = append(m.BytesArr, decodedData...)
		}
	}
}

func (m *ByteBuilder) AppendBase58(b58StrArr ...string) {
	for idx := range b58StrArr {
		if b58StrArr[idx] == "" {
			continue
		}

		m.BytesArr = append(m.BytesArr, base58.Decode(b58StrArr[idx])...)
	}
}

// old style

func (m *ByteBuilder) GetBytes() []byte {
	return m.BytesArr
}

func (m *ByteBuilder) GetString() string {
	return m.String()
}

// new style

func (m *ByteBuilder) Bytes() []byte {
	return m.BytesArr
}

func (m *ByteBuilder) String() string {
	if m.BytesArr == nil {
		return ""
	}

	return string(m.BytesArr)
}

func (m *ByteBuilder) Len() int {
	if m.BytesArr == nil {
		return 0
	}

	return len(m.BytesArr)
}

func (m *ByteBuilder) Base58() string {
	if m.BytesArr == nil {
		return ""
	}

	return base58.Encode(m.BytesArr)
}

func (m *ByteBuilder) Base64() string {
	if m.BytesArr == nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(m.BytesArr)
}

func (m *ByteBuilder) Hex() string {
	if m.BytesArr == nil {
		return ""
	}

	return hex.EncodeToString(m.BytesArr)
}
