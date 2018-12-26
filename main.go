package main

import (
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
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

	input := []byte(os.Args[len(os.Args)-1])

	unmarshalled := &protodec.Empty{}
	for {
		applied := false

		if err := proto.Unmarshal(input, unmarshalled); err == nil {
			// TODO: remove control characters (unfortunately, they are all valid strings here)
			input = []byte(unmarshalled.String())
			logVerbose("applied proto decoding:\n%s\n\n", unmarshalled.String())
			// we can't decode more on top of proto
			break
		}

		if b, err := base64.StdEncoding.DecodeString(string(input)); err == nil {
			input = b
			logVerbose("applied base64 decoding:\n%v\n\n", string(b))
			applied = true
		}

		if b, err := hex.DecodeString(strings.ToLower(strings.TrimPrefix(string(input), "0x"))); err == nil {
			input = b
			logVerbose("applied hex decoding:\n%v\n\n", string(b))
			applied = true
		}

		if !applied {
			break
		}
	}
	logVerbose("result:\n")
	fmt.Printf("%v\n", string(input))
}

func logVerbose(format string, v ...interface{}) {
	if !verbose {
		return
	}
	fmt.Printf(format, v...)
}
