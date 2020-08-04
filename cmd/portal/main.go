package main

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
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

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	bkt := client.Bucket("my-dumb-bucket-69696")
	obj := bkt.Object("gross-domestic-product-march-2020-quarter.csv")
	src, err := obj.NewReader(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer src.Close()

	for src.Remain() > 0 {
		cnt, err := io.CopyN(pipe, src, blocksize)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		fmt.Fprintf(os.Stderr, "wrote %d bytes\n", cnt)
	}

}
