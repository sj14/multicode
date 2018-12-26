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
	flag.BoolVar(&verbose, "v", false, "verbose ouput mode")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadBytes('\n')
	if err != nil {
		log.Fatalln("failed to read input")
	}

	result := input
	var enc decode.Encryption
	for result, enc = decode.Decode(result); enc != decode.None; result, enc = decode.Decode(result) {
		logVerbose("applied decoding %v:\n%s\n\n", enc, result)
	}

	if bytes.Compare(input, result) == 0 {
		fmt.Println("failed to decode")
		os.Exit(1)
	}

	logVerbose("result:\n")
	fmt.Printf("%v\n", string(result))
}

func logVerbose(format string, v ...interface{}) {
	if !verbose {
		return
	}
	fmt.Printf(format, v...)
}
