package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sj14/multicode/decode"
)

var (
	version = "dev version"
	commit  = "none"
	date    = "unknown"
)

func main() {
	var (
		bitDec      = flag.Bool("bit", true, "use bit decoding")
		byteDec     = flag.Bool("byte", true, "use byte decoding")
		hexDec      = flag.Bool("hex", true, "use hex decoding")
		base64Dec   = flag.Bool("base64", true, "use base64 decoding")
		protoDec    = flag.Bool("proto", true, "use proto decoding")
		verbose     = flag.Bool("verbose", false, "verbose output mode")
		versionFlag = flag.Bool("version", false, fmt.Sprintf("print version information of this release (%v)", version))
		// none := flag.Bool("none", false, "disable all decodings") // TODO: not working yet
	)
	flag.Parse()

	// print version info and exit
	if *versionFlag {
		fmt.Printf("version: %v\n", version)
		fmt.Printf("commit: %v\n", commit)
		fmt.Printf("date: %v\n", date)
		os.Exit(0)
	}

	var input []byte

	// read program input
	if flag.NArg() == 0 { // from stdin (also pipe)
		reader := bufio.NewReader(os.Stdin)
		var err error
		input, err = reader.ReadBytes('\n')
		if err != nil {
			log.Fatalln("failed to read input")
		}
	} else { // from argument
		if flag.NArg() > 1 {
			log.Fatalln("takes at most one input")
		}
		input = []byte(flag.Arg(0))
	}
	// can't trim spaces in general, this would mess up proto decoding!
	if strings.TrimSpace(string(input)) == "" {
		log.Fatalln("empty input")
	}

	// Default decoder enables all decodings.
	// Disable all and only enable specified ones.
	// Flags are set to true by default.
	var opts []decode.Option
	opts = append(opts, decode.WithoutAll())

	// Enable specifified decodings.
	if *bitDec {
		opts = append(opts, decode.WithByte())
	}
	if *byteDec {
		opts = append(opts, decode.WithByte())
	}
	if *hexDec {
		opts = append(opts, decode.WithHex())
	}
	if *base64Dec {
		opts = append(opts, decode.WithBase64())
	}
	if *protoDec {
		opts = append(opts, decode.WithProto())
	}

	// decoding
	var (
		decoder = decode.New(opts...)
		result  = input
		enc     decode.Encoding
	)
	for result, enc = decoder.Decode(result); enc != decode.None; result, enc = decoder.Decode(result) {
		logVerbose(*verbose, "- applied decoding '%v':\n%s\n\n", enc, result)
	}

	// check if any kind of decryption was applied
	if bytes.Equal(input, result) {
		log.Fatalln("failed to decode")
	}

	// output result
	logVerbose(*verbose, "- result:\n")
	fmt.Printf("%v\n", string(result))
}

func logVerbose(verbose bool, format string, v ...interface{}) {
	if !verbose {
		return
	}
	fmt.Printf(format, v...)
}
