package bytesbuilder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	originDataForB64 = "Hello base64"
	base64Str        = "SGVsbG8gYmFzZTY0"
	originDataForB58 = "Hello base58"
	base58Str        = "2NEpo7TZRXJTtJ5BD"
	originDataFroHex = "Hello Hex"
	hexStr           = "0x48656c6c6f20486578"
)

func TestByteBuilder_AppendBaseXX(t *testing.T) {
	builder := NewBuilder()

	builder.AppendBase64(base64Str)
	builder.AppendBase58(base58Str)

	bytesData := builder.GetBytes()
	strData := string(bytesData)

	assert.Equal(t, strData, originDataForB64+originDataForB58, "They should be Equal")
}

func TestByteBuilder_AppendHex(t *testing.T) {
	builder := NewBuilder()

	builder.AppendHex(hexStr)

	bytesData := builder.GetBytes()
	strData := string(bytesData)

	assert.Equal(t, strData, originDataFroHex, "They should be Equal")
}

func TestByteBuilder_CommonAppend(t *testing.T) {
	builder := NewBuilder()

	var u8 uint8
	var u8Arr []uint8
	var i32 int32
	var u32 uint32
	var i64 int64
	var u64 uint64
	var str string
	var boolTrue bool
	var boolFalse bool

	i32 = 0x44332211
	i64 = 0x0011223344556677
	u8 = 0x11
	u8Arr = []uint8{0x99, 0x99}
	u32 = 0x11223344
	u64 = 0x1122334455667788
	str = "AAA"
	boolTrue = true
	boolFalse = false

	builder.Append(i32)
	builder.Append(i64)
	builder.Append(u8)
	builder.Append(u8Arr)
	builder.Append(u32)
	builder.Append(u64)
	builder.Append(str)
	builder.Append(boolTrue)
	builder.Append(boolFalse)

	expected := []byte{
		0x44, 0x33, 0x22, 0x11,
		0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77,
		0x11,
		0x99, 0x99,
		0x11, 0x22, 0x33, 0x44,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88,
		0x41, 0x41, 0x41,
		0x01, 0x00}

	bytesData := builder.GetBytes()

	assert.Equal(t, expected, bytesData, "They should be Equal")
}
