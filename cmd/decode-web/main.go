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

	result := decode.DecodeAll([]byte(input))

	if input != "" && bytes.Equal([]byte(input), result) {
		result = []byte("failed to decode")
	}

	data := struct {
		Input   string
		Decoded string
	}{
		Input:   input,
		Decoded: string(result),
	}
	if err := t.Execute(w, data); err != nil {
		log.Printf("failed to execute template: %v", err)
	}
}
