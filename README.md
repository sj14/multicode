# multicode

`multicode` allows to input a (nested) `base64`, `hex` or `proto` (protocol buffers) decoded sequence and will recursively try to encode it. This is helpful when you get encoded data but don't exactly know how it was encoded or encoding might lead to cumbersome command concatenation.

## Examples

First, let's encode a string with hex and base64 encoding:

```text
$ echo hello there | xxd -p | base64
Njg2NTZjNmM2ZjIwNzQ2ODY1NzI2NTBhCg==
```

Decode:

```text
$ decode Njg2NTZjNmM2ZjIwNzQ2ODY1NzI2NTBhCg==
hello there
```

Decode using the pipe:

```text
$ echo Njg2NTZjNmM2ZjIwNzQ2ODY1NzI2NTBhCg== | decode
hello there
```

Decode in verbose mode:

```text
$ decode -v Njg2NTZjNmM2ZjIwNzQ2ODY1NzI2NTBhCg==
- applied decoding 'base64':
68656C6C6F207468657265

- applied decoding 'hex':
hello there

- result:
hello there
```

Disable hex decoding:

```text
$ decode -v -hex=false Njg2NTZjNmM2ZjIwNzQ2ODY1NzI2NTBhCg==
- applied decoding 'base64':
68656C6C6F207468657265

- result:
68656C6C6F207468657265
```

## Usage

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
