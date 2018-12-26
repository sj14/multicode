package main

import (
	"bufio"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/sj14/multidecode/protodec"
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
	result, _ := decode(input)

	logVerbose("result:\n")
	fmt.Printf("%v\n", string(result))
}

func decode(input []byte) (string, int) {
	unmarshalled := &protodec.Empty{}
	appliedCount := 0

	for {
		applied := false
		if err := proto.Unmarshal(input, unmarshalled); err == nil {
			// TODO: remove control characters (unfortunately, they are all valid strings here)
			input = []byte(unmarshalled.String())
			logVerbose("applied proto decoding:\n%s\n\n", unmarshalled.String())
			// we can't decode more on top of proto
			appliedCount++
			return unmarshalled.String(), appliedCount
		}

		// TODO: need to fix this
		if b, err := base64.StdEncoding.DecodeString(string(input)); err == nil {
			input = b
			logVerbose("applied base64 decoding:\n%v\n\n", string(b))
			appliedCount++
			applied = true
		}

		if b, err := hex.DecodeString(strings.ToLower(strings.TrimPrefix(string(input), "0x"))); err == nil {
			input = b
			logVerbose("applied hex decoding:\n%v\n\n", string(b))
			appliedCount++
			applied = true
		}

		if !applied {
			return string(input), appliedCount
		}
	}
}

func logVerbose(format string, v ...interface{}) {
	if !verbose {
		return
	}
	fmt.Printf(format, v...)
}
