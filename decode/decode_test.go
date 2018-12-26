package decode

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"testing"
)

func TestDecode(t *testing.T) {
	testCases := []struct {
		description string
		encoded     []byte
		decodings   int
		applied     Encryption
		decoded     []byte
	}{
		{
			description: "no encoding",
			encoded:     []byte("no encoding"),
			applied:     None,
			decoded:     []byte("no encoding"),
		},
		{
			description: "base64",
			encoded:     []byte(base64.StdEncoding.EncodeToString([]byte("This is a base64 test"))),
			applied:     Base64,
			decoded:     []byte("This is a base64 test"),
		},
		{
			description: "hex",
			encoded:     []byte(hex.EncodeToString([]byte("This is a hex test"))),
			applied:     Hex,
			decoded:     []byte("This is a hex test"),
		},
		{
			description: "0x hex",
			encoded:     []byte("0x" + hex.EncodeToString([]byte("This is a 0xHEX test"))),
			applied:     Hex,
			decoded:     []byte("This is a 0xHEX test"),
		},
		{
			description: "proto/hex",
			encoded:     []byte("0a03617364102a20042a200a1468747470733a2f2f657861206d706c652e636f6d1208657861206d706c653a00"),
			decodings:   2,
			applied:     Proto,
			decoded:     []byte(`1:"asd" 2:42 4:4 5:"\n\x14https://exa mple.com\x12\bexa mple" 7:"" `),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.description, func(t *testing.T) {
			// initial decoding
			decoded, n := Decode(tt.encoded)

			// apply as many decodings as specified
			for i := 1; i < tt.decodings; i++ {
				decoded, n = Decode(decoded)
			}

			// assert results
			if n != tt.applied {
				t.Fatalf("expeced decoding %v, got %v", tt.applied, n)
			}
			if bytes.Compare(decoded, tt.decoded) != 0 {
				t.Fatalf("expected '%v', got '%v'", tt.decoded, decoded)
			}
		})
	}
}
