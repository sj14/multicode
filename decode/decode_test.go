package decode

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"testing"
)

func TestDecode(t *testing.T) {
	type expect struct {
		encryption Encoding
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
			description: "bit",
			given:       []byte("01010100 01101000 01100101 00100000 01110001 01110101 01101001 01100011 01101011 00100000 01100010 01110010 01101111 01110111 01101110 00100000 11110000 10011111 10100110 10001010 00100000 01101010 01110101 01101101 01110000 01110011 00100000 01101111 01110110 01100101 01110010 00100000 00110001 00110011 00100000 01101100 01100001 01111010 01111001 00100000 11110000 10011111 10010000 10110110 00101110"),
			expect: expect{
				encryption: Bit,
				output:     []byte("The quick brown 🦊 jumps over 13 lazy 🐶."),
			},
		},
		{
			description: "bytes",
			given:       []byte("84 104 101 32 113 117 105 99 107 32 98 114 111 119 110 32 240 159 166 138 32 106 117 109 112 115 32 111 118 101 114 32 49 51 32 108 97 122 121 32 240 159 144 182 46"),
			expect: expect{
				encryption: Byte,
				output:     []byte("The quick brown 🦊 jumps over 13 lazy 🐶."),
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
			description: "hex/base64",
			given:       []byte("Njg2NTZDNkM2RjIwNzQ2ODY1NzI2NQ=="),
			decodings:   2,
			expect: expect{
				encryption: Hex,
				output:     []byte("hello there"),
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
				t.Fatalf("expected decoding %v, got %v", tt.expect.encryption, enc)
			}
			if !bytes.Equal(decoded, tt.expect.output) {
				t.Fatalf("expected '%v', got '%v'", tt.expect.output, decoded)
			}
		})
	}
}
