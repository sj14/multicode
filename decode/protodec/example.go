package protodec

import (
	fmt "fmt"
	"log"

	proto "github.com/golang/protobuf/proto"
)

func PrintExample() {
	pb := ComplexMessage{
		Query:      "my query",
		PageNumber: 42,
		Corpus:     ComplexMessage_VIDEO,
	}

	b, err := proto.Marshal(&pb)
	if err != nil {
		log.Fatalf("failed to marhsal example: %v", err)
	}
	fmt.Printf("%x", b)
}
