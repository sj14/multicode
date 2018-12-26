package main

import (
	"encoding/base64"
	"encoding/hex"
	"testing"
)

func TestDecode(t *testing.T) {
	testCases := []struct {
		description string
		encoded     []byte
		applied     int
		decoded     string
	}{
		{
			description: "base64",
			encoded:     []byte(base64.StdEncoding.EncodeToString([]byte("This is a base64 test"))),
			applied:     1,
			decoded:     "This is a base64 test",
		},
		{
			description: "hex",
			encoded:     []byte(hex.EncodeToString([]byte("This is a hex test"))),
			applied:     1,
			decoded:     "This is a hex test",
		},
		{
			description: "0x hex",
			encoded:     []byte("0x" + hex.EncodeToString([]byte("This is a 0xHEX test"))),
			applied:     1,
			decoded:     "This is a 0xHEX test",
		},
		{
			description: "proto/hex",
			encoded:     []byte("0a03617364102a20042a200a1468747470733a2f2f657861206d706c652e636f6d1208657861206d706c653a00"),
			applied:     2,
			decoded:     `1:"asd" 2:42 4:4 5:"\n\x14https://exa mple.com\x12\bexa mple" 7:"" `,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.description, func(t *testing.T) {
			decoded, n := decode(tt.encoded)
			if n != tt.applied {
				t.Fatalf("expeced %v decodings, got %v", tt.applied, n)
			}
			if decoded != tt.decoded {
				t.Fatalf("expeced '%v', got '%v'", tt.decoded, decoded)
			}
		})
	}
}
