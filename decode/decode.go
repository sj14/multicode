package decode

import (
	"encoding/base64"
	"encoding/hex"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/sj14/multicode/protodec"
)

type Encryption int

const (
	None = iota
	Proto
	Hex
	Base64
)

// Decode the given input as proto message, hex or base64 (applied in this order).
func Decode(input []byte) ([]byte, Encryption) {
	unmarshalled := &protodec.Empty{}

	if err := proto.Unmarshal(input, unmarshalled); err == nil {
		// TODO: remove control characters (unfortunately, they are all valid strings here)
		input = []byte(unmarshalled.String())
		return []byte(unmarshalled.String()), Proto
	}

	if b, err := hex.DecodeString(strings.TrimSpace(strings.TrimPrefix(string(input), "0x"))); err == nil {
		return b, Hex
	}

	// TODO: many false-positives. Decodes it when no base64 was given.
	// Keep it as one of the last decodings. Maybe even 'continue' on the
	// applied decoding before, so e.g. nested hex encodings won't reach here this early.
	if b, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(input))); err == nil {
		return b, Base64
	}

	return input, None
}
