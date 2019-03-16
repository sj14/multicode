package decode

import (
	"encoding/base64"
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/sj14/multicode/decode/protodec"
)

// Encoding handled by the decoding.
type Encoding string

const (
	// None decoding.
	None = "none"
	// Proto decoding.
	Proto = "proto"
	// Hex decoding.
	Hex = "hex"
	// Byte decoding.
	Byte = "byte"
	// Base64 decoding.
	Base64 = "base64"
)

// Decoder implementation.
type Decoder struct {
	proto   bool
	byteDec bool
	hex     bool
	base64  bool
}

var defaultDecoder = Decoder{
	proto:   true,
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

// WithHex decoding.
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

	unmarshalled := &protodec.Empty{}

	if d.proto {
		if err := proto.Unmarshal(input, unmarshalled); err == nil {
			// TODO: remove control characters (unfortunately, they are all valid strings here)
			return []byte(unmarshalled.String()), Proto
		}
	}

	if d.hex {
		if b, err := hex.DecodeString(strings.TrimSpace(strings.TrimPrefix(string(input), "0x"))); err == nil {
			return b, Hex
		}
	}

	if d.byteDec {
		if b, err := AsBytes(string(input)); err == nil {
			return b, Byte
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

func AsBytes(input string) ([]byte, error) {
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
