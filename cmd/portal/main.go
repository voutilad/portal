package main

import (
	"fmt"
	"github.com/voutilad/portal"
	"io"
	"log"
	"os"
	"syscall"
)

const blocksize = 1 << 20

func main() {
	err := syscall.Mkfifo("/tmp/junk", 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer syscall.Unlink("/tmp/junk")

	pipe, err := os.OpenFile("/tmp/junk", os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer pipe.Close()

	src, err := portal.NewGcsPortal("my-dumb-bucket-69696", "gross-domestic-product-march-2020-quarter.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer src.Close()

	cnt, err := io.Copy(pipe, src)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(os.Stderr, "wrote %d bytes\n", cnt)
}
