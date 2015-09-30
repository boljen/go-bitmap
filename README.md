# Bitmap (Go)

Package bitmap implements (thread-safe) bitmap functions and abstractions

## Install

    go get github.com/boljen/go-bitmap

## Documentation

See [godoc](https://godoc.org/github.com/boljen/go-bitmap)

## Example

    package main

    import (
        "fmt"
        "github.com/boljen/go-bitmap"
    )

    func main() {
        bm := bitmap.New(100)
        bm.Set(0, true)
        fmt.Println(bm.Get(0))
    }

## License

MIT
