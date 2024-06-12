package decode

import (
	"encoding/base64"
	"encoding/hex"
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// Encoding handled by the decoding.
type Encoding string

const (
	// None decoding.
	None = "none"
	// Bit decoding.
	Bit = "bit"
	// Byte decoding.
	Byte = "byte"
	// Hex decoding.
	Hex = "hex"
	// Base64 decoding.
	Base64 = "base64"
	// Proto decoding.
	Proto = "proto"
)

// Decoder implementation.
type Decoder struct {
	proto   bool
	bitDec  bool
	byteDec bool
	hex     bool
	base64  bool
}

var defaultDecoder = Decoder{
	proto:   true,
	bitDec:  true,
	byteDec: true,
	hex:     true,
	base64:  true,
}

// Option for decoding.
type Option func(*Decoder)

// WithoutAll disables all decodings.
func WithoutAll() Option {
	return func(d *Decoder) {
		d.proto = false
		d.byteDec = false
		d.hex = false
		d.base64 = false
	}
}

// WithBit decoding.
func WithBit() Option {
	return func(d *Decoder) {
		d.bitDec = true
	}
}

// WithByte decoding.
func WithByte() Option {
	return func(d *Decoder) {
		d.byteDec = true
	}
}

// WithHex decoding.
func WithHex() Option {
	return func(d *Decoder) {
		d.hex = true
	}
}

// WithBase64 decoding.
func WithBase64() Option {
	return func(d *Decoder) {
		d.base64 = true
	}
}

// WithProto decoding.
func WithProto() Option {
	return func(d *Decoder) {
		d.proto = true
	}
}

// New decoder with all decodings enabled by default.
func New(opts ...Option) Decoder {
	d := defaultDecoder
	for _, o := range opts {
		o(&d)
	}
	return d
}

// DecodeAll decodes the given input recursively as long as a decoding was applied.
func DecodeAll(input []byte, opts ...Option) []byte {
	if len(input) == 0 {
		return []byte{}
	}

	var (
		decoder = New(opts...)
		result  = input
		enc     Encoding
	)

	for result, enc = decoder.Decode(result); enc != None; result, enc = decoder.Decode(result) {
		// continue decoding as long a a decoder was applied (not 'None')
	}
	return result
}

// Decode the given input as proto message, hex or base64 (applied in this order).
func (d *Decoder) Decode(input []byte) ([]byte, Encoding) {
	if len(input) == 0 {
		return []byte{}, None
	}

	unmarshalled := &anypb.Any{}

	if d.proto {
		if err := proto.Unmarshal(input, unmarshalled); err == nil {
			// TODO: remove control characters (unfortunately, they are all valid strings here)
			return []byte(unmarshalled.String()), Proto
		}
	}

	if d.bitDec {
		byteIn := strings.Trim(string(input), "[]") // [32 87 111 114 108 100] -> 32 87 111 114 108 100

		if b, err := Base2AsBytes(byteIn); err == nil {
			return b, Bit
		}
	}

	// byte before hex, hex might contains letters, which are not valid in byte dec
	if d.byteDec {
		byteIn := strings.Trim(string(input), "[]") // [32 87 111 114 108 100] -> 32 87 111 114 108 100

		if b, err := Base10AsBytes(byteIn); err == nil {
			return b, Byte
		}
	}

	// hex after byte
	if d.hex {
		hexIn := strings.TrimSpace(string(input))   // e.g. new line
		hexIn = strings.TrimPrefix(hexIn, "0x")     // hex prefix
		hexIn = strings.Replace(hexIn, " ", "", -1) // bd b2 3d bc 20 e2 8c 98 -> bdb23dbc20e28c98

		if b, err := hex.DecodeString(hexIn); err == nil {
			return b, Hex
		}
	}

	// TODO: many false-positives. Decodes it when no base64 was given.
	// Keep it as one of the last decodings.
	if d.base64 {
		if b, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(input))); err == nil {
			return b, Base64
		}
	}

	return input, None
}

func Base10AsBytes(input string) ([]byte, error) {
	input = strings.TrimSpace(input)
	splitted := strings.Split(input, " ")
	var result []byte

	for _, i := range splitted {
		byteAsInt, err := strconv.Atoi(i)
		if err != nil {
			return nil, err
		}
		result = append(result, byte(byteAsInt))
	}
	return result, nil
}

func Base2AsBytes(input string) ([]byte, error) {
	input = strings.TrimSpace(input)
	splitted := strings.Split(input, " ")
	var result []byte

	for _, i := range splitted {
		byteAsInt, err := strconv.ParseInt(i, 2, 0)
		if err != nil {
			return nil, err
		}
		result = append(result, byte(byteAsInt))
	}
	return result, nil
}
