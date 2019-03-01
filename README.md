# multicode

[![Build Status](https://travis-ci.org/sj14/multicode.svg?branch=master)](https://travis-ci.org/sj14/multicode)
[![Go Report Card](https://goreportcard.com/badge/github.com/sj14/multicode)](https://goreportcard.com/report/github.com/sj14/multicode)
[![GoDoc](https://godoc.org/github.com/sj14/multicode/decode?status.png)](https://godoc.org/github.com/sj14/multicode/decode)

`multicode` allows to input a (nested) `base64`, `hex` or `proto` (protocol buffers) decoded sequence and will recursively try to encode it. This is helpful when you get encoded data but don't exactly know how it was encoded or encoding might lead to cumbersome command concatenation.

## Installation

### CLI

``` text
go get -u github.com/sj14/multicode/cmd/decode
```

### Web Interface

[Demo](https://multicode.herokuapp.com/)

``` text
go get -u github.com/sj14/multicode/cmd/decode-web
```

## CLI Usage

``` text
  -base64
        use base64 decoding (default true)
  -hex
        use hex decoding (default true)
  -none
        disable all decodings
  -proto
        use proto decoding (default true)
  -v    verbose ouput mode
```

## CLI Examples

First, let's encode a string with hex and base64 encoding:

``` text
$ echo hello there | xxd -p | base64
Njg2NTZjNmM2ZjIwNzQ2ODY1NzI2NTBhCg==
```

Decode:

``` text
$ decode Njg2NTZjNmM2ZjIwNzQ2ODY1NzI2NTBhCg==
hello there
```

Decode using the pipe:

``` text
$ echo Njg2NTZjNmM2ZjIwNzQ2ODY1NzI2NTBhCg== | decode
hello there
```

Decode in verbose mode:

``` text
$ decode -v Njg2NTZjNmM2ZjIwNzQ2ODY1NzI2NTBhCg==
- applied decoding 'base64':
68656C6C6F207468657265

- applied decoding 'hex':
hello there

- result:
hello there
```

Disable hex decoding:

``` text
$ decode -v -hex=false Njg2NTZjNmM2ZjIwNzQ2ODY1NzI2NTBhCg==
- applied decoding 'base64':
68656C6C6F207468657265

- result:
68656C6C6F207468657265
```

## Protobuf

We can decode protocol buffer encodings without specifying a proto file. Based on the missing definition file, it's unfortunately not possible, to output the field names. Field names will be replaced by the field id.

Let's assume the following proto message:

```text
message Message {
  string query = 1;
  int32 page_number = 2;
  int32 result_per_page = 3;
  enum Corpus {
      UNIVERSAL = 0;
      WEB = 1;
  }
  Corpus corpus = 4;
}
```

And we initialize the message like this:

```go
Message{
  Query:      "my query",
  PageNumber: 42,
  Corpus:     ComplexMessage_NEWS,
}
```

The hex decoded proto message (`0a086d79207175657279102a2006`) will be decoded as:

```text
$ decode 0a086d79207175657279102a2004
1:"my query" 2:42 4:6
```

## Using with Docker

### CLI Version

```text
docker build -f Dockerfile.decode -t decode .
docker run --rm -it decode
```

### Web Version

```text
docker build -f Dockerfile.decode-web -t decode-web .
docker run --rm -it -p 8080:8080 decode-web
```