package portal

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"strconv"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

// Represents an abstract Portal source
type Portal interface {
	io.ReadCloser
	io.WriterTo
}

// Implementation of a Google Cloud Storage portal source
type gcsPortal struct {
	client     *storage.Client
	src        *storage.Reader
	BucketName string
	ObjectName string
}

// Constructs a new Portal, backed by an object in Google Cloud Storage
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
	err := portal.src.Close()
	if err != nil {
		return err
	}

	return portal.client.Close()
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

func (portal gcsPortal) String() string {
	return fmt.Sprintf("GCS { bucket: %s, object: %s }",
		portal.BucketName, portal.ObjectName)
}

const LatestVersion = -1

type gsmPortal struct {
	buffer     *bytes.Buffer
	ProjectId  string
	SecretName string
	Version    int
}

func NewGsmPortal(projectId, secretName string, version int) (Portal, error) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	name := "projects/" + projectId + "/secrets/" + secretName + "/versions/"

	if version == LatestVersion {
		name += "latest"
	} else {
		name += strconv.Itoa(version)
	}

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, err
	}

	portal := &gsmPortal{
		buffer:     bytes.NewBuffer(result.Payload.Data),
		ProjectId:  projectId,
		SecretName: secretName,
		Version:    version,
	}

	return portal, nil
}

func (portal gsmPortal) Read(p []byte) (n int, err error) {
	return portal.buffer.Read(p)
}

func (portal gsmPortal) Close() error {
	return nil
}

func (portal gsmPortal) WriteTo(w io.Writer) (n int64, err error) {
	return portal.buffer.WriteTo(w)
}

func (portal gsmPortal) String() string {
	return fmt.Sprintf("GSM { projectId: %s, secret: %s, version %d }",
		portal.ProjectId, portal.SecretName, portal.Version)
}
