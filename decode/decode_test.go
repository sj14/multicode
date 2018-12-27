package decode

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"testing"
)

func TestDecode(t *testing.T) {
	type expect struct {
		encryption Encryption
		output     []byte
	}

	testCases := []struct {
		description string
		given       []byte
		decodings   int
		expect      expect
	}{
		{
			description: "no encoding",
			given:       []byte("no encoding"),
			expect: expect{
				encryption: None,
				output:     []byte("no encoding"),
			},
		},
		{
			description: "base64",
			given:       []byte(base64.StdEncoding.EncodeToString([]byte("This is a base64 test"))),
			expect: expect{
				encryption: Base64,
				output:     []byte("This is a base64 test"),
			},
		},
		{
			description: "hex",
			given:       []byte(hex.EncodeToString([]byte("This is a hex test"))),
			expect: expect{
				encryption: Hex,
				output:     []byte("This is a hex test"),
			},
		},
		{
			description: "0x hex",
			given:       []byte("0x" + hex.EncodeToString([]byte("This is a 0xHEX test"))),
			expect: expect{
				encryption: Hex,
				output:     []byte("This is a 0xHEX test"),
			},
		},
		{
			description: "proto/hex",
			given:       []byte("0a03617364102a20042a200a1468747470733a2f2f657861206d706c652e636f6d1208657861206d706c653a00"),
			decodings:   2,
			expect: expect{
				encryption: Proto,
				output:     []byte(`1:"asd" 2:42 4:4 5:"\n\x14https://exa mple.com\x12\bexa mple" 7:"" `),
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.description, func(t *testing.T) {
			// initial decoding
			decoder := New()
			decoded, enc := decoder.Decode(tt.given)

			// apply as many decodings as specified
			for i := 1; i < tt.decodings; i++ {
				decoded, enc = decoder.Decode(decoded)
			}

			// assert results
			if enc != tt.expect.encryption {
				t.Fatalf("expeced decoding %v, got %v", tt.expect.encryption, enc)
			}
			if bytes.Compare(decoded, tt.expect.output) != 0 {
				t.Fatalf("expected '%v', got '%v'", tt.expect.output, decoded)
			}
		})
	}
}
