package portal

import (
	"cloud.google.com/go/storage"
	"context"
	"io"
)

type Portal interface {
	io.ReadCloser
	io.WriterTo
}

type gcsPortal struct {
	client     *storage.Client
	src        *storage.Reader
	BucketName string
	ObjectName string
}

func NewGcsPortal(bucket, object string) (Portal, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	bkt := client.Bucket(bucket)
	obj := bkt.Object(object)
	src, err := obj.NewReader(ctx)
	if err != nil {
		return nil, err
	}

	portal := &gcsPortal{
		client:     client,
		src:        src,
		BucketName: bucket,
		ObjectName: object,
	}
	return portal, nil
}

func (portal gcsPortal) Read(p []byte) (n int, err error) {
	return portal.src.Read(p)
}

func (portal gcsPortal) Close() error {
	return portal.src.Close()
}

func (portal gcsPortal) WriteTo(w io.Writer) (n int64, err error) {
	var cnt int64
	src := portal.src

	for src.Remain() > 0 {
		c, err := io.CopyN(w, src, 4096)
		if err != nil && err != io.EOF {
			return cnt, err
		}
		cnt += c
	}

	return cnt, nil
}
