package bytesbuilder

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"reflect"

	"github.com/btcsuite/btcutil/base58"
)

type ByteBuilder struct {
	BytesArr []byte
}

func NewBuilder() *ByteBuilder {
	b := ByteBuilder{}
	return &b
}

func (builder *ByteBuilder) Clear() {
	builder.BytesArr = nil
}

//TODO : exception handling

func (builder *ByteBuilder) Append(data interface{}) {
	var byteData []byte
	switch t := data.(type) {
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
	builder.BytesArr = append(builder.BytesArr, byteData...)
}

func (builder *ByteBuilder) AppendHex(hexStr string) {
	data, _ := hex.DecodeString(hexStr[2:])
	builder.BytesArr = append(builder.BytesArr, data...)
}

func (builder *ByteBuilder) AppendBase64(b64Str string) {
	decodedData, _ := base64.StdEncoding.DecodeString(b64Str)
	builder.BytesArr = append(builder.BytesArr, decodedData...)
}

func (builder *ByteBuilder) AppendBase58(b58Str string) {
	decodedData := base58.Decode(b58Str)
	builder.BytesArr = append(builder.BytesArr, decodedData...)
}

func (builder *ByteBuilder) GetBytes() []byte {
	return builder.BytesArr
}

func (builder *ByteBuilder) GetString() string {
	return string(builder.BytesArr)
}
