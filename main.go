package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/sj14/multicode/decode"
)

var (
	verbose bool
)

func main() {
	// init flags
	flag.BoolVar(&verbose, "v", false, "verbose ouput mode")
	flag.Parse()

	// read program input
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadBytes('\n')
	if err != nil {
		log.Fatalln("failed to read input")
	}

	// decoding
	result := input
	var enc decode.Encryption
	for result, enc = decode.Decode(result); enc != decode.None; result, enc = decode.Decode(result) {
		logVerbose("applied decoding %v:\n%s\n\n", enc, result)
	}

	// check if any kind of decryption was applied
	if bytes.Compare(input, result) == 0 {
		log.Fatalln("failed to decode")
	}

	// output result
	logVerbose("result:\n")
	fmt.Printf("%v\n", string(result))
}

func logVerbose(format string, v ...interface{}) {
	if !verbose {
		return
	}
	fmt.Printf(format, v...)
}
