package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/sj14/multicode/decode"
)

const tmpl string = `<!DOCTYPE html>
<html>

<head>
    <meta charset="UTF-8">
    <title>Decode hex, base64 and protobuf recursively</title>
</head>

<body>
    <form action="/">
        <textarea rows="10" cols="100" name="input" placeholder="Your encoded input." required>{{.Input}}</textarea>
        <br>
        <input type="submit">
    </form>


    <br>
    <textarea rows="10" cols="100" placeholder="The decoded result will be shown here." disabled>{{.Decoded}}</textarea>

    <br>
    {{if .Encryptions}}
    Applied Decodings:
    <ol>
        {{range .Encryptions}}
        <li> {{.}} </li>
        {{end}}
    </ol>

    {{end}}

    <br>
    Examples:
	<ul>
		<li> <a
				href="/?input=01010100+01101000+01100101+00100000+01110001+01110101+01101001+01100011+01101011+00100000+01100010+01110010+01101111+01110111+01101110+00100000+11110000+10011111+10100110+10001010+00100000+01101010+01110101+01101101+01110000+01110011+00100000+01101111+01110110+01100101+01110010+00100000+00110001+00110011+00100000+01101100+01100001+01111010+01111001+00100000+11110000+10011111+10010000+10110110+00101110">Bits</a>
		</li>
        <li> <a
                href="/?input=84+104+101+32+113+117+105+99+107+32+98+114+111+119+110+32+240+159+166+138+32+106+117+109+112+115+32+111+118+101+114+32+49+51+32+108+97+122+121+32+240+159+144+182+46">Bytes</a>
        </li>
        <li> <a
                href="/?input=54686520717569636b2062726f776e20f09fa68a206a756d7073206f766572203133206c617a7920f09f90b62e">Hex</a>
        </li>
        <li> <a href="/?input=VGhlIHF1aWNrIGJyb3duIPCfpooganVtcHMgb3ZlciAxMyBsYXp5IPCfkLYu">Base64</a></li>
        <li> <a
                href="/?input=NTQ2ODY1MjA3MTc1Njk2MzZiMjA2MjcyNmY3NzZlMjBmMDlmYTY4YTIwNmE3NTZkNzA3MzIwNmY3NjY1NzIyMDMxMzMyMDZjNjE3YTc5MjBmMDlmOTBiNjJl">Base64
                and Hex</a></li>
        <li> <a href="/?input=0a086d79207175657279102a2006">Proto</a>
        </li>
    </ul>

    <a href="https://github.com/sj14/multicode">Source-Code</a>
</body>

</html>`

// TODO: buttons to enable/disable specific decodings
func main() {
	http.HandleFunc("/", handleDecode)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("listening on port %v\n", port)

	// TODO: use custom server with timeout
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}

func handleDecode(w http.ResponseWriter, r *http.Request) {
	// TODO: cache the template
	t, err := template.New("index").Parse(tmpl)
	if err != nil {
		log.Printf("failed to parse template: %v", err)
	}
	input := r.FormValue("input")
	input = strings.TrimSpace(input)

	var (
		decoder    = decode.New()
		result     = []byte(input)
		enc        decode.Encoding
		appliedEnc []decode.Encoding
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
		Encryptions []decode.Encoding
	}{
		Input:       input,
		Decoded:     string(result),
		Encryptions: appliedEnc,
	}
	if err := t.Execute(w, data); err != nil {
		log.Printf("failed to execute template: %v", err)
	}
}
