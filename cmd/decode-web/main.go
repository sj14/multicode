package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/sj14/multicode/decode"
)

// TODO: buttons to enable/disable specific decodings
func main() {
	http.HandleFunc("/", handleDecode)

	// TODO: use custom server with timeout
	// TODO: address/port as flag
	log.Fatalln(http.ListenAndServe(":8080", nil))
}

func handleDecode(w http.ResponseWriter, r *http.Request) {
	// TODO: cache the template
	t := template.Must(template.ParseFiles("./cmd/decode-web/index.html"))

	input := r.FormValue("input")
	input = strings.TrimSpace(input)

	var (
		decoder    = decode.New()
		result     = []byte(input)
		enc        decode.Encryption
		appliedEnc []decode.Encryption
	)
	for result, enc = decoder.Decode(result); enc != decode.None; result, enc = decoder.Decode(result) {
		appliedEnc = append(appliedEnc, enc)
	}

	if input != "" && bytes.Equal([]byte(input), result) {
		result = []byte("Failed to decode!")
	}

	data := struct {
		Input       string
		Decoded     string
		Encryptions []decode.Encryption
	}{
		Input:       input,
		Decoded:     string(result),
		Encryptions: appliedEnc,
	}
	if err := t.Execute(w, data); err != nil {
		log.Printf("failed to execute template: %v", err)
	}
}
