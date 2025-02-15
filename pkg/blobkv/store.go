package blobkv

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type KeyValue struct {
	Key       string
	Value     []byte
	Version   int64
	CreateRev int64
	ModRev    int64
	Lease     int64
}

type Store struct {
	client     *minio.Client
	bucket     string
	createRev  int64
	modRev     int64
	watchChans map[int64]chan *KeyValue
}

func NewStore(blobURL string) (*Store, error) {
	parsed, err := ParseURL(blobURL)
	if err != nil {
		return nil, fmt.Errorf("invalid blob URL: %w", err)
	}

	endpoint := fmt.Sprintf("%s:%s", parsed.Host, parsed.Port)
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(parsed.Username, parsed.Password, ""),
		Secure: parsed.Secure,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, parsed.Bucket)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = client.MakeBucket(ctx, parsed.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}

	return &Store{
		client:     client,
		bucket:     parsed.Bucket,
		watchChans: make(map[int64]chan *KeyValue),
	}, nil
}

func (s *Store) Get(ctx context.Context, key string) (*KeyValue, error) {
	obj, err := s.client.GetObject(ctx, s.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	var kv KeyValue
	err = json.NewDecoder(obj).Decode(&kv)
	if err != nil {
		return nil, err
	}

	return &kv, nil
}

func (s *Store) Put(ctx context.Context, kv *KeyValue) error {
	data, err := json.Marshal(kv)
	if err != nil {
		return err
	}

	_, err = s.client.PutObject(ctx, s.bucket, kv.Key, bytes.NewBuffer(data), int64(len(data)), minio.PutObjectOptions{
		ContentType: "application/json",
	})
	if err != nil {
		return err
	}

	s.notifyWatchers(kv)

	return nil
}

func (s *Store) Delete(ctx context.Context, key string) error {
	return s.client.RemoveObject(ctx, s.bucket, key, minio.RemoveObjectOptions{})
}

func (s *Store) Watch(ctx context.Context, key string) (<-chan *KeyValue, error) {
	ch := make(chan *KeyValue, 100)
	id := time.Now().UnixNano()
	s.watchChans[id] = ch

	go func() {
		<-ctx.Done()
		delete(s.watchChans, id)
		close(ch)
	}()

	return ch, nil
}

func (s *Store) notifyWatchers(kv *KeyValue) {
	for _, ch := range s.watchChans {
		select {
		case ch <- kv:
		default:
		}
	}
}
