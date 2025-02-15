package blobkv

import (
	"fmt"
	"net/url"
	"strings"
)

type BlobURL struct {
	Scheme   string
	Host     string
	Port     string
	Username string
	Password string
	Bucket   string
	Secure   bool
}

func ParseURL(rawURL string) (*BlobURL, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "s3" {
		return nil, fmt.Errorf("unsupported scheme: %s", u.Scheme)
	}

	host := u.Hostname()
	port := u.Port()
	if port == "" {
		port = "9000"
	}

	bucket := strings.TrimPrefix(u.Path, "/")
	if bucket == "" {
		return nil, fmt.Errorf("bucket name required in URL path")
	}

	password, _ := u.User.Password()

	return &BlobURL{
		Scheme:   u.Scheme,
		Host:     host,
		Port:     port,
		Username: u.User.Username(),
		Password: password,
		Bucket:   bucket,
		Secure:   false, // Could be determined by scheme (s3s://) in the future
	}, nil
}
