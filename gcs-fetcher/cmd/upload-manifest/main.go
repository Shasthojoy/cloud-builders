package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
	"github.com/golang/glog"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"

	"github.com/GoogleCloudPlatform/cloud-builders/gcs-fetcher/pkg/uploader"
)

const (
	userAgent = "gcs-uploader"
)

var (
	dir    = flag.String("dir", ".", "Directory of files to upload")
	bucket = flag.String("bucket", "", "GCS bucket to upload files and manifest to")
)

func main() {
	flag.Parse()

	if *bucket == "" {
		log.Fatalf("--bucket must be specified")
	}

	ctx := context.Background()
	hc, err := buildHTTPClient(ctx)
	if err != nil {
		glog.Info(err)
		os.Exit(2)
	}

	client, err := storage.NewClient(ctx, option.WithHTTPClient(hc), option.WithUserAgent(userAgent))
	if err != nil {
		glog.Infof("Failed to create new GCS client: %v", err)
		os.Exit(2)
	}

	u := uploader.GCSUploader{
		GCS:    realGCS{client},
		OS:     realOS{},
		Root:   *dir,
		Bucket: *bucket,
	}

	manifest, err := u.Upload(ctx)
	if err != nil {
		log.Fatalf("Failed to upload: %v", err)
	}

	log.Println("Uploaded manifest: %s", manifest)
}

func buildHTTPClient(ctx context.Context) (*http.Client, error) {
	hc, err := google.DefaultClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create default client: %v", err)
	}

	ts, err := google.DefaultTokenSource(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create default token source: %v", err)
	}

	hc.Transport = &oauth2.Transport{
		Base:   http.DefaultTransport,
		Source: ts,
	}

	return hc, nil
}

// realGCS is a wrapper over the GCS client functions.
type realGCS struct {
	client *storage.Client
}

func (gp realGCS) NewWriter(ctx context.Context, bucket, object string) io.WriteCloser {
	return gp.client.Bucket(bucket).Object(object).
		If(storage.Conditions{DoesNotExist: true}). // Skip upload if already exists.
		NewWriter(ctx)
}

// realOS merely wraps the os package implementations.
type realOS struct{}

func (realOS) Walk(root string, fn filepath.WalkFunc) error {
	return filepath.Walk(root, fn)
}
